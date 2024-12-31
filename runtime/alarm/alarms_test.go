package alarm

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/chenxiio/chenxi/models"
	"github.com/stretchr/testify/assert"
)

func TestGetAlarms(t *testing.T) {
	// 创建测试用的数据库连接
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	// 创建 Alarms 对象
	alarms := AlarmsDao{db: db}
	err = alarms.createTable()
	if err != nil {
		t.Fatal(err)
	}
	// 插入一些测试数据
	_, err = db.Exec(`INSERT INTO alarms ( Module, Level, Aid, ClearType, Text, OnTime)
		VALUES ( 'module1', 1, 1, -1, 'alarm1', ?);`,
		time.Now().UnixNano())
	if err != nil {
		t.Fatal(err)
	}
	// 调用 GetAlarms 方法
	ctx := context.Background()
	result, err := alarms.GetAlarms(ctx)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(result)

	err = alarms.ClearAlarms(-1, 255)
	if err != nil {
		t.Fatal(err)
	}

	// 断言结果是否符合预期
	expected := []models.Alarm{
		{
			Sn:         1,
			Module:     "module1",
			Level:      1,
			Aid:        1,
			ClearType:  -1,
			Text:       "alarm1",
			CreateTime: time.Now().UnixNano(),
			OffTime:    time.Now().UnixNano(),
		},
	}
	assert.Equal(t, expected, result)
}

func TestClearAlarms(t *testing.T) {
	// 创建测试数据库连接
	db, err := createTestDBConnection()
	if err != nil {
		t.Fatalf("failed to create test database connection: %v", err)
	}
	defer db.Close()

	// 创建 Alarms 实例
	alarms := &AlarmsDao{db: db}

	// 插入测试数据
	err = insertTestAlarmsData(db)
	if err != nil {
		t.Fatalf("failed to insert test alarms data: %v", err)
	}

	// 模拟上下文
	ctx := context.Background()

	// 清除所有警报
	err = alarms.ClearAlarms(0, 255)
	if err != nil {
		t.Fatalf("failed to clear alarms: %v", err)
	}

	// 验证清除所有警报后的数据状态
	alarmsAfterClear, err := alarms.GetAlarms(ctx)
	if err != nil {
		t.Fatalf("failed to get alarms after clear: %v", err)
	}
	if len(alarmsAfterClear) != 0 {
		t.Fatalf("expected 0 alarms after clear, got %d", len(alarmsAfterClear))
	}

	// 清除指定警报
	err = alarms.ClearAlarms(1, 1)
	if err != nil {
		t.Fatalf("failed to clear alarm: %v", err)
	}

	// 验证清除指定警报后的数据状态
	alarmsAfterClear, err = alarms.GetAlarms(ctx)
	if err != nil {
		t.Fatalf("failed to get alarms after clear: %v", err)
	}
	if len(alarmsAfterClear) != 1 {
		t.Fatalf("expected 1 alarm after clear, got %d", len(alarmsAfterClear))
	}
	// 验证其他数据是否正确

	// 进行其他断言和验证
}

// 创建测试数据库连接
func createTestDBConnection() (*sql.DB, error) {
	// 连接数据库并返回数据库连接

	db, err := sql.Open("sqlite3", "database.db")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	return db, err
}

// 插入测试警报数据
func insertTestAlarmsData(db *sql.DB) error {
	// 插入测试警报数据到数据库

	_, err := db.Exec(`INSERT INTO alarms ( Module, Level, Aid, ClearType, Text, OnTime)
		VALUES ( 'module1', 1, 1, -1, 'alarm1', ?);`,
		time.Now().UnixNano())

	return err
}
func TestGetAlarmsHistory(t *testing.T) {
	// 创建一个 Alarms 实例
	// 创建测试数据库连接
	db, err := createTestDBConnection()
	if err != nil {
		t.Fatalf("failed to create test database connection: %v", err)
	}
	defer db.Close()

	// 创建 Alarms 实例
	alarms := &AlarmsDao{db: db}
	// 定义测试用例的输入和期望的输出
	start := int64(0)
	end := time.Now().UnixNano()
	expectedResult := []models.Alarm{
		// 在这里定义你期望的告警历史记录结果
	}
	// 调用函数获取实际结果
	result, err := alarms.GetAlarmsHistory(context.Background(), start, end)
	if err != nil {
		t.Errorf("GetAlarmsHistory returned an error: %v", err)
	}
	// 比较实际结果和期望结果
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("GetAlarmsHistory returned unexpected result: got %v, want %v", result, expectedResult)
	}
}
func TestSetAlarms(t *testing.T) {
	// 创建测试数据库连接
	db, err := createTestDBConnection()
	if err != nil {
		t.Fatalf("failed to create test database connection: %v", err)
	}
	defer db.Close()

	// 创建 Alarms 实例
	//alarms := &AlarmsDao{db: db}
	// 定义测试用例的输入和期望的输出
	// module := "testModule"
	// text := "testText"
	// level := 1
	// aid := 123
	// expectedResult := int64(1)
	// // 调用函数获取实际结果
	// result, err := alarms.Insert(module, text, level, aid, time.Now().UnixNano())
	// if err != nil {
	// 	t.Errorf("SetAlarms returned an error: %v", err)
	// }
	// // 比较实际结果和期望结果
	// if result != expectedResult {
	// 	t.Errorf("SetAlarms returned unexpected result: got %v, want %v", result, expectedResult)
	// }
}
