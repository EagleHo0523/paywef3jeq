package dbquery

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type CoinInfo struct {
	Idx   int
	Sname string
	Fname string
	Note  interface{}
}
type BalanceInfo struct {
	Cid     int
	Balance string
}
type HistoryInfo struct {
	Tid         string
	From        string
	To          string
	Value       string
	Coin_idx    int
	Ttype       string
	Create_time int64
}

func CreateUserWallet(conn, uid, wid string, timestamp int64) error {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStr2 := "INSERT INTO user_wallet (uid, wid, create_time) VALUES (?, ?, ?)"
	_, err = db.Exec(sqlStr2, uid, wid, timestamp)
	if err != nil {
		return err
	}

	return nil
}
func CoinbaseWalletGetInfo(conn string) ([]CoinInfo, error) {
	var coinbases []CoinInfo

	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlStr := "SELECT idx,sname,fname,note FROM coin_info WHERE utype='N'"
	rows, err := db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c CoinInfo
		if err := rows.Scan(&c.Idx, &c.Sname, &c.Fname, &c.Note); err != nil {
			return nil, err
		}
		coinbases = append(coinbases, c)
	}
	return coinbases, nil
}
func BalanceWalletGetInfo(conn, uid string) ([]BalanceInfo, error) {
	var balances []BalanceInfo

	db, err := sql.Open("mysql", conn)
	if err != nil {
		return balances, err
	}
	defer db.Close()

	sqlStr := "SELECT b.coin_idx,b.balance FROM wallet_balance as b, user_wallet as u WHERE b.wid=u.wid AND u.uid=?"
	rows, err := db.Query(sqlStr, uid)
	if err != nil {
		return balances, err
	}
	defer rows.Close()

	for rows.Next() {
		var b BalanceInfo
		if err := rows.Scan(&b.Cid, &b.Balance); err != nil {
			return balances, err
		}
		balances = append(balances, b)
	}

	return balances, nil
}
func HistoryWalletGetInfo(conn, uid string, cid int) ([]HistoryInfo, error) {
	var histories []HistoryInfo

	db, err := sql.Open("mysql", conn)
	if err != nil {
		return histories, err
	}
	defer db.Close()

	sqlStr := ""
	var rows *sql.Rows
	if cid == 0 {
		sqlStr = "SELECT ts.tid,fu.uid as from_uid,tu.uid as to_uid,ts.coin_idx,ts.tvalue,ts.ttype,ts.create_time from trans_list as ts,(select * from user_wallet where uid=?) as u,user_wallet as tu,user_wallet as fu where (ts.tto=u.wid or ts.tfrom=u.wid) and ts.tto=tu.wid and ts.tfrom=fu.wid order by ts.create_time desc limit 50"
		rows, err = db.Query(sqlStr, uid)
	} else {
		sqlStr = "SELECT ts.tid,fu.uid as from_uid,tu.uid as to_uid,ts.coin_idx,ts.tvalue,ts.ttype,ts.create_time from trans_list as ts,(select * from user_wallet where uid=?) as u,user_wallet as tu,user_wallet as fu where (ts.tto=u.wid or ts.tfrom=u.wid) and coin_idx=? and ts.tto=tu.wid and ts.tfrom=fu.wid order by ts.create_time desc limit 50"
		rows, err = db.Query(sqlStr, uid, cid)
	}

	if err != nil {
		return histories, err
	}
	defer rows.Close()

	for rows.Next() {
		var h HistoryInfo
		if err := rows.Scan(&h.Tid, &h.From, &h.To, &h.Coin_idx, &h.Value, &h.Ttype, &h.Create_time); err != nil {
			return histories, err
		}
		histories = append(histories, h)
	}

	return histories, nil
}
