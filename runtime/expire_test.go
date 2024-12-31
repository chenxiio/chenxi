package runtime

import (
	"database/sql"

	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestExpire(t *testing.T) {
	// 创建内存数据库
	db, err := sql.Open("sqlite3", "db.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	// 创建测试表
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS test_table (id INTEGER PRIMARY KEY, name TEXT,createtime INTEGER)")
	if err != nil {
		t.Fatal(err)
	}
	// 插入一些测试数据
	_, err = db.Exec("INSERT INTO test_table (name,createtime) VALUES ('John',?)", time.Now().UnixNano())
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO test_table (name,createtime) VALUES ('John',?)", time.Now().UnixNano())
	if err != nil {
		t.Fatal(err)
	}
	// 创建Expire对象
	expire := &Expire{db: db, cfg: ExpireCfg{Enable: true, Day: 0}}
	// 调用Expire方法进行测试
	err = expire.Expire()
	if err != nil {
		t.Errorf("Expire failed: %s", err.Error())
	}
	// 检查是否成功执行了ExpireTable方法
	// 这里假设ExpireTable方法会删除表中的数据
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM test_table").Scan(&count)
	if err != nil {
		t.Errorf("Failed to query test_table: %s", err.Error())
	}
	if count != 0 {
		t.Errorf("ExpireTable failed: test_table still has %d rows", count)
	}
}
