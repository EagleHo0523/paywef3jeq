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

type reqNolimitCreate struct {
	Cid       int    `json:"coin_idx"`
	Value     string `json:"value"`
	Expire    int64  `json:"expire,omitempty"`
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
}
type reqNolimitConfirm struct {
	From      string `json:"from"`
	Info      string `json:"info"`
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
}

type dataCreateNolimit struct {
	Pg_conn string
	PsnInfo tradePsnInfo
	Data    reqNolimitCreate
}
type dataConfirmNolimit struct {
	Data    reqNolimitConfirm
	Tid     string
	Pg_conn string
	PsnInfo tradePsnInfo
}

type respCreateNolimit struct {
	Tid       string `json:"tid"`
	Timestamp int64  `json:"timestamp"`
}
type respConfirmNolimit struct {
	Cid       int    `json:"coin_idx"`
	Increment string `json:"increment"`
	Balance   string `json:"balance"`
	Timestamp int64  `json:"timestamp"`
}

func NolimitProcessRequest(r *http.Request) (*RequestService, error) {
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
func (rs *RequestService) NolimitMethod(url string) (interface{}, error) {
	var rtn interface{}
	var err error = nil

	start := time.Now()
	switch strings.ToUpper(rs.Method) {
	case "CREATE":
		rtn, err = processCreateNolimit(url, rs.Params)
	case "CONFIRM":
		rtn, err = processConfirmNolimit(url, rs.Params)
	default:
		err = errors.New("method required or not exist.")
		return rtn, err
	}
	log.Printf("Nolimit: %s\t%s", rs.Method, time.Since(start))

	if err != nil {
		return rtn, err
	}

	return rtn, nil
}

func processCreateNolimit(url string, params requestParameters) (interface{}, error) {
	var resp responsePayment

	s, err := createNolimitCheckData(url, params)
	if err != nil {
		return nil, err
	}

	err = s.createNolimitCheck()
	if err != nil {
		return nil, err
	}

	resp.Result, err = s.createNolimitCreateInfo()
	if err != nil {
		return nil, err
	}

	return resp, nil
}
func processConfirmNolimit(url string, params requestParameters) (interface{}, error) {
	var resp responsePayment

	s, err := confirmNolimitCheckData(url, params)
	if err != nil {
		return nil, err
	}

	resp.Result, err = s.confirmNolimitUpdateInfo()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func createNolimitCheckData(url string, params requestParameters) (*dataCreateNolimit, error) {
	uid := params.Uid

	psn, info, err := deCreateNolimit(url, uid, params.Data)
	if err != nil {
		return nil, err
	}

	return &dataCreateNolimit{
		Pg_conn: url,
		PsnInfo: *psn,
		Data:    *info,
	}, nil
}
func deCreateNolimit(url, uid, data string) (*tradePsnInfo, *reqNolimitCreate, error) {
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

	var req reqNolimitCreate
	err = json.Unmarshal([]byte(s), &req)
	if err != nil {
		return nil, nil, err
	}

	if strings.Compare(d.Token, req.Token) != 0 {
		return nil, nil, errors.New("check data failure.")
	}

	return &tradePsnInfo{
		Uid:         uid,
		Sys_pubkey:  d.Sys_pubkey,
		Sys_privkey: d.Sys_privkey,
		User_pubkey: d.Psn_pubkey,
	}, &req, nil
}
func (dc *dataCreateNolimit) createNolimitCheck() error {
	wid, err := db.GetWidFromUid(dc.Pg_conn, dc.PsnInfo.Uid)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return errors.New("user wallet does not exist.")
		} else {
			fmt.Println("database error:", err)
			return errors.New("database error, please notify the administrator.")
		}
	}

	value, err := strconv.ParseFloat(dc.Data.Value, 64)
	if err != nil {
		return err
	}
	if value <= 0 {
		return errors.New("trading value error.")
	}

	dc.PsnInfo.Wid = wid

	return nil
}
func (dc *dataCreateNolimit) createNolimitCreateInfo() (string, error) {
	timestamp := util.Timestamp()
	tid := createTradeID(dc.PsnInfo.Wid, timestamp)
	ttoken := calcNolimitTtokenHash(tid, dc.PsnInfo.Uid, dc.Data.Value, dc.Data.Cid, dc.Data.Expire, timestamp)

	err := db.CreateNolimitInfo(dc.Pg_conn, tid, dc.PsnInfo.Wid, dc.Data.Value, ttoken, dc.Data.Cid, dc.Data.Expire, timestamp)
	if err != nil {
		fmt.Println("database error:", err)
		return "", errors.New("database error, please notify the administrator.")
	}

	var resp respCreateNolimit
	resp.Tid, err = enCreateNolimitTid(tid, dc.PsnInfo.Sys_pubkey)
	if err != nil {
		return "", err
	}
	resp.Timestamp = util.Timestamp()

	b, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	veri, err := ver.ImportPubKey(dc.PsnInfo.User_pubkey)
	if err != nil {
		return "", err
	}
	pt, err := veri.Encrypt(string(b))
	if err != nil {
		return "", err
	}

	return pt, nil
}
func enCreateNolimitTid(sTid, sys_pubkey string) (string, error) {
	veri, err := ver.ImportPubKey(sys_pubkey)
	if err != nil {
		return "", err
	}
	pt, err := veri.Encrypt(sTid)
	if err != nil {
		return "", err
	}

	return pt, nil
}

func confirmNolimitCheckData(url string, params requestParameters) (*dataConfirmNolimit, error) {
	to_uid := params.Uid

	psn, from, err := deConfirmNolimit(url, to_uid, params.Data)
	if err != nil {
		return nil, err
	}

	tid, err := deConfirmNolimitTid(url, from.From, from.Info)
	if err != nil {
		return nil, err
	}

	return &dataConfirmNolimit{
		Pg_conn: url,
		PsnInfo: *psn,
		Data:    *from,
		Tid:     tid,
	}, nil
}
func deConfirmNolimit(url, uid, data string) (*tradePsnInfo, *reqNolimitConfirm, error) {
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

	var req reqNolimitConfirm
	err = json.Unmarshal([]byte(s), &req)
	if err != nil {
		return nil, nil, err
	}

	if strings.Compare(d.Token, req.Token) != 0 {
		return nil, nil, errors.New("check data failure.")
	}

	return &tradePsnInfo{
		Uid:         uid,
		Sys_pubkey:  d.Sys_pubkey,
		Sys_privkey: d.Sys_privkey,
		User_pubkey: d.Psn_pubkey,
	}, &req, nil
}
func deConfirmNolimitTid(url, uid, data string) (string, error) {
	if uid == "" || len(uid) != 40 {
		return "", errors.New("uid: non-regular value.")
	}

	d, err := db.GetInfoUserSys(url, uid)
	if err != nil {
		fmt.Println("database error:", err)
		return "", errors.New("database error, please notify the administrator.")
	}

	veri, err := ver.ImportPrivKey(d.Sys_privkey)
	if err != nil {
		return "", err
	}
	tid, err := veri.Decrypt(data)
	if err != nil {
		return "", err
	}

	return tid, nil
}
func (dc *dataConfirmNolimit) confirmNolimitUpdateInfo() (string, error) {
	d, err := db.ConfirmNolimitInfo(dc.Pg_conn, dc.Tid, dc.PsnInfo.Uid)
	if err != nil {
		fmt.Println("database error:", err)
		return "", errors.New("database error, please notify the administrator.")
	}

	var resp respConfirmNolimit
	resp.Cid = d.Cid
	resp.Balance = d.Balance
	resp.Increment = d.Increment
	resp.Timestamp = util.Timestamp()

	b, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	veri, err := ver.ImportPubKey(dc.PsnInfo.User_pubkey)
	if err != nil {
		return "", err
	}
	pt, err := veri.Encrypt(string(b))
	if err != nil {
		return "", err
	}

	return pt, nil
}

func calcNolimitTtokenHash(tid, from_uid, value string, cid int, expire, ctime int64) string {
	s := tid + "&" + from_uid + "&" + value + "&" + strconv.Itoa(cid) + "&" + strconv.FormatInt(expire, 10) + "&" + strconv.FormatInt(ctime, 10)
	return util.CreateSHA1Hash(s)
}
