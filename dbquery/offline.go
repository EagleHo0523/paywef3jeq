package dbquery

import (
	"database/sql"
	"errors"

	util "../util"
	_ "github.com/go-sql-driver/mysql"
)

func CreateOfflineInfo(conn, tid, from_wid, to_wid, value, ttoken string, cid int, timestamp int64) error {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return err
	}
	defer db.Close()

	balance := ""
	sqlStr1 := "SELECT balance FROM wallet_balance WHERE wid=? AND coin_idx=?"
	err = db.QueryRow(sqlStr1, from_wid, cid).Scan(&balance)
	if err != nil {
		return err
	}
	b := checkValueBalance(value, balance)
	if !b {
		return errors.New("coin balance insufficient.")
	}

	sqlStr2 := "INSERT INTO offline_info (tid,tfrom,tto,coin_idx,tvalue,ttoken,create_time) VALUES (?,?,?,?,?,?,?)"
	_, err = db.Exec(sqlStr2, tid, from_wid, to_wid, cid, value, ttoken, timestamp)
	if err != nil {
		return err
	}

	return nil
}
func ConfirmOfflineInfo(conn, tid string) (*TradeBalance, error) {
	timestamp := util.Timestamp()

	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	balance, err := calcConfirmOffline(db, tid)
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
	_, err = tx.Exec(sqlStr3, tid, balance.To.Wid, balance.From.Wid, balance.Trade.Coin_idx, balance.Trade.Value, balance.Trade.Ttoken, 0, balance.Trade.Create_time, 0, timestamp, "O")
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	sqlStr4 := "DELETE FROM offline_info WHERE tid=?"
	_, err = tx.Exec(sqlStr4, tid)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &TradeBalance{
		Cid:     balance.Trade.Coin_idx,
		Balance: balance.To.Balance,
	}, nil
}

func calcConfirmOffline(db *sql.DB, tid string) (*balanceInfo, error) {
	var info TradeInfo
	info.Tid = tid

	sqlStr1 := "SELECT tfrom,tto,coin_idx,tvalue,ttoken,create_time FROM offline_info WHERE tid=?"
	err := db.QueryRow(sqlStr1, info.Tid).Scan(&info.From, &info.To, &info.Coin_idx, &info.Value, &info.Ttoken, &info.Create_time)
	if err != nil {
		return nil, err
	}

	to := userBalance{NoRow: false}
	to.Wid = info.To
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
