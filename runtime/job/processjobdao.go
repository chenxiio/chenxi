package job

import (
	"database/sql"
	"encoding/json"

	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/exp/slog"
)

type ProcessJobDAO struct {
	db *sql.DB
}

var ProcessJobDAOInstance *ProcessJobDAO
var pjdaoonce sync.Once

func InitProcessJobDAOInstance(db *sql.DB) *ProcessJobDAO {
	pjdaoonce.Do(func() {

		udo := ProcessJobDAO{db: db}
		err := udo.createTable()
		if err != nil {
			slog.Error(err.Error())
			panic(err)
		}
		ProcessJobDAOInstance = &udo
	})

	return ProcessJobDAOInstance
}
func (dao *ProcessJobDAO) createTable() error {
	query := `CREATE TABLE IF NOT EXISTS processjobs (
		PROCESS_JOB_ID TEXT PRIMARY KEY,
		PROCESS_JOB_NAME TEXT,
		PROCESS_DEFINITION_ID TEXT,
		PARAMETER_LIST TEXT,
		MATERIAL_LIST TEXT,
		PRIORITY TEXT,
		PROCESS_JOB_NOTES TEXT,
		State TEXT,
		Createtime INTEGER,
		Starttime INTEGER,
		Endtime INTEGER
	);	
	CREATE INDEX IF NOT EXISTS idx_processjobs_createtime ON processjobs (createtime);`
	_, err := dao.db.Exec(query)
	if err != nil {
		slog.Debug("Error creating processjobs table:", err)
		return err
	}
	return nil
}

func (dao *ProcessJobDAO) Insert(processJob *ProcessJob) error {
	tx, err := dao.db.Begin()
	if err != nil {
		slog.Debug("Error beginning transaction:", err)
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO processjobs (PROCESS_JOB_ID, PROCESS_JOB_NAME, PROCESS_DEFINITION_ID, PARAMETER_LIST, MATERIAL_LIST, PRIORITY, PROCESS_JOB_NOTES, State, Createtime, Starttime, Endtime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		slog.Debug("Error preparing statement:", err)
		tx.Rollback()
		return err
	}
	processJob.Createtime = time.Now().UnixNano()
	if len(processJob.PROCESS_JOB_ID) == 0 {
		processJob.PROCESS_JOB_ID = uuid.New().String()
	}

	// Convert PARAMETER_LIST to JSON
	parameterListJSON, err := json.Marshal(processJob.PARAMETER_LIST)
	if err != nil {
		slog.Debug("Error converting PARAMETER_LIST to JSON:", err)
		tx.Rollback()
		return err
	}

	// Convert MATERIAL_LIST to JSON
	materialListJSON, err := json.Marshal(processJob.MATERIAL_LIST)
	if err != nil {
		slog.Debug("Error converting MATERIAL_LIST to JSON:", err)
		tx.Rollback()
		return err
	}

	_, err = stmt.Exec(processJob.PROCESS_JOB_ID, processJob.PROCESS_JOB_NAME, processJob.PROCESS_DEFINITION_ID, string(parameterListJSON), string(materialListJSON), processJob.PRIORITY, processJob.PROCESS_JOB_NOTES, processJob.State, processJob.Createtime, processJob.Starttime, processJob.Endtime)
	if err != nil {
		slog.Debug("Error executing statement:", err)
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		slog.Debug("Error committing transaction:", err)
		return err
	}

	return nil
}

func (dao *ProcessJobDAO) GetByPJID(processJobID string) (*ProcessJob, error) {
	query := "SELECT * FROM processjobs WHERE PROCESS_JOB_ID = ?"
	row := dao.db.QueryRow(query, processJobID)
	processJob := &ProcessJob{}
	var parameterListJSON string
	var materialListJSON string
	err := row.Scan(&processJob.PROCESS_JOB_ID, &processJob.PROCESS_JOB_NAME, &processJob.PROCESS_DEFINITION_ID, &parameterListJSON, &materialListJSON, &processJob.PRIORITY, &processJob.PROCESS_JOB_NOTES, &processJob.State, &processJob.Createtime, &processJob.Starttime, &processJob.Endtime)
	if err != nil {
		slog.Debug("Error scanning processjob:", err)
		return nil, err
	}
	// Convert parameterListJSON from JSON to map[string]string
	err = json.Unmarshal([]byte(parameterListJSON), &processJob.PARAMETER_LIST)
	if err != nil {
		slog.Debug("Error converting PARAMETER_LIST from JSON:", err)
		return nil, err
	}
	// Convert materialListJSON from JSON to []string
	err = json.Unmarshal([]byte(materialListJSON), &processJob.MATERIAL_LIST)
	if err != nil {
		slog.Debug("Error converting MATERIAL_LIST from JSON:", err)
		return nil, err
	}
	return processJob, nil
}
func (dao *ProcessJobDAO) Update(processJob *ProcessJob) error {
	// tx, err := dao.db.Begin()
	// if err != nil {
	// 	slog.Debug("Error beginning transaction:", err)
	// 	return err
	// }
	stmt, err := dao.db.Prepare("UPDATE processjobs SET PROCESS_JOB_NAME=?, PROCESS_DEFINITION_ID=?, PARAMETER_LIST=?, MATERIAL_LIST=?, PRIORITY=?, PROCESS_JOB_NOTES=?, State=?, Createtime=?, Starttime=?, Endtime=? WHERE PROCESS_JOB_ID=?")
	if err != nil {
		slog.Debug("Error preparing update statement:", err)
		return err
	}
	defer stmt.Close()
	// Convert PARAMETER_LIST to JSON
	parameterListJSON, err := json.Marshal(processJob.PARAMETER_LIST)
	if err != nil {
		slog.Debug("Error converting PARAMETER_LIST to JSON:", err)
		return err
	}
	// Convert MATERIAL_LIST to JSON
	materialListJSON, err := json.Marshal(processJob.MATERIAL_LIST)
	if err != nil {
		slog.Debug("Error converting MATERIAL_LIST to JSON:", err)
		return err
	}
	_, err = stmt.Exec(processJob.PROCESS_JOB_NAME, processJob.PROCESS_DEFINITION_ID, string(parameterListJSON), string(materialListJSON), processJob.PRIORITY, processJob.PROCESS_JOB_NOTES, processJob.State, processJob.Createtime, processJob.Starttime, processJob.Endtime, processJob.PROCESS_JOB_ID)
	if err != nil {
		slog.Debug("Error executing update statement:", err)
		//tx.Rollback()
		return err
	}
	// err = tx.Commit()
	// if err != nil {
	// 	slog.Debug("Error committing transaction:", err)
	// 	return err
	// }
	return nil
}

func (dao *ProcessJobDAO) GetList(from int64, to int64, parm map[string]any) ([]*ProcessJob, error) {
	if to == 0 {
		to = time.Now().UnixNano()
	}
	if from == 0 {
		from = to - (7 * 24 * 60 * 60)
	}
	query := "SELECT * FROM processjobs WHERE createtime BETWEEN ? and ?  "
	args := []interface{}{from, to}
	clauses := []string{}
	for key, value := range parm {
		clauses = append(clauses, key+"=?")
		args = append(args, value)
	}
	if len(clauses) > 0 {
		query += " AND " + strings.Join(clauses, " AND ")
	}
	rows, err := dao.db.Query(query, args...)
	if err != nil {
		slog.Debug("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()
	processJobs := []*ProcessJob{}
	for rows.Next() {
		processJob := &ProcessJob{}
		var parameterListJSON string
		var materialListJSON string
		err := rows.Scan(&processJob.PROCESS_JOB_ID, &processJob.PROCESS_JOB_NAME, &processJob.PROCESS_DEFINITION_ID, &parameterListJSON, &materialListJSON, &processJob.PRIORITY, &processJob.PROCESS_JOB_NOTES, &processJob.State, &processJob.Createtime, &processJob.Starttime, &processJob.Endtime)
		if err != nil {
			slog.Debug("Error scanning processjob:", err)
			return nil, err
		}
		// Convert parameterListJSON from JSON to map[string]string
		err = json.Unmarshal([]byte(parameterListJSON), &processJob.PARAMETER_LIST)
		if err != nil {
			slog.Debug("Error converting PARAMETER_LIST from JSON:", err)
			return nil, err
		}
		// Convert materialListJSON from JSON to []string
		err = json.Unmarshal([]byte(materialListJSON), &processJob.MATERIAL_LIST)
		if err != nil {
			slog.Debug("Error converting MATERIAL_LIST from JSON:", err)
			return nil, err
		}
		processJobs = append(processJobs, processJob)
	}

	return processJobs, nil
}
