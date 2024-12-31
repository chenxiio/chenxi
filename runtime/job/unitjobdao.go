package job

import (
	"database/sql"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"golang.org/x/exp/slog"
)

type UnitjobDAO struct {
	db *sql.DB
}

var UnitjobDAOInstance *UnitjobDAO
var ujdaoonce sync.Once

func InitUnitjobDAOInstance(db *sql.DB) *UnitjobDAO {
	ujdaoonce.Do(func() {
		udo := UnitjobDAO{db: db}
		err := udo.createTable()
		if err != nil {
			slog.Error(err.Error())
			panic(err)
		}
		UnitjobDAOInstance = &udo
	})

	return UnitjobDAOInstance
}

func (dao *UnitjobDAO) createTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS unitjobs (
			unitjob_id INTEGER,
			process_job_id TEXT,
			material_list TEXT,
			unit TEXT,
			state TEXT,
			stepname INTEGER,  
			createtime INTEGER,
			starttime INTEGER,
			endtime INTEGER,
			PRIMARY KEY(unitjob_id AUTOINCREMENT)
		);
		CREATE INDEX IF NOT EXISTS idx_unitjobs_process_job_id ON unitjobs (process_job_id);
		CREATE INDEX IF NOT EXISTS idx_unitjobs_createtime ON unitjobs (createtime);
	`

	_, err := dao.db.Exec(query)
	return err
}

func (dao *UnitjobDAO) Insert(unitjob *Unitjob) error {
	materialListStr, err := json.Marshal(unitjob.MATERIAL_LIST)
	if err != nil {
		slog.Debug("Error marshaling MATERIAL_LIST:", err)
		return err
	}
	unitjob.Createtime = time.Now().UnixNano()
	//unitjob.Starttime = time.Now().UnixNano()
	query := "INSERT INTO unitjobs (PROCESS_JOB_ID, MATERIAL_LIST, Unit, State,stepname, Createtime, Starttime, Endtime) VALUES (?, ?, ?, ?, ?,?, ?, ?)"
	result, err := dao.db.Exec(query, unitjob.PROCESS_JOB_ID, string(materialListStr),
		unitjob.Unit, unitjob.State, unitjob.StepName, unitjob.Createtime, unitjob.Starttime, unitjob.Endtime)
	if err != nil {
		slog.Debug("Error inserting unitjobs:", err)
		return err
	}
	unitjobID, err := result.LastInsertId()
	if err != nil {
		slog.Debug("Error getting last insert ID:", err)
		return err
	}
	unitjob.Unitjob_id = unitjobID
	return nil
}
func (dao *UnitjobDAO) GetByID(unitjobID int64) (*Unitjob, error) {
	query := "SELECT * FROM unitjobs WHERE Unitjob_id = ?"
	row := dao.db.QueryRow(query, unitjobID)
	var unitjob Unitjob
	var materialListStr string
	err := row.Scan(&unitjob.Unitjob_id, &unitjob.PROCESS_JOB_ID, &materialListStr, &unitjob.Unit, &unitjob.State, &unitjob.StepName, &unitjob.Createtime, &unitjob.Starttime, &unitjob.Endtime)
	if err != nil {
		slog.Debug("Error scanning unitjobs:", err)
		return nil, err
	}
	err = json.Unmarshal([]byte(materialListStr), &unitjob.MATERIAL_LIST)
	if err != nil {
		slog.Debug("Error unmarshaling MATERIAL_LIST:", err)
		return nil, err
	}
	return &unitjob, nil
}
func (dao *UnitjobDAO) GetByPJID(processJobID string) ([]*Unitjob, error) {
	query := "SELECT * FROM unitjobs WHERE PROCESS_JOB_ID = ?"
	rows, err := dao.db.Query(query, processJobID)
	if err != nil {
		slog.Debug("Error querying unitjobs:", err)
		return nil, err
	}
	defer rows.Close()
	unitjobs := []*Unitjob{}
	for rows.Next() {
		var unitjob Unitjob
		var materialListStr string
		err := rows.Scan(&unitjob.Unitjob_id, &unitjob.PROCESS_JOB_ID, &materialListStr, &unitjob.Unit, &unitjob.State, &unitjob.StepName, &unitjob.Createtime, &unitjob.Starttime, &unitjob.Endtime)
		if err != nil {
			slog.Debug("Error scanning unitjob:", err)
			return nil, err
		}
		err = json.Unmarshal([]byte(materialListStr), &unitjob.MATERIAL_LIST)
		if err != nil {
			slog.Debug("Error unmarshaling MATERIAL_LIST:", err)
			return nil, err
		}
		unitjobs = append(unitjobs, &unitjob)
	}
	return unitjobs, nil
}
func (dao *UnitjobDAO) Update(unitjob *Unitjob) error {
	materialListStr, err := json.Marshal(unitjob.MATERIAL_LIST)
	if err != nil {
		slog.Debug("Error marshaling MATERIAL_LIST:", err)
		return err
	}
	query := "UPDATE unitjobs SET PROCESS_JOB_ID = ?, MATERIAL_LIST = ?, Unit = ?, State = ?,StepName = ?, Createtime = ?, Starttime = ?, Endtime = ? WHERE Unitjob_id = ?"
	_, err = dao.db.Exec(query, unitjob.PROCESS_JOB_ID, string(materialListStr), unitjob.Unit, unitjob.State, unitjob.StepName, unitjob.Createtime, unitjob.Starttime, unitjob.Endtime, unitjob.Unitjob_id)
	if err != nil {
		slog.Debug("Error updating unitjobs:", err)
		return err
	}
	return nil
}
func (dao *UnitjobDAO) GetList(from int64, to int64, parm map[string]interface{}) ([]*Unitjob, error) {
	if to == 0 {
		to = time.Now().UnixNano()
	}
	if from == 0 {
		from = to - (7 * 24 * 60 * 60)
	}
	query := "SELECT * FROM unitjobs WHERE createtime BETWEEN ? and ?  "
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
	unitjobs := []*Unitjob{}
	for rows.Next() {
		unitjob := &Unitjob{}
		var materialListJSON string
		err := rows.Scan(&unitjob.Unitjob_id, &unitjob.PROCESS_JOB_ID, &materialListJSON, &unitjob.Unit, &unitjob.State, &unitjob.StepName, &unitjob.Createtime, &unitjob.Starttime, &unitjob.Endtime)
		if err != nil {
			slog.Debug("Error scanning unitjob:", err)
			return nil, err
		}
		// Convert materialListJSON from JSON to []string
		err = json.Unmarshal([]byte(materialListJSON), &unitjob.MATERIAL_LIST)
		if err != nil {
			slog.Debug("Error converting MATERIAL_LIST from JSON:", err)
			return nil, err
		}
		unitjobs = append(unitjobs, unitjob)
	}
	// Get additional information
	// Add any additional information you want to include in the response
	// For example:
	// additionalInfo["totalCount"] = len(unitjobs)
	// additionalInfo["from"] = from
	// additionalInfo["to"] = to
	return unitjobs, nil
}
