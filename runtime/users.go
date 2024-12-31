package runtime

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slog"
)

type User struct {
	Username string
	Fullname string
	Password string
	Groups   string
}

type UsersAPI interface {
	RemoveUser(ctx context.Context, usr string) error                     //perm:none
	SetUser(ctx context.Context, usr, fullname, pwd, groups string) error //perm:none
	GetUsers(ctx context.Context, user *User) ([]User, error)             //perm:none
	SetDefault(ctx context.Context) error                                 //perm:none
}

type Users struct {
	db *sql.DB
}

func (u *Users) Close() error {
	return u.db.Close()
}

/**
 * 从数据库中删除用户
 */
func (u *Users) RemoveUser(ctx context.Context, usr string) error {
	// 准备查询
	sql := "DELETE FROM users WHERE username = '" + usr + "'"
	_, err := u.db.Exec(sql)
	if err != nil {
		slog.Error("usrstorage.remove 失败！" + err.Error())
		return err
	}
	return nil
}
func (u *Users) SetUser(ctx context.Context, usr, fullname, pwd, groups string) error {
	// 准备查询
	exist := false
	data, err := u.GetUsers(context.TODO(), &User{Username: usr})
	if err != nil {
		return err
	}
	if len(data) > 0 {
		exist = true
	}
	pwdhash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	var sql string
	if pwd != "" {
		sql = "INSERT OR REPLACE INTO users (username, fullname, password, groups) VALUES('" + usr + "','" + fullname + "','" + string(pwdhash) + "','" + groups + "');"
		if exist {
			sql = "UPDATE users SET password = '" + string(pwdhash) + "', groups = '" + groups + "', fullname = '" + fullname + "' WHERE username = '" + usr + "';"
		}
	} else {
		sql = "INSERT OR REPLACE INTO users (username, fullname, groups) VALUES('" + usr + "','" + fullname + "','" + groups + "');"
		if exist {
			sql = "UPDATE users SET groups = '" + groups + "', fullname = '" + fullname + "' WHERE username = '" + usr + "';"
		}
	}
	_, err = u.db.Exec(sql)
	if err != nil {
		slog.Error("usrstorage.set 失败！" + err.Error())
		return err
	}
	return nil
}
func (u *Users) GetUsers(ctx context.Context, user *User) ([]User, error) {
	var rows *sql.Rows
	var err error

	if user != nil && user.Username != "" {
		sql := fmt.Sprintf("SELECT username, fullname, password, groups FROM users WHERE username = '%s'", user.Username)
		rows, err = u.db.Query(sql)
	} else {
		sql := "SELECT username, fullname,  password,groups FROM users"
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
		user1 := User{
			Username: username,
			Fullname: fullname,

			Groups: groups,
		}
		if user != nil && user.Username != "" {
			user1.Password = password
		}
		users = append(users, user1)
	}

	return users, nil
}

func (u *Users) Init() error {

	// prepare query
	sql := "CREATE TABLE IF NOT EXISTS users (username TEXT PRIMARY KEY, fullname TEXT, password TEXT, groups INTEGER);"
	_, err := u.db.Exec(sql)
	if err != nil {
		slog.Error("usrstorage.bind failed!", "e", err)
		return err
	}

	return nil
}
func (u *Users) SetDefault(ctx context.Context) error {
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
