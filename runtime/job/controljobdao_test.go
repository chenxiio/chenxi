package job

import (
	"database/sql"
	"testing"
	"time"

	"github.com/chenxiio/chenxi/models"
	"github.com/google/uuid"
)

func TestControlJobDAO(t *testing.T) {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	dao := InitControlJobDAOInstance(db)
	if err != nil {
		t.Fatal(err)
	}
	err = dao.createTable()
	if err != nil {
		t.Fatal(err)
	}
	controlJob := &ControlJob{
		ControlJob: models.ControlJob{
			CARRIER_ID: []string{"c1", "c2"},

			PROCESS_JOB_LIST: []string{"p1", "p2", "p3"},
			PRIORITY:         1,
			Mode:             "Auto",
			State:            "OUEUED",
			Createtime:       1000,
			Starttime:        2000,
			Endtime:          3000},
	}
	// for i := 0; i < 1000; i++ {
	// 	controlJob.CONTROL_JOB_ID = uuid.New().String()
	// 	err = dao.Insert(controlJob)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// }
	controlJob.CONTROL_JOB_ID = uuid.New().String()
	err = dao.Insert(controlJob)
	if err != nil {
		t.Fatal(err)
	}
	controlJob.PRIORITY = 5

	dao.Update(controlJob)

	time.Sleep(time.Second)
	// retrievedJob, err := dao.GetByCJID(controlJob.CONTROL_JOB_ID)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if retrievedJob.CONTROL_JOB_ID != controlJob.CONTROL_JOB_ID {
	// 	t.Errorf("Expected PROCESS_JOB_NAME to be %s, got %s", controlJob.CONTROL_JOB_ID, retrievedJob.CONTROL_JOB_ID)
	// }
	// retrievedJobs, err := dao.GetList(100, time.Now().UnixNano(), nil, nil)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Println(len(retrievedJobs))
	// retrievedJobs2, err := dao.GetList(100, time.Now().UnixNano(), map[string]any{"CONTROL_JOB_ID": 1}, nil)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Println(len(retrievedJobs2))
	// if retrievedJob.PROCESS_JOB_NAME != processJob.PROCESS_JOB_NAME {
	// 	t.Errorf("Expected PROCESS_JOB_NAME to be %s, got %s", processJob.PROCESS_JOB_NAME, retrievedJob.PROCESS_JOB_NAME)
	// }
	// More assertions as necessary...
}
