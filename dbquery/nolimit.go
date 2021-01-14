package dbquery

import (
	"database/sql"
	"errors"
	"strings"

	util "../util"
	_ "github.com/go-sql-driver/mysql"
)

type NolimitBalance struct {
	Cid       int
	Increment string
	Balance   string
}

func CreateNolimitInfo(conn, tid, from_wid, value, ttoken string, cid int, expire, timestamp int64) error {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStr := "INSERT INTO nolimit_info (tid,tfrom,coin_idx,tvalue,ttoken,expire_time,create_time) values (?,?,?,?,?,?,?)"
	_, err = db.Exec(sqlStr, tid, from_wid, cid, value, ttoken, expire, timestamp)
	if err != nil {
		return err
	}

	return nil
}
func ConfirmNolimitInfo(conn, tid, to_uid string) (*NolimitBalance, error) {
	timestamp := util.Timestamp()

	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	balance, err := calcConfirmNolimit(db, tid, to_uid)
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	sqlStr1 := ""
	if balance.To.NoRow {
		sqlStr1 = "INSERT INTO wallet_balance (balance,update_time,wid,coin_idx) VALUES (?,?,?,?)"
	} else {
		sqlStr1 = "UPDATE wallet_balance SET balance=?, update_time=? WHERE wid=? AND coin_idx=?"
	}
	r, err := tx.Exec(sqlStr1, balance.To.Balance, timestamp, balance.To.Wid, balance.Trade.Coin_idx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	count, _ := r.RowsAffected()
	if count != 1 {
		tx.Rollback()
		return nil, errors.New("update balance error.")
	}

	sqlStr2 := "UPDATE wallet_balance SET balance=?, update_time=? WHERE wid=? AND coin_idx=?"
	r, err = tx.Exec(sqlStr2, balance.From.Balance, timestamp, balance.From.Wid, balance.Trade.Coin_idx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	count, _ = r.RowsAffected()
	if count != 1 {
		tx.Rollback()
		return nil, errors.New("update balance error.")
	}

	sqlStr3 := "INSERT INTO trans_list (tid,tto,tfrom,coin_idx,tvalue,ttoken,info_expire_time,info_create_time,info_update_time,create_time,ttype) VALUES (?,?,?,?,?,?,?,?,?,?,?)"
	_, err = tx.Exec(sqlStr3, tid, balance.To.Wid, balance.From.Wid, balance.Trade.Coin_idx, balance.Trade.Value, balance.Trade.Ttoken, balance.Trade.Expire, balance.Trade.Create_time, 0, timestamp, "N")
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	sqlStr4 := "DELETE FROM nolimit_info WHERE tid=?"
	_, err = tx.Exec(sqlStr4, tid)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &NolimitBalance{
		Cid:       balance.Trade.Coin_idx,
		Increment: balance.Trade.Value,
		Balance:   balance.To.Balance,
	}, nil
}

func calcConfirmNolimit(db *sql.DB, tid, to_uid string) (*balanceInfo, error) {
	var info TradeInfo
	info.Tid = tid

	sqlStr1 := "SELECT tfrom,coin_idx,tvalue,ttoken,expire_time,create_time FROM nolimit_info WHERE tid=?"
	err := db.QueryRow(sqlStr1, info.Tid).Scan(&info.From, &info.Coin_idx, &info.Value, &info.Ttoken, &info.Expire, &info.Create_time)
	if err != nil {
		return nil, err
	}

	to := userBalance{NoRow: false}
	sqlStr := "SELECT wid FROM user_wallet WHERE uid=?"
	err = db.QueryRow(sqlStr, to_uid).Scan(&to.Wid)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, errors.New("user wallet does not exist.")
		} else {
			return nil, err
		}
	}

	sqlStr2 := "SELECT balance FROM wallet_balance WHERE wid=? AND coin_idx=?"
	err = db.QueryRow(sqlStr2, to.Wid, info.Coin_idx).Scan(&to.Balance)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if to.Balance == "" {
		to.NoRow = true
	}

	var from userBalance
	from.Wid = info.From
	sqlStr3 := "SELECT balance FROM wallet_balance WHERE wid=? AND coin_idx=?"
	err = db.QueryRow(sqlStr3, from.Wid, info.Coin_idx).Scan(&from.Balance)
	if err != nil {
		return nil, err
	}

	to.Balance, from.Balance = calcBalance(info.Value, to.Balance, from.Balance)

	return &balanceInfo{
		From:  from,
		To:    to,
		Trade: info,
	}, nil
}
