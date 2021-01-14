package dbquery

import (
	"database/sql"
	"errors"

	util "../util"
	_ "github.com/go-sql-driver/mysql"
)

func CreateTradeInfo(conn, tid, wid_to, value, ttoken string, cid int, expire, timestamp int64) error {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStr := "INSERT INTO trade_info (tid,tto,coin_idx,tvalue,ttoken,expire_time,create_time) values (?,?,?,?,?,?,?)"
	_, err = db.Exec(sqlStr, tid, wid_to, cid, value, ttoken, expire, timestamp)
	if err != nil {
		return err
	}

	return nil
}
func UpdateTradeInfo(conn, tid, from string, timestamp int64) (*TradeInfo, error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sname := ""
	wid := ""
	var info TradeInfo
	sqlStr1 := "SELECT t.tid,t.coin_idx,t.tvalue,t.ttoken,t.expire_time,t.create_time,c.sname,t.tto FROM trade_info as t, coin_info as c WHERE t.coin_idx=c.idx AND tid=?"
	row1 := db.QueryRow(sqlStr1, tid)
	err = row1.Scan(&info.Tid, &info.Coin_idx, &info.Value, &info.Ttoken, &info.Expire, &info.Create_time, &sname, &wid)
	if err != nil {
		return nil, err
	}

	balance := ""
	sqlStr2 := "SELECT balance FROM wallet_balance as b, (SELECT * FROM user_wallet WHERE uid=?) as u WHERE b.wid=u.wid AND b.coin_idx=?"
	row2 := db.QueryRow(sqlStr2, from, info.Coin_idx)
	err = row2.Scan(&balance)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	b := checkValueBalance(info.Value, balance)
	if !b {
		return &info, errors.New("coin " + sname + " balance insufficient.")
	}

	sqlStr3 := "UPDATE trade_info SET tfrom=(SELECT wid FROM user_wallet WHERE uid=?), update_time=? WHERE tid=?"
	_, err = db.Exec(sqlStr3, from, timestamp, tid)
	if err != nil {
		return nil, err
	}

	sqlStr4 := "SELECT uid FROM user_wallet WHERE wid=?"
	err = db.QueryRow(sqlStr4, wid).Scan(&info.To)
	if err != nil {
		return nil, err
	}

	return &info, nil
}
func ConfirmTradeInfo(conn, tid string) (*TradeBalance, error) {
	timestamp := util.Timestamp()
	// TODO: 待檢查送confirm的人與tid中的to是同一個人
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	balance, err := calcConfirmTrading(db, tid)
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
	_, err = tx.Exec(sqlStr3, tid, balance.To.Wid, balance.From.Wid, balance.Trade.Coin_idx, balance.Trade.Value, balance.Trade.Ttoken, balance.Trade.Expire, balance.Trade.Create_time, balance.Trade.Update_time, timestamp, "T")
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	sqlStr4 := "DELETE FROM trade_info WHERE tid=?"
	_, err = tx.Exec(sqlStr4, tid)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &TradeBalance{
		Cid:     balance.Trade.Coin_idx,
		Balance: balance.From.Balance,
	}, nil
}

func calcConfirmTrading(db *sql.DB, tid string) (*balanceInfo, error) {
	var info TradeInfo
	info.Tid = tid

	sqlStr1 := "SELECT tfrom,tto,coin_idx,tvalue,ttoken,expire_time,create_time,update_time FROM trade_info WHERE tid=?"
	err := db.QueryRow(sqlStr1, info.Tid).Scan(&info.From, &info.To, &info.Coin_idx, &info.Value, &info.Ttoken, &info.Expire, &info.Create_time, &info.Update_time)
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
