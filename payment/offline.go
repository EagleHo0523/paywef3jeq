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

type reqOfflineCheckData struct {
	From      string `json:"from"`
	Info      string `json:"info"`
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
}
type reqOfflineCheckInfo struct {
	To        string `json:"to"`
	Cid       int    `json:"coin_idx"`
	Value     string `json:"value"`
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
}
type reqOfflineConfirm struct {
	Tid       string `json:"tid"`
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
}

type dataCheckOffline struct {
	Pg_conn string
	PsnInfo tradePsnInfo
	Data    reqOfflineCheckData
	Info    reqOfflineCheckInfo
}
type dataConfirmOffline struct {
	Pg_conn string
	Data    reqOfflineConfirm
	PsnInfo tradePsnInfo
}

type respCheckOffline struct {
	Tid         string `json:"tid"`
	From        string `json:"from"`
	To          string `json:"to"`
	Cid         int    `json:"coin_idx"`
	Value       string `json:"value"`
	Ttoken      string `json:"ttoken"`
	Create_time int64  `json:"create_time"`
	Message     string `json:"message"`
	Timestamp   int64  `json:"timestamp"`
}

func OfflineProcessRequest(r *http.Request) (*RequestService, error) {
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
func (rs *RequestService) OfflineMethod(url string) (interface{}, error) {
	var rtn interface{}
	var err error = nil

	start := time.Now()
	switch strings.ToUpper(rs.Method) {
	case "CHECK":
		rtn, err = processCheckOffline(url, rs.Params)
	case "CONFIRM":
		rtn, err = processConfirmOffline(url, rs.Params)
	default:
		err = errors.New("method required or not exist.")
		return rtn, err
	}
	log.Printf("Offline: %s\t%s", rs.Method, time.Since(start))

	if err != nil {
		return rtn, err
	}

	return rtn, nil
}

func processCheckOffline(url string, params requestParameters) (interface{}, error) {
	var resp responsePayment

	s, err := checkOfflineCheckData(url, params)
	if err != nil {
		return nil, err
	}

	resp.Result, err = s.checkOfflineCreateInfo()
	if err != nil {
		return nil, err
	}

	return resp, nil
}
func processConfirmOffline(url string, params requestParameters) (interface{}, error) {
	var resp responsePayment

	s, err := confirmOfflineCheckData(url, params)
	if err != nil {
		return nil, err
	}

	resp.Result, err = s.confirmOfflineUpdateInfo()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func checkOfflineCheckData(url string, params requestParameters) (*dataCheckOffline, error) {
	to_uid := params.Uid

	psn, data, err := deCheckOffline(url, to_uid, params.Data)
	if err != nil {
		return nil, err
	}

	if strings.Compare(to_uid, data.From) == 0 {
		return nil, errors.New("trade information error.")
	}

	d, err := deCheckOfflineReq(url, data.From, data.Info)
	if err != nil {
		return nil, err
	}

	if strings.Compare(to_uid, d.To) != 0 {
		return nil, errors.New("trade information does not match.")
	}

	return &dataCheckOffline{
		Pg_conn: url,
		PsnInfo: *psn,
		Data:    *data,
		Info:    *d,
	}, nil
}
func deCheckOffline(url, to_uid, data string) (*tradePsnInfo, *reqOfflineCheckData, error) {
	if to_uid == "" || len(to_uid) != 40 {
		return nil, nil, errors.New("uid: non-regular value.")
	}

	d, err := db.GetInfoUserSys(url, to_uid)
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

	var req reqOfflineCheckData
	err = json.Unmarshal([]byte(s), &req)
	if err != nil {
		return nil, nil, err
	}

	if strings.Compare(d.Token, req.Token) != 0 {
		return nil, nil, errors.New("check data failure.")
	}

	return &tradePsnInfo{
		Uid:         to_uid,
		Sys_pubkey:  d.Sys_pubkey,
		Sys_privkey: d.Sys_privkey,
		User_pubkey: d.Psn_pubkey,
	}, &req, nil

}
func deCheckOfflineReq(url, from_uid, info string) (*reqOfflineCheckInfo, error) {
	if from_uid == "" || len(from_uid) != 40 {
		return nil, errors.New("uid: non-regular value.")
	}

	d, err := db.GetInfoUserSys(url, from_uid)
	if err != nil {
		fmt.Println("database error:", err)
		return nil, errors.New("database error, please notify the administrator.")
	}

	veri, err := ver.ImportPrivKey(d.Sys_privkey)
	if err != nil {
		return nil, err
	}
	s, err := veri.Decrypt(info)
	if err != nil {
		return nil, err
	}

	var req reqOfflineCheckInfo
	err = json.Unmarshal([]byte(s), &req)
	if err != nil {
		return nil, err
	}

	if strings.Compare(d.Token, req.Token) != 0 {
		return nil, errors.New("check data failure.")
	}

	return &req, nil
}
func (dc *dataCheckOffline) checkOfflineCreateInfo() (string, error) {
	var resp respCheckOffline

	to_wid, err := db.GetWidFromUid(dc.Pg_conn, dc.PsnInfo.Uid)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return "", errors.New("user wallet does not exist.")
		} else {
			fmt.Println("database error:", err)
			return "", errors.New("database error, please notify the administrator.")
		}
	}

	from_wid, err := db.GetWidFromUid(dc.Pg_conn, dc.Data.From)
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
	ttoken := calcOfflineTtokenHash(tid, dc.Data.From, dc.PsnInfo.Uid, dc.Info.Value, dc.Info.Cid, timestamp)
	err = db.CreateOfflineInfo(dc.Pg_conn, tid, from_wid, to_wid, dc.Info.Value, ttoken, dc.Info.Cid, timestamp)
	resp.From = dc.Data.From
	resp.To = dc.PsnInfo.Uid
	resp.Cid = dc.Info.Cid
	resp.Value = dc.Info.Value
	resp.Cid = dc.Info.Cid
	resp.Ttoken = ttoken
	resp.Tid = tid
	resp.Message = "ok"
	if err != nil {
		if strings.Contains(err.Error(), "balance insufficient.") {
			resp.Message = err.Error()
			resp.Ttoken = ""
			resp.Tid = ""
		} else {
			fmt.Println("database error:", err)
			return "", errors.New("database error, please notify the administrator.")
		}
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

func confirmOfflineCheckData(url string, params requestParameters) (*dataConfirmOffline, error) {
	uid := params.Uid

	psn, info, err := deConfirmOffline(url, uid, params.Data)
	if err != nil {
		return nil, err
	}

	return &dataConfirmOffline{
		Pg_conn: url,
		PsnInfo: *psn,
		Data:    *info,
	}, nil
}
func deConfirmOffline(url, uid, data string) (*tradePsnInfo, *reqOfflineConfirm, error) {
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

	var req reqOfflineConfirm
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
func (dc *dataConfirmOffline) confirmOfflineUpdateInfo() (string, error) {
	d, err := db.ConfirmOfflineInfo(dc.Pg_conn, dc.Data.Tid)
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

func calcOfflineTtokenHash(tid, from, to, value string, cid int, ctime int64) string {
	s := tid + "&" + from + "&" + to + "&" + value + "&" + strconv.Itoa(cid) + "&" + strconv.FormatInt(ctime, 10)
	return util.CreateSHA1Hash(s)
}
