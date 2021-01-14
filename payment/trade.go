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

type tradePsnInfo struct {
	Uid         string
	Token       string
	Wid         string
	User_pubkey string
	Sys_pubkey  string
	Sys_privkey string
}

type reqTradeCreate struct {
	To        string `json:"to"`
	Cid       int    `json:"coin_idx"`
	Value     string `json:"value"`
	Expire    int64  `json:"expire,omitempty"`
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
}
type reqTradeCheck struct {
	To        string `json:"to"`
	Info      string `json:"info"`
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
}
type reqTradeCheckTid struct {
	Tid string `json:"tid"`
}
type reqTradeConfirm struct {
	Tid       string `json:"tid"`
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
}

type respCreateTrading struct {
	Tid       string `json:"tid"`
	Timestamp int64  `json:"timestamp"`
}
type respCheckTrading struct {
	Tid         string `json:"tid"`
	To          string `json:"to"`
	Cid         int    `json:"coin_idx"`
	Value       string `json:"value"`
	Ttoken      string `json:"ttoken"`
	Create_time int64  `json:"create_time"`
	Expire      int64  `json:"expire,omitempty"`
	Message     string `json:"message"`
	Timestamp   int64  `json:"timestamp"`
}
type respConfirmTrading struct {
	Cid       int    `json:"coin_idx"`
	Balance   string `json:"balance"`
	Timestamp int64  `json:"timestamp"`
}

type dataCreateTrading struct {
	Data    reqTradeCreate
	Pg_conn string
	PsnInfo tradePsnInfo
}
type dataCheckTrading struct {
	Data    reqTradeCheck
	Tid     string
	Pg_conn string
	PsnInfo tradePsnInfo
}
type dataConfirmTrading struct {
	Pg_conn string
	Data    reqTradeConfirm
	PsnInfo tradePsnInfo
}

func TradeProcessRequest(r *http.Request) (*RequestService, error) {
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
func (rs *RequestService) TradeMethod(url string) (interface{}, error) {
	var rtn interface{}
	var err error = nil

	start := time.Now()
	switch strings.ToUpper(rs.Method) {
	case "CREATE":
		rtn, err = processCreateTrading(url, rs.Params)
	case "CHECK":
		rtn, err = processCheckTrading(url, rs.Params)
	case "CONFIRM":
		rtn, err = processConfirmTrading(url, rs.Params)
	default:
		err = errors.New("method required or not exist.")
		return rtn, err
	}
	log.Printf("Trade: %s\t%s", rs.Method, time.Since(start))

	if err != nil {
		return rtn, err
	}

	return rtn, nil
}

func processCreateTrading(url string, params requestParameters) (interface{}, error) {
	var resp responsePayment

	s, err := createTradingCheckData(url, params)
	if err != nil {
		return nil, err
	}

	err = s.createTradingGet()
	if err != nil {
		return nil, err
	}

	resp.Result, err = s.createTradingCreateInfo()
	if err != nil {
		return nil, err
	}

	return resp, nil
}
func processCheckTrading(url string, params requestParameters) (interface{}, error) {
	var resp responsePayment

	s, err := checkTradingCheckData(url, params)
	if err != nil {
		return nil, err
	}

	resp.Result, err = s.checkTradingUpdateInfo()
	if err != nil {
		return nil, err
	}

	return resp, nil
}
func processConfirmTrading(url string, params requestParameters) (interface{}, error) {
	var resp responsePayment

	s, err := confirmTradingCheckData(url, params)
	if err != nil {
		return nil, err
	}

	resp.Result, err = s.confirmTradingUpdateInfo()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func createTradingCheckData(url string, params requestParameters) (*dataCreateTrading, error) {
	uid := params.Uid

	if uid == "" || len(uid) != 40 {
		return nil, errors.New("uid: non-regular value.")
	}

	d, err := db.GetInfoUserSys(url, uid)
	if err != nil {
		fmt.Println("database error:", err)
		return nil, errors.New("database error, please notify the administrator.")
	}

	veri, err := ver.ImportPrivKey(d.Sys_privkey)
	if err != nil {
		return nil, err //errors.New("import key failure.")
	}
	s, err := veri.Decrypt(params.Data)
	if err != nil {
		return nil, err // errors.New("decrypt data failure.")
	}

	var data reqTradeCreate
	err = json.Unmarshal([]byte(s), &data)
	if err != nil {
		return nil, err // errors.New("unmarshal data failure.")
	}

	if strings.Compare(d.Token, data.Token) != 0 {
		return nil, errors.New("check data failure.")
	}

	if strings.Compare(uid, data.To) != 0 {
		return nil, errors.New("trading information error.")
	}

	return &dataCreateTrading{
		Pg_conn: url,
		Data:    data,
		PsnInfo: tradePsnInfo{
			Uid:         uid,
			Sys_pubkey:  d.Sys_pubkey,
			Sys_privkey: d.Sys_privkey,
			User_pubkey: d.Psn_pubkey,
		},
	}, nil
}
func (dc *dataCreateTrading) createTradingGet() error {
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
func (dc *dataCreateTrading) createTradingCreateInfo() (string, error) {
	timestamp := util.Timestamp()
	tid := createTradeID(dc.PsnInfo.Wid, timestamp)
	ttoken := calcTradeTtokenHash(tid, dc.Data.To, dc.Data.Value, dc.Data.Cid, dc.Data.Expire, timestamp)

	err := db.CreateTradeInfo(dc.Pg_conn, tid, dc.PsnInfo.Wid, dc.Data.Value, ttoken, dc.Data.Cid, dc.Data.Expire, timestamp)
	if err != nil {
		fmt.Println("database error:", err)
		return "", errors.New("database error, please notify the administrator.")
	}

	var resp respCreateTrading
	resp.Tid, err = enCreateTradingTid(tid, dc.PsnInfo.Sys_pubkey)
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
func enCreateTradingTid(sTid, sys_pubkey string) (string, error) {
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

func checkTradingCheckData(url string, params requestParameters) (*dataCheckTrading, error) {
	uid := params.Uid

	psn, from, err := deCheckTrading(url, uid, params.Data)
	if err != nil {
		return nil, err
	}

	tid, err := deCheckTradingTid(url, from.To, from.Info)
	if err != nil {
		return nil, err
	}

	return &dataCheckTrading{
		Pg_conn: url,
		PsnInfo: *psn,
		Data:    *from,
		Tid:     tid,
	}, nil
}
func deCheckTrading(url, uid, data string) (*tradePsnInfo, *reqTradeCheck, error) {
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

	var req reqTradeCheck
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
func deCheckTradingTid(url, uid, data string) (string, error) {
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
func (dc *dataCheckTrading) checkTradingUpdateInfo() (string, error) {
	var resp respCheckTrading

	timestamp := util.Timestamp()
	info, err := db.UpdateTradeInfo(dc.Pg_conn, dc.Tid, dc.PsnInfo.Uid, timestamp)
	if err != nil {
		fmt.Println("database error:", err)
		return "", errors.New("database error, please notify the administrator.")
	}
	resp.Tid = info.Tid
	resp.To = info.To
	resp.Cid = info.Coin_idx
	resp.Value = info.Value
	resp.Ttoken = info.Ttoken
	resp.Expire = info.Expire
	resp.Create_time = info.Create_time
	resp.Message = "ok"
	if err != nil && strings.Contains(err.Error(), "balance insufficient.") {
		resp.Message = err.Error()
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

func confirmTradingCheckData(url string, params requestParameters) (*dataConfirmTrading, error) {
	uid := params.Uid

	psn, info, err := deConfirmTrading(url, uid, params.Data)
	if err != nil {
		return nil, err
	}

	return &dataConfirmTrading{
		Pg_conn: url,
		PsnInfo: *psn,
		Data:    *info,
	}, nil
}
func deConfirmTrading(url, uid, data string) (*tradePsnInfo, *reqTradeConfirm, error) {
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

	var req reqTradeConfirm
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
func (dc *dataConfirmTrading) confirmTradingUpdateInfo() (string, error) {
	d, err := db.ConfirmTradeInfo(dc.Pg_conn, dc.Data.Tid)
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

func calcTradeTtokenHash(tid, to, value string, cid int, expire, ctime int64) string {
	s := tid + "&" + to + "&" + value + "&" + strconv.Itoa(cid) + "&" + strconv.FormatInt(expire, 10) + "&" + strconv.FormatInt(ctime, 10)
	return util.CreateSHA1Hash(s)
}
func createTradeID(to string, timestamp int64) string {
	s := "EagleHo&" + to + "&" + strconv.FormatInt(timestamp, 10)
	return util.CreateSHA1Hash(s)
}
