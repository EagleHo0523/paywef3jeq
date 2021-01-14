package register

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	db "../dbquery"
	hw "../hdwallet"
	util "../util"
	ver "../verification"
)

type RequestService struct {
	Method string            `json:"method,omitempty"`
	Params requestParameters `json:"params,omitempty"`
}
type requestParameters struct {
	Account string `json:"account,omitempty"`
	Data    string `json:"data,omitempty"`
	Uid     string `json:"uid,omitempty"`
}

type dataSetup struct {
	Pg_conn      string
	Pg_password  string
	Psn_account  string
	Psn_pubkey   string
	Psn_password string
	Psn_token    string
	Uid          string
	Mnemonic     string
	Sys_pubkey   string
	Sys_privkey  string
	Sys_token    string
}
type dataVerify struct {
	Pg_conn    string
	Uid        string
	VerifyCode string
	Token      string
	Account    string
	Psn_pubkey string
}
type dataLogin struct {
	Pg_conn    string
	Uid        string
	Token      string
	Account    string
	Psn_pubkey string
}

type reqSetup struct {
	PublicKey string `json:"public_key,omitempty"`
	Password  string `json:"password,omitempty"`
	Token     string `json:"token,omitempty"`
}
type reqVerify struct {
	Code  string `json:"code,omitempty"`
	Token string `json:"token,omitempty"`
}

type respSetup struct {
	Uid        string `json:"uid,omitempty"`
	Sys_pubkey string `json:"sys_pubkey,omitempty"`
	Token      string `json:"token,omitempty"`
	Timestamp  int64  `json:"timestamp,omitempty"`
}
type respVerify struct {
	Status    string `json:"status,omitempty"`
	Token     string `json:"token,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

type responseCreate struct {
	TempKey string `json:"key,omitempty"`
	Token   string `json:"token,omitempty"`
}
type responseSetup struct {
	Result string `json:"result,omitempty"`
}
type responseVerify struct {
	Result string `json:"result,omitempty"`
}
type responseLogin struct {
	Result string `json:"result,omitempty"`
}

func RegProcessRequest(r *http.Request) (*RequestService, error) {
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
func (rs *RequestService) ProcessMethod(url, pwd string) (interface{}, error) {
	var rtn interface{}
	var err error = nil

	start := time.Now()
	switch strings.ToUpper(rs.Method) {
	case "CREATE":
		rtn, err = processCreate(url, rs.Params)
	case "SETUP":
		rtn, err = processSetup(url, pwd, rs.Params)
	case "VERIFY":
		rtn, err = processVerify(url, rs.Params)
	default:
		err = errors.New("method required or not exist.")
		return rtn, err
	}
	log.Printf("Authorization: %s\t%s", rs.Method, time.Since(start))

	if err != nil {
		return rtn, err
	}

	return rtn, nil
}

func processCreate(url string, params requestParameters) (interface{}, error) {
	var respCreate responseCreate
	account := params.Account

	if account == "" || len(account) != 42 {
		return nil, errors.New("account: non-regular address.")
	}
	// 查db看看account是否已經存在
	err := db.CheckAccountRegTemp(url, account)
	if err != nil {
		fmt.Println("database error:", err)
		return nil, errors.New("database error, please notify the administrator.")
	}

	acc, err := ver.GenerateKeyPair() // create tmp_privkey/tmp_pubkey
	if err != nil {
		return nil, err
	}
	tmp_pubkey := acc.PublicKey()
	tmp_privkey := acc.PrivateKey()
	token := calcHashValue(account, tmp_pubkey)

	err = db.CreateRegTemp(url, account, tmp_pubkey, tmp_privkey, token)
	if err != nil {
		fmt.Println("database error:", err)
		return nil, errors.New("database error, please notify the administrator.")
	}

	respCreate.TempKey = tmp_pubkey
	respCreate.Token = token

	return respCreate, nil
}
func processSetup(url, pwd string, params requestParameters) (interface{}, error) {
	var respSetup responseSetup
	s, err := processCheckSetup(url, pwd, params)
	if err != nil {
		return nil, err
	}

	err = s.setupCreateAccount()
	if err != nil {
		return nil, err
	}

	respSetup.Result, err = s.setupCalcResp()
	if err != nil {
		return nil, err
	}

	return respSetup, nil
}
func processVerify(url string, params requestParameters) (interface{}, error) {
	var respVerify responseVerify

	v, err := processCheckVerify(url, params)
	if err != nil {
		return nil, err
	}

	err = v.verifyCheckCode()
	if err != nil {
		return nil, err
	}

	respVerify.Result, err = v.verifyRenewInfo()
	if err != nil {
		return nil, err
	}

	return respVerify, nil
}

// func processLogin(url string, params requestParameters) (interface{}, error) {
// 	var respLogin responseLogin

// 	d, err := processCheckLogin(url, params)
// 	if err != nil {
// 		return nil, err
// 	}

// 	s, err := d.loginRenewToken()
// 	if err != nil {
// 		return nil, err
// 	}
// 	respLogin.Result = s

// 	return respLogin, nil
// }

func processCheckSetup(url, pwd string, params requestParameters) (*dataSetup, error) {
	account := params.Account

	if account == "" || len(account) != 42 {
		return nil, errors.New("account: non-regular address.")
	}

	// 拿account從db取tmp_privkey
	d, err := db.GetSetupRegTemp(url, account)
	if err != nil {
		fmt.Println("database error:", err)
		return nil, errors.New("database error, please notify the administrator.")
	}

	veri, err := ver.ImportPrivKey(d.Private_key)
	if err != nil {
		return nil, err //errors.New("import key failure.")
	}
	s, err := veri.Decrypt(params.Data)
	if err != nil {
		return nil, err // errors.New("decrypt data failure.")
	}

	var data reqSetup
	err = json.Unmarshal([]byte(s), &data)
	if err != nil {
		return nil, err // errors.New("unmarshal data failure.")
	}

	if strings.Compare(d.Token, data.Token) != 0 {
		return nil, errors.New("check data failure.")
	}

	err = setupVerifyAccPub(account, data.PublicKey)
	if err != nil {
		return nil, err // errors.New("unmarshal data failure.")
	}

	return &dataSetup{
		Pg_conn:      url,
		Pg_password:  pwd,
		Psn_account:  account,
		Psn_pubkey:   data.PublicKey,
		Psn_password: calcHashValue(data.Password, "EagleHo"),
	}, nil
}
func setupVerifyAccPub(account, pubkey string) error {
	veri, err := ver.ImportPubKey(pubkey)
	if err != nil {
		return err
	}

	err = veri.VerifyAccountWithPubKey(account)
	if err != nil {
		return err
	}

	return nil
}
func (ds *dataSetup) setupCreateAccount() error {
	mnemonic, err := hw.NewMnemonic()
	if err != nil {
		return err
	}
	wallet, err := hw.NewFromMnemonic(mnemonic, ds.Pg_password+"EagleHo")
	if err != nil {
		return err
	}
	acc, err := wallet.GenerateAccount("ETH", 0)
	if err != nil {
		return err
	}
	// sys_addr := acc.Address()
	ds.Sys_pubkey = acc.PublicKey()
	ds.Sys_privkey = acc.PrivateKey()

	encode, err := util.EncryptKey(mnemonic, ds.Pg_password)
	if err != nil {
		return err
	}

	s := calcHashValue(mnemonic, ds.Psn_account)
	ds.Uid = calcHashValue(s, "EagleHo")
	ds.Sys_token = calcHashValue(ds.Uid, ds.Psn_account)

	err = db.CreateUserSys(ds.Pg_conn, ds.Psn_account, ds.Psn_pubkey, ds.Psn_password, ds.Uid, encode, ds.Pg_password, ds.Sys_pubkey, ds.Sys_privkey)
	if err != nil {
		fmt.Println("database error:", err)
		return errors.New("database error, please notify the administrator.")
	}

	return nil
}
func (ds *dataSetup) setupCalcResp() (string, error) {
	resp := respSetup{
		Uid:        ds.Uid,
		Sys_pubkey: ds.Sys_pubkey,
		Token:      ds.Sys_token,
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	veri, err := ver.ImportPubKey(ds.Psn_pubkey)
	if err != nil {
		return "", err
	}
	pt, err := veri.Encrypt(string(b))
	if err != nil {
		return "", err
	}

	return pt, nil
}

func processCheckVerify(url string, params requestParameters) (*dataVerify, error) {
	uid := params.Uid

	if uid == "" || len(uid) != 40 {
		return nil, errors.New("uid: non-regular format.")
	}

	d, err := db.GetVerifyUserReg(url, uid)
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

	var data reqVerify
	err = json.Unmarshal([]byte(s), &data)
	if err != nil {
		return nil, err // errors.New("unmarshal data failure.")
	}

	return &dataVerify{
		Pg_conn:    url,
		Uid:        uid,
		VerifyCode: data.Code,
		Token:      data.Token,
		Account:    d.Account,
		Psn_pubkey: d.Psn_pubkey,
	}, nil
}
func (dv *dataVerify) verifyCheckCode() error {
	code := calcHashValue(dv.Account, dv.Psn_pubkey)
	if strings.Compare(code, dv.VerifyCode) != 0 {
		return errors.New("check data failure.")
	}
	token := calcHashValue(dv.Uid, dv.Account)
	if strings.Compare(token, dv.Token) != 0 {
		return errors.New("check data failure.")
	}
	return nil
}
func (dv *dataVerify) verifyRenewInfo() (string, error) {
	timestamp := util.Timestamp()
	sys_token := calcSysToken(dv.Uid, dv.Account, timestamp)

	err := db.RenewRegInfo(dv.Pg_conn, dv.Uid, sys_token, timestamp)
	if err != nil {
		fmt.Println("database error:", err)
		return "", errors.New("database error, please notify the administrator.")
	}

	var resp respVerify
	resp.Status = "success"
	resp.Token = sys_token

	b, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	veri, err := ver.ImportPubKey(dv.Psn_pubkey)
	if err != nil {
		return "", err
	}
	pt, err := veri.Encrypt(string(b))
	if err != nil {
		return "", err
	}

	return pt, nil
}

// func processCheckLogin(url string, params requestParameters) (*dataLogin, error) {
// 	uid := params.Uid

// 	d, err := db.GetInfoUserReg(url, uid)
// 	if err != nil {
// 		return nil, errors.New("database error, please notify the administrator.")
// 	}

// 	veri, err := ver.ImportPrivKey(d.Sys_privkey)
// 	if err != nil {
// 		return nil, err
// 	}
// 	s, err := veri.Decrypt(params.Data)
// 	if err != nil {
// 		return nil, err // errors.New("decrypt data failure.")
// 	}

// 	if strings.Compare(s, d.Account) != 0 {
// 		return nil, errors.New("check data failure.")
// 	}

// 	return &dataLogin{
// 		Pg_conn:    url,
// 		Uid:        uid,
// 		Account:    d.Account,
// 		Psn_pubkey: d.Psn_pubkey,
// 	}, nil
// }
// func (dl *dataLogin) loginRenewToken() (string, error) {
// 	timestamp := time.Now().UTC().Unix()
// 	sys_token := calcSysToken(dl.Uid, dl.Account, timestamp)

// 	err := db.RenewTokenUserAuth(dl.Pg_conn, dl.Uid, sys_token, timestamp)
// 	if err != nil {
// 		return "", errors.New("database error, please notify the administrator.")
// 	}

// 	veri, err := ver.ImportPubKey(dl.Psn_pubkey)
// 	if err != nil {
// 		return "", err
// 	}

// 	s, err := veri.Encrypt(sys_token)
// 	if err != nil {
// 		return "", err
// 	}

// 	return s, nil
// }

func calcHashValue(text1, text2 string) string {
	s := sha1.Sum([]byte(text1 + "&" + text2))
	return hex.EncodeToString(s[:])
}
func calcSysToken(uid, account string, timestamp int64) string {
	s := calcHashValue(uid, account)
	return calcHashValue(s, string(timestamp))
}
