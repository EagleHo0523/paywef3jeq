package dbquery

import (
	"database/sql"
	"strconv"

	util "../util"
	_ "github.com/go-sql-driver/mysql"
)

type TradeInfo struct {
	Tid         string
	To          string
	From        string
	Coin_idx    int
	Value       string
	Ttoken      string
	Expire      int64
	Create_time int64
	Update_time int64
}
type TradeBalance struct {
	Cid     int
	Balance string
}

type userBalance struct {
	Wid     string
	Balance string
	NoRow   bool
}
type balanceInfo struct {
	From  userBalance
	To    userBalance
	Trade TradeInfo
}

func GetInfoUserSys(conn, uid string) (*InfoUserSys, error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// var count int64
	// sqlStr1 := "SELECT COUNT(*) FROM user_wallet WHERE uid='" + uid + "'"
	// err = db.QueryRow(sqlStr1).Scan(&count)
	// if err != nil {
	// 	return nil, err
	// }
	// if count > 0 {
	// 	return nil, errors.New("user wallet exist.")
	// }

	var info InfoUserSys
	sqlStr := "SELECT reg.account,reg.public_key,sys.public_key,sys.private_key,auth.token FROM user_sys AS sys, (SELECT * FROM user_reg WHERE status='N' AND uid=?) AS reg, (SELECT * FROM user_auth WHERE uid=?) AS auth WHERE sys.uid=reg.uid AND sys.uid=auth.uid"
	row := db.QueryRow(sqlStr, uid, uid)
	err = row.Scan(&info.Account, &info.Psn_pubkey, &info.Sys_pubkey, &info.Sys_privkey, &info.Token)
	if err != nil {
		return nil, err
	}

	return &info, nil
}
func GetWidFromUid(conn, uid string) (string, error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return "", err
	}
	defer db.Close()

	wid := ""
	sqlStr := "SELECT wid FROM user_wallet WHERE uid=?"
	row := db.QueryRow(sqlStr, uid)
	err = row.Scan(&wid)
	if err != nil {
		return "", err
	}

	return wid, nil
}

func checkValueBalance(v, b string) bool {
	rtn := false

	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return rtn
	}

	balance := 0.0
	if b != "" {
		balance, err = strconv.ParseFloat(b, 64)
		if err != nil {
			return rtn
		}
	}

	if balance >= value {
		rtn = true
	}

	return rtn
}
func calcBalance(value, toValue, fromValue string) (string, string) {
	v := util.ToWei(value, 18)
	t := util.ToWei(toValue, 18)
	f := util.ToWei(fromValue, 18)
	t.Add(t, v)
	f.Sub(f, v)
	return util.ToDecimal(t, 18).String(), util.ToDecimal(f, 18).String()
}
