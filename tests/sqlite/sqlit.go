package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 打开数据库连接
	db, err := sql.Open("sqlite3", "example.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	// 创建表
	createTable := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		age INTEGER
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 插入数据
	insertData := `INSERT INTO users (name, age) VALUES (?, ?)`
	_, err = db.Exec(insertData, "John Doe", 30)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 查询数据
	selectData := `SELECT * FROM users`
	rows, err := db.Query(selectData)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		var age int
		err = rows.Scan(&id, &name, &age)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", id, name, age)
	}
	// 更新数据
	updateData := `UPDATE users SET age = ? WHERE name = ?`
	_, err = db.Exec(updateData, 35, "John Doe")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 删除数据
	deleteData := `DELETE FROM users WHERE name = ?`
	_, err = db.Exec(deleteData, "John Doe")
	if err != nil {
		fmt.Println(err)
		return
	}
}
