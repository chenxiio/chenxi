package runtime

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/exp/slog"
)

type ExpireCfg struct {
	Enable bool
	Day    int
}
type Expire struct {
	db  *sql.DB
	cfg ExpireCfg
}

func (e *Expire) ExpireTable(tableName string) error {

	//int.
	ctime := time.Now().Add(time.Hour * 24 * time.Duration(e.cfg.Day))
	_, err := e.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE createtime < ?", tableName), ctime.Unix())

	// 检查遍历过程中是否出现错误

	return err
}

func (e *Expire) Expire() error {
	if !e.cfg.Enable {
		return nil
	}
	// 查询所有表名
	tables, err := e.getalltableName()
	if err != nil {
		return err
	}
	// 遍历查询结果并打印表名
	for _, v := range tables {

		err = e.ExpireTable(v)
		if err != nil {
			slog.Error(err.Error())
		}
	}

	return nil
}

func (e *Expire) getalltableName() ([]string, error) {
	// 查询所有表名
	rows, err := e.db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {

		return nil, err
	}
	defer rows.Close()
	ret := []string{}
	// 遍历查询结果并打印表名
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		ret = append(ret, tableName)
	}

	// 检查遍历过程中是否出现错误
	err = rows.Err()
	if err != nil {
		slog.Error(err.Error())
	}
	return ret, nil
}
