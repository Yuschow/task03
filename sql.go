package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func taskSql() {
	db, err := sql.Open("sqlite3", "task03.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `CREATE TABLE IF NOT EXISTS students (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					name TEXT NOT NULL,
					age INTEGER,
					grade TEXT
				);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlStmt)
	}
	sqlInsert := `insert into students(name, age, grade) values(?, ?, ?)`
	_, err = db.Exec(sqlInsert, "张三", 20, "三年级")
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlStmt)
	}

	sqlSelect := `select * from students where age > 18`
	rows, err := db.Query(sqlSelect)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlSelect)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var age int
		var grade string
		err := rows.Scan(&id, &name, &age, &grade)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("id: %d, name: %s, age: %d, grade: %s\n", id, name, age, grade)
	}

	sqlUpdate := "update students set grade = ? where grade = ?"
	_, err = db.Exec(sqlUpdate, "四年级", "三年级")
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlUpdate)
	}

	sqlDelete := "delete from students where age > ?"
	_, err = db.Exec(sqlDelete, 15)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlDelete)
	}

	sqlAcc := `CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY AUTOINCREMENT,       -- 主键，自增
		balance NUMERIC(15, 2) NOT NULL DEFAULT 0  -- 账户余额，保留两位小数，默认 0
	);`

	sqlTran := `CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY AUTOINCREMENT,       -- 主键，自增
    from_account_id INTEGER NOT NULL,  -- 转出账户 ID
    to_account_id INTEGER NOT NULL,    -- 转入账户 ID
    amount NUMERIC(15, 2) NOT NULL CHECK (amount > 0), -- 金额必须为正

    -- 外键约束，确保账户必须存在
    FOREIGN KEY (from_account_id) REFERENCES accounts(id),
    FOREIGN KEY (to_account_id) REFERENCES accounts(id)
	);`
	_, err = db.Exec(sqlAcc)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlAcc)
	}
	_, err = db.Exec(sqlTran)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlTran)
	}

	// sqlInsertAcc := `insert into accounts(id, balance) values(1, 100)`
	// _, err = db.Exec(sqlInsertAcc)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// sqlInsertAcc2 := `insert into accounts(id, balance) values(2, 80)`
	// _, err = db.Exec(sqlInsertAcc2)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	var balance float64
	err = tx.QueryRow("SELECT balance FROM accounts WHERE id = $1", 1).Scan(&balance)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	if balance < 100 {
		tx.Rollback()
		log.Print("余额不足")
		return
	}

	// 扣除 A
	_, err = tx.Exec("UPDATE accounts SET balance = balance - 100 WHERE id = $1", 1)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	// 增加 B
	_, err = tx.Exec("UPDATE accounts SET balance = balance + 100 WHERE id = $1", 2)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	// 插入转账记录
	_, err = tx.Exec(`
		INSERT INTO transactions (from_account_id, to_account_id, amount)
		VALUES ($1, $2, $3)
	`, 1, 2, 100)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

}
