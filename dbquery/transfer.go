package dbquery

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

func ConfirmTransferInfo(conn, tid, from_wid, to_wid, value, ttoken string, cid int, timestamp int64) (*TradeBalance, error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	balance, err := calcConfirmTransfer(db, from_wid, to_wid, value, cid)
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
	r, err := tx.Exec(sqlStr1, balance.To.Balance, timestamp, to_wid, cid)
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
	r, err = tx.Exec(sqlStr2, balance.From.Balance, timestamp, from_wid, cid)
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
	_, err = tx.Exec(sqlStr3, tid, to_wid, from_wid, cid, value, ttoken, 0, timestamp, timestamp, timestamp, "F")
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &TradeBalance{
		Cid:     cid,
		Balance: balance.From.Balance,
	}, nil
}

func calcConfirmTransfer(db *sql.DB, from_wid, to_wid, value string, cid int) (*balanceInfo, error) {
	var info TradeInfo
	info.Value = value

	var from userBalance
	sqlStr2 := "SELECT balance FROM wallet_balance WHERE wid=? AND coin_idx=?"
	err := db.QueryRow(sqlStr2, from_wid, cid).Scan(&from.Balance)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if from.Balance == "" {
		from.NoRow = true
	}
	b := checkValueBalance(info.Value, from.Balance)
	if !b {
		return nil, errors.New("coin balance insufficient.")
	}

	to := userBalance{NoRow: false}
	sqlStr3 := "SELECT balance FROM wallet_balance WHERE wid=? AND coin_idx=?"
	err = db.QueryRow(sqlStr3, to_wid, cid).Scan(&to.Balance)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if to.Balance == "" {
		to.NoRow = true
	}

	to.Balance, from.Balance = calcBalance(info.Value, to.Balance, from.Balance)

	return &balanceInfo{
		From:  from,
		To:    to,
		Trade: info,
	}, nil
}
