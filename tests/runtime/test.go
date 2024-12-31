package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slog"
)

type Users struct {
	db *sql.DB
}
type User struct {
	Username string
	Fullname string
	Password string
	Groups   string
}

func (u *Users) getUsers(user *User) ([]User, error) {
	var rows *sql.Rows
	var err error

	if user != nil && user.Username != "" {
		sql := fmt.Sprintf("SELECT username, fullname, password, groups FROM users WHERE username = '%s'", user.Username)
		rows, err = u.db.Query(sql)
	} else {
		sql := "SELECT username, fullname, groups FROM users"
		rows, err = u.db.Query(sql)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var username, fullname, password, groups string
		err = rows.Scan(&username, &fullname, &password, &groups)
		if err != nil {
			return nil, err
		}
		user := User{
			Username: username,
			Fullname: fullname,
			Password: password,
			Groups:   groups,
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *Users) Init(path string) error {

	// prepare query
	sql := "CREATE TABLE IF NOT EXISTS users (username TEXT PRIMARY KEY, fullname TEXT, password TEXT, groups INTEGER);"
	_, err := u.db.Exec(sql)
	if err != nil {
		slog.Error("usrstorage.bind failed!", "e", err)
		return err
	}

	return nil
}
func (u *Users) SetDefault() error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	// if err != nil {
	// 	return fmt.Errorf("密码加密失败:", err)

	// }
	sql := ""
	sql += "INSERT OR REPLACE INTO users (username, fullname, password, groups) VALUES('admin', 'Administrator Account', '" + string(hashedPassword) + "','-1');"
	_, err := u.db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}
func main() {
	// Initialize settings and logger
	// settings := make(map[string]string)
	// logger := log.New(os.Stdout, "", log.LstdFlags)
	// // Call the init function
	// err := Init(settings, logger)
	// if err != nil {
	// 	fmt.Println("Initialization failed:", err)
	// 	return
	// }
	// Continue with the rest of your program
	// ...
}
