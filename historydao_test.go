package chenxi

import (
	"database/sql"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/chenxiio/chenxi/models"
	_ "github.com/mattn/go-sqlite3"
)

func TestHistoriesDao_Insert(t *testing.T) {
	// 创建一个内存数据库进行测试
	db, err := sql.Open("sqlite3", "his.db")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()
	// 创建HistoriesDao对象
	historiesDao := &HistoriesDao{
		db: db,
	}
	err = historiesDao.createTable()
	if err != nil {
		t.Fatalf("failed to insert history: %v", err)
	}
	// 创建测试数据
	history := &models.His{
		Parm:       "example",
		Value:      1,
		CreateTime: 1234567890,
	}
	// 插入数据
	err = historiesDao.Insert(history)
	if err != nil {
		t.Fatalf("failed to insert history: %v", err)
	}

	list, err := historiesDao.GetList(0, time.Now().UnixNano(), nil)
	if err != nil {
		t.Fatalf("failed to insert history: %v", err)
	}
	fmt.Println(list)
	history.Value = 0.111
	// 插入数据
	err = historiesDao.Insert(history)
	if err != nil {
		t.Fatalf("failed to insert history: %v", err)
	}
	list, err = historiesDao.GetList(0, math.MaxInt64, nil)
	if err != nil {
		t.Fatalf("failed to insert history: %v", err)
	}
	fmt.Println(list)
	history.Value = "888"
	// 插入数据
	err = historiesDao.Insert(history)
	if err != nil {
		t.Fatalf("failed to insert history: %v", err)
	}
	list, err = historiesDao.GetList(0, math.MaxInt64, nil)
	if err != nil {
		t.Fatalf("failed to insert history: %v", err)
	}
	fmt.Println(list)
	// 可以根据需要进行查询验证
	// ...
	// 其他断言和验证
	// ...
}
