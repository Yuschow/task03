package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var schema = `
DROP TABLE IF EXISTS employees;
CREATE TABLE employees (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	department TEXT NOT NULL,
	salary DECIMAL
);
INSERT INTO employees (name, department, salary) VALUES
('Alice', 'Engineering', 75000.00),
('Bob', 'Engineering', 72000.00),
('Charlie', 'HR', 58000.00),
('Diana', 'HR', 60000.00),
('Eve', 'Marketing', 65000.00),
('Frank', 'Marketing', 62000.00),
('Grace', 'Finance', 70000.00),
('Heidi', 'Finance', 73000.00),
('Ivan', 'Engineering', 71000.00),
('Judy', 'Marketing', 64000.00);

`

type Employee struct {
	ID         int     `db:"id"`
	Name       string  `db:"name"`
	Department string  `db:"department"`
	Salary     float64 `db:"salary"`
}

func taskSqlx() {
	db, err := sqlx.Connect("sqlite3", "task03.db")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	db.MustExec(schema)
	var employees []Employee
	err = db.Select(&employees, "select * from employees where department = 'Engineering'")
	if err != nil {
		log.Fatalln(err)
	}
	for _, e := range employees {
		fmt.Printf("ID: %d, Name: %s, Dept: %s, Salary: %.2f\n",
			e.ID, e.Name, e.Department, e.Salary)
	}
}
