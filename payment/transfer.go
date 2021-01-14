package payment

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	db "../dbquery"
	util "../util"
	ver "../verification"
)

type reqTransferConfirm struct {
	To        string `json:"to"`
	Cid       int    `json:"coin_idx"`
	Value     string `json:"value"`
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
}

type dataConfirmTransfer struct {
	Pg_conn string
	Data    reqTransferConfirm
	PsnInfo tradePsnInfo
}

func TransferProcessRequest(r *http.Request) (*RequestService, error) {
	decoder := json.NewDecoder(r.Body)
	decoder.UseNumber()

	var req RequestService
	err := decoder.Decode(&req)
	if err != nil {
		return nil, err
	}

	return &RequestService{
		Method: req.Method,
		Params: req.Params,
	}, nil
}
func (rs *RequestService) TransferMethod(url string) (interface{}, error) {
	var rtn interface{}
	var err error = nil

	start := time.Now()
	switch strings.ToUpper(rs.Method) {
	case "CONFIRM":
		rtn, err = processConfirmTransfer(url, rs.Params)
	default:
		err = errors.New("method required or not exist.")
		return rtn, err
	}
	log.Printf("Transfer: %s\t%s", rs.Method, time.Since(start))

	if err != nil {
		return rtn, err
	}

	return rtn, nil
}

func processConfirmTransfer(url string, params requestParameters) (interface{}, error) {
	var resp responsePayment

	s, err := confirmTransferCheckData(url, params)
	if err != nil {
		return nil, err
	}

	resp.Result, err = s.confirmTransferUpdateInfo()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func confirmTransferCheckData(url string, params requestParameters) (*dataConfirmTransfer, error) {
	uid := params.Uid

	psn, info, err := deConfirmTransfer(url, uid, params.Data)
	if err != nil {
		return nil, err
	}

	return &dataConfirmTransfer{
		Pg_conn: url,
		PsnInfo: *psn,
		Data:    *info,
	}, nil
}

func deConfirmTransfer(url, uid, data string) (*tradePsnInfo, *reqTransferConfirm, error) {
	if uid == "" || len(uid) != 40 {
		return nil, nil, errors.New("uid: non-regular value.")
	}

	d, err := db.GetInfoUserSys(url, uid)
	if err != nil {
		fmt.Println("database error:", err)
		return nil, nil, errors.New("database error, please notify the administrator.")
	}

	veri, err := ver.ImportPrivKey(d.Sys_privkey)
	if err != nil {
		return nil, nil, err
	}
	s, err := veri.Decrypt(data)
	if err != nil {
		return nil, nil, err
	}

	var req reqTransferConfirm
	err = json.Unmarshal([]byte(s), &req)
	if err != nil {
		return nil, nil, err
	}

	if strings.Compare(d.Token, req.Token) != 0 {
		return nil, nil, errors.New("check data failure.")
	}

	if strings.Compare(uid, req.To) == 0 {
		return nil, nil, errors.New("trading information error.")
	}

	return &tradePsnInfo{
		Uid:         uid,
		Sys_pubkey:  d.Sys_pubkey,
		Sys_privkey: d.Sys_privkey,
		User_pubkey: d.Psn_pubkey,
	}, &req, nil
}
func (df *dataConfirmTransfer) confirmTransferUpdateInfo() (string, error) {
	to_wid, err := db.GetWidFromUid(df.Pg_conn, df.Data.To)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return "", errors.New("user wallet does not exist.")
		} else {
			fmt.Println("database error:", err)
			return "", errors.New("database error, please notify the administrator.")
		}
	}

	from_wid, err := db.GetWidFromUid(df.Pg_conn, df.PsnInfo.Uid)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return "", errors.New("user wallet does not exist.")
		} else {
			fmt.Println("database error:", err)
			return "", errors.New("database error, please notify the administrator.")
		}
	}

	timestamp := util.Timestamp()
	tid := createTradeID(to_wid, timestamp)
	ttoken := calcTransferTtokenHash(tid, df.Data.To, df.Data.Value, df.Data.Cid, 0, timestamp)

	d, err := db.ConfirmTransferInfo(df.Pg_conn, tid, from_wid, to_wid, df.Data.Value, ttoken, df.Data.Cid, timestamp)
	if err != nil {
		fmt.Println("database error:", err)
		return "", errors.New("database error, please notify the administrator.")
	}

	var resp respConfirmTrading
	resp.Cid = d.Cid
	resp.Balance = d.Balance
	resp.Timestamp = util.Timestamp()

	b, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	veri, err := ver.ImportPubKey(df.PsnInfo.User_pubkey)
	if err != nil {
		return "", err
	}
	pt, err := veri.Encrypt(string(b))
	if err != nil {
		return "", err
	}

	return pt, nil
}

func calcTransferTtokenHash(tid, to, value string, cid int, expire, ctime int64) string {
	s := tid + "&" + to + "&" + value + "&" + strconv.Itoa(cid) + "&" + strconv.FormatInt(expire, 10) + "&" + strconv.FormatInt(ctime, 10)
	return util.CreateSHA1Hash(s)
}
