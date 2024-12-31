package job

import (
	"database/sql"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type ControlJobDAO struct {
	db *sql.DB
}

var ControlJobDAOInstance *ControlJobDAO
var cjdaoonce sync.Once

func InitControlJobDAOInstance(db *sql.DB) *ControlJobDAO {
	cjdaoonce.Do(func() {
		cjdo := ControlJobDAO{db: db}
		err := cjdo.createTable()
		if err != nil {
			slog.Error(err.Error())
			panic(err)
		}
		ControlJobDAOInstance = &cjdo
	})
	return ControlJobDAOInstance
}
func (dao *ControlJobDAO) createTable() error {
	query := `CREATE TABLE IF NOT EXISTS controljobs (
		CONTROL_JOB_ID TEXT PRIMARY KEY,
		CARRIER_ID TEXT,
		PRIORITY INTEGER,
		PROCESS_JOB_LIST TEXT,
		Mode TEXT,
		State TEXT,
		Createtime INTEGER,
		Starttime INTEGER,
		Endtime INTEGER
	);	
	CREATE INDEX IF NOT EXISTS idx_controljobs_createtime ON controljobs (createtime);`
	_, err := dao.db.Exec(query)
	if err != nil {
		slog.Debug("Error creating controljobs table:", err)
		return err
	}
	return nil
}

func (dao *ControlJobDAO) Insert(controlJob *ControlJob) error {

	stmt, err := dao.db.Prepare("INSERT INTO controljobs (CONTROL_JOB_ID, CARRIER_ID, PRIORITY, PROCESS_JOB_LIST, Mode, State, Createtime, Starttime, Endtime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		slog.Debug("Error preparing statement:", err)

		return err
	}
	controlJob.Createtime = time.Now().UnixNano()
	if len(controlJob.CONTROL_JOB_ID) == 0 {
		controlJob.CONTROL_JOB_ID = uuid.New().String()
	}
	// Convert CARRIER_ID and PROCESS_JOB_LIST to JSON
	carrierIDJSON, err := json.Marshal(controlJob.CARRIER_ID)
	if err != nil {
		slog.Debug("Error converting CARRIER_ID to JSON:", err)

		return err
	}
	processJobListJSON, err := json.Marshal(controlJob.PROCESS_JOB_LIST)
	if err != nil {
		slog.Debug("Error converting PROCESS_JOB_LIST to JSON:", err)

		return err
	}
	_, err = stmt.Exec(controlJob.CONTROL_JOB_ID, string(carrierIDJSON), controlJob.PRIORITY, string(processJobListJSON), controlJob.Mode, controlJob.State, controlJob.Createtime, controlJob.Starttime, controlJob.Endtime)
	if err != nil {
		slog.Debug("Error executing statement:", err)

		return err
	}

	return nil
}
func (dao *ControlJobDAO) GetByCJID(controlJobID string) (*ControlJob, error) {
	query := "SELECT * FROM controljobs WHERE CONTROL_JOB_ID = ?"
	row := dao.db.QueryRow(query, controlJobID)
	controlJob := &ControlJob{}
	var carrierIDJSON string
	var processJobListJSON string
	err := row.Scan(&controlJob.CONTROL_JOB_ID, &carrierIDJSON, &controlJob.PRIORITY, &processJobListJSON, &controlJob.Mode, &controlJob.State, &controlJob.Createtime, &controlJob.Starttime, &controlJob.Endtime)
	if err != nil {
		slog.Debug("Error scanning controljob:", err)
		return nil, err
	}
	// Convert carrierIDJSON and processJobListJSON from JSON to []string
	err = json.Unmarshal([]byte(carrierIDJSON), &controlJob.CARRIER_ID)
	if err != nil {
		slog.Debug("Error converting CARRIER_ID from JSON:", err)
		return nil, err
	}
	err = json.Unmarshal([]byte(processJobListJSON), &controlJob.PROCESS_JOB_LIST)
	if err != nil {
		slog.Debug("Error converting PROCESS_JOB_LIST from JSON:", err)
		return nil, err
	}
	return controlJob, nil
}
func (dao *ControlJobDAO) Update(controlJob *ControlJob) error {
	// tx, err := dao.db.Begin()
	// if err != nil {
	// 	slog.Debug("Error beginning transaction:", err)
	// 	return err
	// }
	stmt, err := dao.db.Prepare("UPDATE controljobs SET CARRIER_ID=?, PRIORITY=?, PROCESS_JOB_LIST=?, Mode=?, State=?, Createtime=?, Starttime=?, Endtime=? WHERE CONTROL_JOB_ID=?")
	if err != nil {
		slog.Debug("Error preparing update statement:", err)
		return err
	}
	defer stmt.Close()
	// Convert CARRIER_ID and PROCESS_JOB_LIST to JSON
	carrierIDJSON, err := json.Marshal(controlJob.CARRIER_ID)
	if err != nil {
		slog.Debug("Error converting CARRIER_ID to JSON:", err)
		return err
	}
	processJobListJSON, err := json.Marshal(controlJob.PROCESS_JOB_LIST)
	if err != nil {
		slog.Debug("Error converting PROCESS_JOB_LIST to JSON:", err)
		return err
	}
	_, err = stmt.Exec(string(carrierIDJSON), controlJob.PRIORITY, string(processJobListJSON), controlJob.Mode, controlJob.State, controlJob.Createtime, controlJob.Starttime, controlJob.Endtime, controlJob.CONTROL_JOB_ID)
	if err != nil {
		slog.Debug("Error executing update statement:", err)

		return err
	}
	// err = tx.Commit()
	// if err != nil {
	// 	slog.Debug("Error committing transaction:", err)
	// 	return err
	// }
	return nil
}

func (dao *ControlJobDAO) GetList(from int64, to int64, parm map[string]interface{}, neq map[string]interface{}) ([]*ControlJob, error) {
	if to == 0 {
		to = time.Now().UnixNano()
	}
	if from == 0 {
		from = to - (3 * 24 * 60 * 60)
	}
	query := "SELECT * FROM controljobs WHERE createtime BETWEEN ? and ?  "
	args := []interface{}{from, to}
	clauses := []string{}
	for key, value := range parm {
		clauses = append(clauses, key+"=?")
		args = append(args, value)
	}
	if len(clauses) > 0 {
		query += " AND " + strings.Join(clauses, " AND ")
	}
	neqClauses := []string{}
	neqArgs := []interface{}{}
	for key, value := range neq {
		neqClauses = append(neqClauses, key+"!=?")
		neqArgs = append(neqArgs, value)
	}
	if len(neqClauses) > 0 {
		query += " AND " + strings.Join(neqClauses, " AND ")
		args = append(args, neqArgs...)
	}
	rows, err := dao.db.Query(query, args...)
	if err != nil {
		slog.Debug("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()
	controlJobs := []*ControlJob{}
	for rows.Next() {
		controlJob := &ControlJob{}
		var carrierIDJSON string
		var processJobListJSON string
		err := rows.Scan(&controlJob.CONTROL_JOB_ID, &carrierIDJSON, &controlJob.PRIORITY, &processJobListJSON, &controlJob.Mode, &controlJob.State, &controlJob.Createtime, &controlJob.Starttime, &controlJob.Endtime)
		if err != nil {
			slog.Debug("Error scanning controljob:", err)
			return nil, err
		}
		// Convert carrierIDJSON and processJobListJSON from JSON to []string
		err = json.Unmarshal([]byte(carrierIDJSON), &controlJob.CARRIER_ID)
		if err != nil {
			slog.Debug("Error converting CARRIER_ID from JSON:", err)
			return nil, err
		}
		err = json.Unmarshal([]byte(processJobListJSON), &controlJob.PROCESS_JOB_LIST)
		if err != nil {
			slog.Debug("Error converting PROCESS_JOB_LIST from JSON:", err)
			return nil, err
		}
		controlJobs = append(controlJobs, controlJob)
	}
	return controlJobs, nil
}
