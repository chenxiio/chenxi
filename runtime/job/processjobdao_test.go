package job

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/chenxiio/chenxi/models"
	_ "github.com/mattn/go-sqlite3"
)

func TestProcessJobDAO(t *testing.T) {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	dao := InitProcessJobDAOInstance(db)
	if err != nil {
		t.Fatal(err)
	}
	err = dao.createTable()
	if err != nil {
		t.Fatal(err)
	}
	processJob := &ProcessJob{ProcessJob: models.ProcessJob{
		PROCESS_JOB_NAME:      "Test Job",
		PROCESS_DEFINITION_ID: "Test Definition",
		PARAMETER_LIST: map[string]string{
			"recipe1": "value1",
			"recipe2": "value2",
		},
		MATERIAL_LIST:     []string{"A1", "B2", "C3"},
		PRIORITY:          "High",
		PROCESS_JOB_NOTES: "Test notes",
		State:             "wait",
		Createtime:        1000,
		Starttime:         2000,
		Endtime:           3000},
	}
	err = dao.Insert(processJob)
	if err != nil {
		t.Fatal(err)
	}
	retrievedJob, err := dao.GetByPJID(processJob.PROCESS_JOB_ID)
	if err != nil {
		t.Fatal(err)
	}
	if retrievedJob.PROCESS_JOB_NAME != processJob.PROCESS_JOB_NAME {
		t.Errorf("Expected PROCESS_JOB_NAME to be %s, got %s", processJob.PROCESS_JOB_NAME, retrievedJob.PROCESS_JOB_NAME)
	}
	retrievedJobs, err := dao.GetList(100, time.Now().UnixNano(), nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(retrievedJobs))
	retrievedJobs2, err := dao.GetList(100, time.Now().UnixNano(), map[string]any{"process_job_id": 1})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(retrievedJobs2))
	// if retrievedJob.PROCESS_JOB_NAME != processJob.PROCESS_JOB_NAME {
	// 	t.Errorf("Expected PROCESS_JOB_NAME to be %s, got %s", processJob.PROCESS_JOB_NAME, retrievedJob.PROCESS_JOB_NAME)
	// }
	// More assertions as necessary...
}
