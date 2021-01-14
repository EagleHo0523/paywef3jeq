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

type reqCreateWallet struct {
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
}
type reqBalanceWallet struct {
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
}
type reqHistoryWallet struct {
	Coin_idx   int    `json:"coin_idx,omitempty"`
	Start_time int64  `json:"start_time,omitempty"`
	End_time   int64  `json:"end_time,omitempty"`
	Token      string `json:"token"`
	Timestamp  int64  `json:"timestamp"`
}

type dataCreateWallet struct {
	Pg_conn     string
	Psn_account string
	Psn_pubkey  string
	Uid         string
	Sys_privkey string
	Sys_token   string
}
type dataBalanceWallet struct {
	Pg_conn string
	Data    reqBalanceWallet
	PsnInfo tradePsnInfo
	Wid     string
}
type dataHistoryWallet struct {
	Pg_conn string
	Data    reqHistoryWallet
	PsnInfo tradePsnInfo
}

type coinInfo struct {
	Cid         int    `json:"coin_idx"`
	Sname       string `json:"sname"`
	Fname       string `json:"fname"`
	Description string `json:"description"`
}
type balanceInfo struct {
	Cid     int    `json:"coin_idx"`
	Balance string `json:"balance"`
}
type historyInfo struct {
	Tid         string `json:"tid"`
	From        string `json:"from"`
	To          string `json:"to"`
	Coin_idx    int    `json:"coin_idx"`
	Value       string `json:"value"`
	Type        string `json:"type"`
	Create_time int64  `json:"time"`
}

func PayProcessRequest(r *http.Request) (*RequestService, error) {
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
func (rs *RequestService) PayMethod(url string) (interface{}, error) {
	var rtn interface{}
	var err error = nil

	start := time.Now()
	switch strings.ToUpper(rs.Method) {
	case "CREATE":
		rtn, err = processCreateWallet(url, rs.Params)
	case "COINBASE":
		rtn, err = processCoinbaseWallet(url, rs.Params)
	case "BALANCE":
		rtn, err = processBalanceWallet(url, rs.Params)
	case "HISTORY":
		rtn, err = processHistoryWallet(url, rs.Params)
	default:
		err = errors.New("method required or not exist.")
		return rtn, err
	}
	log.Printf("Wallet: %s\t%s", rs.Method, time.Since(start))

	if err != nil {
		return rtn, err
	}

	return rtn, nil
}

func processCreateWallet(url string, params requestParameters) (interface{}, error) {
	var respCreate responsePayment

	s, err := createWalletCreate(url, params)
	if err != nil {
		return nil, err
	}

	respCreate.Result, err = s.createWalletcheck()
	if err != nil {
		return nil, err
	}

	return respCreate, nil
}
func processCoinbaseWallet(url string, params requestParameters) (interface{}, error) {
	rtn, err := coinbaseWalletGetInfo(url, params)
	if err != nil {
		return nil, err
	}

	return rtn, nil
}
func processBalanceWallet(url string, params requestParameters) (interface{}, error) {
	var respCreate responsePayment

	s, err := balanceWalletCheckData(url, params)
	if err != nil {
		return nil, err
	}

	respCreate.Result, err = s.balanceWalletGetInfo()
	if err != nil {
		return nil, err
	}

	return respCreate, nil
}
func processHistoryWallet(url string, params requestParameters) (interface{}, error) {
	var respCreate responsePayment

	s, err := historyWalletCheckData(url, params)
	if err != nil {
		return nil, err
	}

	respCreate.Result, err = s.historyWalletGetInfo()
	if err != nil {
		return nil, err
	}

	return respCreate, nil
}

func createWalletCreate(url string, params requestParameters) (*dataCreateWallet, error) {
	uid := params.Uid

	if uid == "" || len(uid) != 40 {
		return nil, errors.New("uid: non-regular value.")
	}

	w, err := db.GetWidFromUid(url, uid)
	if err != nil && !strings.Contains(err.Error(), "no rows") {
		fmt.Println("database error:", err)
		return nil, errors.New("database error, please notify the administrator.")
	}
	if w != "" {
		return nil, errors.New("user wallet exist.")
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

	var data reqCreateWallet
	err = json.Unmarshal([]byte(s), &data)
	if err != nil {
		return nil, err // errors.New("unmarshal data failure.")
	}

	if strings.Compare(d.Token, data.Token) != 0 {
		return nil, errors.New("check data failure.")
	}

	timestamp := util.Timestamp()
	if data.Timestamp > timestamp {
		return nil, errors.New("check time error.")
	}

	return &dataCreateWallet{
		Pg_conn:     url,
		Psn_account: d.Account,
		Psn_pubkey:  d.Psn_pubkey,
		Uid:         uid,
		Sys_privkey: d.Sys_privkey,
		Sys_token:   d.Token,
	}, nil
}
func (dl *dataCreateWallet) createWalletcheck() (string, error) {
	timestamp := util.Timestamp()

	s := util.CreateSHA1Hash(dl.Psn_account + "&" + dl.Uid)
	wid := util.CreateSHA1Hash(s + "&" + strconv.FormatInt(timestamp, 10))

	err := db.CreateUserWallet(dl.Pg_conn, dl.Uid, wid, timestamp)
	if err != nil {
		fmt.Println("database error:", err)
		return "", errors.New("database error, please notify the administrator.")
	}

	return "success", nil
}

func coinbaseWalletGetInfo(url string, params requestParameters) (interface{}, error) {
	d, err := db.CoinbaseWalletGetInfo(url)
	if err != nil {
		fmt.Println("database error:", err)
		return nil, errors.New("database error, please notify the administrator.")
	}

	var coins []coinInfo
	for i := 0; i < len(d); i++ {
		var c coinInfo
		c.Cid = d[i].Idx
		c.Sname = d[i].Sname
		c.Fname = d[i].Fname
		if d[i].Note != nil {
			c.Description = fmt.Sprintf("%s", d[i].Note)
		} else {
			c.Description = ""
		}
		coins = append(coins, c)
	}

	return coins, nil
}

func balanceWalletCheckData(url string, params requestParameters) (*dataBalanceWallet, error) {
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

	var data reqBalanceWallet
	err = json.Unmarshal([]byte(s), &data)
	if err != nil {
		return nil, err // errors.New("unmarshal data failure.")
	}

	if strings.Compare(d.Token, data.Token) != 0 {
		return nil, errors.New("check data failure.")
	}

	return &dataBalanceWallet{
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
func (dw *dataBalanceWallet) balanceWalletGetInfo() (string, error) {
	d, err := db.BalanceWalletGetInfo(dw.Pg_conn, dw.PsnInfo.Uid)
	if err != nil {
		fmt.Println("database error:", err)
		return "", errors.New("database error, please notify the administrator.")
	}

	var balances []balanceInfo
	for i := 0; i < len(d); i++ {
		var b balanceInfo
		b.Cid = d[i].Cid
		b.Balance = d[i].Balance
		balances = append(balances, b)
	}

	b, err := json.Marshal(balances)
	if err != nil {
		return "", err
	}

	veri, err := ver.ImportPubKey(dw.PsnInfo.User_pubkey)
	if err != nil {
		return "", err
	}
	pt, err := veri.Encrypt(string(b))
	if err != nil {
		return "", err
	}

	return pt, nil
}

func historyWalletCheckData(url string, params requestParameters) (*dataHistoryWallet, error) {
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

	var data reqHistoryWallet
	err = json.Unmarshal([]byte(s), &data)
	if err != nil {
		return nil, err // errors.New("unmarshal data failure.")
	}

	if strings.Compare(d.Token, data.Token) != 0 {
		return nil, errors.New("check data failure.")
	}

	return &dataHistoryWallet{
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
func (dh *dataHistoryWallet) historyWalletGetInfo() (string, error) {
	d, err := db.HistoryWalletGetInfo(dh.Pg_conn, dh.PsnInfo.Uid, dh.Data.Coin_idx)
	if err != nil {
		fmt.Println("database error:", err)
		return "", errors.New("database error, please notify the administrator.")
	}

	var histories []historyInfo
	for i := 0; i < len(d); i++ {
		var h historyInfo
		h.Tid = d[i].Tid
		h.From = d[i].From
		h.To = d[i].To
		h.Coin_idx = d[i].Coin_idx
		h.Value = d[i].Value
		h.Create_time = d[i].Create_time
		switch d[i].Ttype {
		case "T":
			h.Type = "TRADE"
		case "F":
			h.Type = "TRANSFER"
		case "O":
			h.Type = "OFFLINE"
		case "N":
			h.Type = "NOLIMIT"
		}
		histories = append(histories, h)
	}

	b, err := json.Marshal(histories)
	if err != nil {
		return "", err
	}

	veri, err := ver.ImportPubKey(dh.PsnInfo.User_pubkey)
	if err != nil {
		return "", err
	}
	pt, err := veri.Encrypt(string(b))
	if err != nil {
		return "", err
	}

	return pt, nil
}
