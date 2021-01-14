package dbquery

import (
	"database/sql"
	"errors"

	util "../util"
	_ "github.com/go-sql-driver/mysql"
)

type SetupRegTemp struct {
	Public_key  string
	Private_key string
	Token       string
}
type VerifyUserReg struct {
	Sys_privkey string
	Account     string
	Psn_pubkey  string
}
type InfoUserSys struct {
	Sys_pubkey  string
	Sys_privkey string
	Account     string
	Psn_pubkey  string
	Token       string
}

// type InfoUserReg struct {
// 	Sys_privkey string
// 	Account     string
// 	Psn_pubkey  string
// }

func CheckAccountRegTemp(conn, account string) error {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return err
	}
	defer db.Close()

	var count int64
	sqlStr := "select count(*) from user_reg where account='" + account + "'"
	err = db.QueryRow(sqlStr).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("register account exist.")
	}

	sqlStr = "select count(*) from reg_temp where account='" + account + "'"
	err = db.QueryRow(sqlStr).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("register account exist.")
	}
	return nil
}
func CreateRegTemp(conn, account, pubKey, privKey, token string) error {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return err
	}
	defer db.Close()

	timestamp := util.Timestamp()
	sqlStr := "INSERT INTO reg_temp (account, public_key, private_key, token, create_time) values (?, ?, ?, ?, ?)"
	_, err = db.Exec(sqlStr, account, pubKey, privKey, token, timestamp)
	if err != nil {
		return err
	}

	return nil
}
func GetSetupRegTemp(conn, account string) (*SetupRegTemp, error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	public_key := ""
	private_key := ""
	token := ""

	sqlStr := "select public_key, private_key, token from reg_temp where account=?"
	row := db.QueryRow(sqlStr, account)
	err = row.Scan(&public_key, &private_key, &token)
	if err != nil {
		return nil, err
	}

	return &SetupRegTemp{
		Public_key:  public_key,
		Private_key: private_key,
		Token:       token,
	}, nil
}
func CreateUserSys(conn, psn_acc, psn_pubkey, psn_pwd, uid, mnemonic, sys_pwd, sys_pubkey, sys_privkey string) error {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	timestamp := util.Timestamp()
	sqlStr1 := "insert into user_reg (uid,account,public_key,passwd,create_time,update_time,status) values (?,?,?,?,?,?,?)"
	_, err = tx.Exec(sqlStr1, uid, psn_acc, psn_pubkey, psn_pwd, timestamp, timestamp, "W")
	if err != nil {
		tx.Rollback()
		return err
	}

	sqlStr2 := "insert into user_sys (uid,mnemonic,passwd,public_key,private_key,idx,create_time,update_time) values (?,?,?,?,?,?,?,?)"
	_, err = tx.Exec(sqlStr2, uid, mnemonic, sys_pwd, sys_pubkey, sys_privkey, 0, timestamp, timestamp)
	if err != nil {
		tx.Rollback()
		return err
	}

	sqlStr3 := "delete from reg_temp where account=?"
	_, err = tx.Exec(sqlStr3, psn_acc)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
func GetVerifyUserReg(conn, uid string) (*VerifyUserReg, error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sys_privkey := ""
	// sqlStr1 := "select private_key from user_sys where uid=?"
	// row1 := db.QueryRow(sqlStr1, uid)
	// err = row1.Scan(&sys_privkey)
	// if err != nil {
	// 	return nil, err
	// }

	account := ""
	psn_pubkey := ""
	sqlStr2 := "select reg.account,reg.public_key,sys.private_key from user_sys as sys, (select * from user_reg where status='W' and uid=?) as reg where sys.uid=reg.uid"
	row2 := db.QueryRow(sqlStr2, uid)
	err = row2.Scan(&account, &psn_pubkey, &sys_privkey)
	if err != nil {
		return nil, err
	}

	return &VerifyUserReg{
		Sys_privkey: sys_privkey,
		Account:     account,
		Psn_pubkey:  psn_pubkey,
	}, nil
}
func RenewRegInfo(conn, uid, token string, timestamp int64) error {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return err
	}
	defer db.Close()

	var count int64
	sqlStr1 := "select count(*) from user_auth where uid='" + uid + "'"
	err = db.QueryRow(sqlStr1).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("register information error, please notify the administrator.")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	sqlStr2 := "insert into user_auth (uid,token,create_time,update_time) values (?,?,?,?)"
	_, err = tx.Exec(sqlStr2, uid, token, timestamp, timestamp)
	if err != nil {
		tx.Rollback()
		return err
	}

	sqlStr3 := "update user_reg set status='N', update_time=? where uid=?"
	_, err = tx.Exec(sqlStr3, timestamp, uid)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
func GetInfoUserReg(conn, uid string) (*VerifyUserReg, error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	account := ""
	sys_privkey := ""
	psn_pubkey := ""
	sqlStr := "select reg.account,reg.public_key,sys.private_key from user_sys as sys, (select * from user_reg where status='N' and uid=?) as reg where sys.uid=reg.uid"
	row := db.QueryRow(sqlStr, uid)
	err = row.Scan(&account, &psn_pubkey, &sys_privkey)
	if err != nil {
		return nil, err
	}

	return &VerifyUserReg{
		Account:     account,
		Sys_privkey: sys_privkey,
		Psn_pubkey:  psn_pubkey,
	}, nil
}
func RenewTokenUserAuth(conn, uid, token string, timestamp int64) error {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	sqlStr := "update user_auth set token=?, update_time=? where uid=?"
	_, err = tx.Exec(sqlStr, token, timestamp, uid)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// func f() {
// 	db, err := sql.Open("mysql", "root:123456@/register")
// 	if err != nil {
// 		fmt.Println("db error:", err)
// 	}
// 	defer db.Close()

// 	err = db.Ping()
// 	if err != nil {
// 		fmt.Println("ping error:", err)
// 	}
// 	// fmt.Println("db:", db)
// 	rows, err := db.Query("SELECT * FROM reg_tmp")
// 	if err != nil {
// 		panic(err.Error()) // proper error handling instead of panic in your app
// 	}

// 	// Get column names
// 	columns, err := rows.Columns()
// 	if err != nil {
// 		panic(err.Error()) // proper error handling instead of panic in your app
// 	}

// 	// Make a slice for the values
// 	values := make([]sql.RawBytes, len(columns))

// 	scanArgs := make([]interface{}, len(values))
// 	for i := range values {
// 		scanArgs[i] = &values[i]
// 	}

// 	// Fetch rows
// 	for rows.Next() {
// 		// get RawBytes from data
// 		err = rows.Scan(scanArgs...)
// 		if err != nil {
// 			panic(err.Error()) // proper error handling instead of panic in your app
// 		}

// 		// Now do something with the data.
// 		// Here we just print each column as a string.
// 		var value string
// 		for i, col := range values {
// 			// Here we can check if the value is nil (NULL value)
// 			if col == nil {
// 				value = "NULL"
// 			} else {
// 				value = string(col)
// 			}
// 			fmt.Println(columns[i], ": ", value)
// 		}
// 		fmt.Println("-----------------------------------")
// 	}
// 	if err = rows.Err(); err != nil {
// 		panic(err.Error()) // proper error handling instead of panic in your app
// 	}
// }
