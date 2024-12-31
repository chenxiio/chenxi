package job

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/chenxiio/chenxi/models"
	_ "github.com/mattn/go-sqlite3" // 导入SQLite驱动
	"golang.org/x/exp/slog"
)

func TestCreateTable(t *testing.T) {
	// 创建内存数据库
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	dao := UnitjobDAO{db: db}
	err = dao.createTable()
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	// 在这里添加其他测试逻辑，例如插入测试数据并验证结果
}

func TestUnitjobDAO(t *testing.T) {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	// Create a new instance of the UnitjobDAO
	dao := UnitjobDAO{db: db}
	err = dao.createTable()
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	// Insert a new unitjob
	unitjob := &Unitjob{
		Unitjob: models.Unitjob{PROCESS_JOB_ID: "123456",
			MATERIAL_LIST: []string{"material1", "material2"},
			Unit:          "unit1",
			State:         "state1",
			Createtime:    123456790,
			Starttime:     123456790,
			Endtime:       123456791},
	}

	err = dao.Insert(unitjob)
	if err != nil {
		t.Errorf("Error inserting unitjob: %v", err)
	}

	// Get the unitjob by ID
	unitjobID := unitjob.Unitjob_id
	retrievedUnitjob, err := dao.GetByID(unitjobID)
	if err != nil {
		t.Errorf("Error getting unitjob by ID: %v", err)
	}

	// Verify that the retrieved unitjob matches the inserted unitjob
	if !compareUnitjobs(unitjob, retrievedUnitjob) {
		t.Errorf("Retrieved unitjob does not match the inserted unitjob")
	}

	// Get the unitjobs by PROCESS_JOB_ID
	processJobID := unitjob.PROCESS_JOB_ID
	retrievedUnitjobs, err := dao.GetByPJID(processJobID)
	if err != nil {
		t.Errorf("Error getting unitjobs by PROCESS_JOB_ID: %v", err)
	}

	// Verify that the retrieved unitjobs contain the inserted unitjob
	found := false
	for _, u := range retrievedUnitjobs {
		if compareUnitjobs(unitjob, u) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Inserted unitjob not found in retrieved unitjobs")
	}

	// Update the unitjob
	unitjob.State = "state2"
	err = dao.Update(unitjob)
	if err != nil {
		t.Errorf("Error updating unitjob: %v", err)
	}

	// Get the updated unitjob by ID
	retrievedUpdatedUnitjob, err := dao.GetByID(unitjobID)
	if err != nil {
		t.Errorf("Error getting updated unitjob by ID: %v", err)
	}

	// Verify that the retrieved updated unitjob matches the updated unitjob
	if !compareUnitjobs(unitjob, retrievedUpdatedUnitjob) {
		t.Errorf("Retrieved updated unitjob does not match the updated unitjob")
	}
}

func compareUnitjobs(u1, u2 *Unitjob) bool {
	u1Str, err := json.Marshal(u1)
	if err != nil {
		slog.Debug("Error marshaling unitjob:", err)
		return false
	}
	u2Str, err := json.Marshal(u2)
	if err != nil {
		slog.Debug("Error marshaling unitjob:", err)
		return false
	}
	return string(u1Str) == string(u2Str)
}
