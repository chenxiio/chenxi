package alarm

import (
	"context"
	"database/sql"
	"strings"
	"sync"
	"time"

	"github.com/chenxiio/chenxi/models"
)

var AlarmsDaoInstance *AlarmsDao
var cjdaoonce sync.Once

func InitAlarmsDaoInstance(db *sql.DB) *AlarmsDao {
	cjdaoonce.Do(func() {
		cjdo := AlarmsDao{db: db}
		err := cjdo.createTable()
		if err != nil {
			log.Error(err.Error())
			panic(err)
		}
		AlarmsDaoInstance = &cjdo
	})
	return AlarmsDaoInstance
}

type AlarmsDao struct {
	db *sql.DB
}

func (a *AlarmsDao) createTable() error {

	var sql = "CREATE TABLE if not exists alarms (Sn INTEGER, module TEXT, level INTEGER  DEFAULT 0,aid INTEGER  DEFAULT 0,cleartype INTEGER DEFAULT 0,text TEXT,  createtime INTEGER DEFAULT 0, offtime INTEGER DEFAULT 0, PRIMARY KEY(Sn AUTOINCREMENT));"
	sql += "CREATE INDEX IF NOT EXISTS idx_createtime ON alarms (createtime);"

	_, err := a.db.Exec(sql)
	if err != nil {
		//log.Error("alarmsstorage.failed-to-bind: ", err)
		return err
	}

	return nil
}

/**
 * cleartype -1, Clear all Alarms from table,
 */
func (a *AlarmsDao) ClearAlarms(id int64, cleartype int) error {
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	if cleartype == 255 {

		// 设置最近7天createtime的未消除的报警offtime为当前时间

		sql := "UPDATE alarms SET offtime = ?, cleartype = ? WHERE createtime >= ? AND cleartype = 0;"
		_, err := a.db.Exec(sql, time.Now().UnixNano(), cleartype, sevenDaysAgo.Unix())
		if err != nil {
			return err
		}

	} else {
		sql := "UPDATE alarms SET offtime = ?, cleartype = ? WHERE Sn = ?;"
		_, err := a.db.Exec(sql, time.Now().UnixNano(), cleartype, sevenDaysAgo.Unix(), id)
		if err != nil {
			return err
		}
	}
	return nil
}

/**
 * 设置报警
 */
// func (a *AlarmsDao) SetAlarms(module, text string, level, aid int, createtime int64) (int64, error) {
// 	r, err := a.db.Exec(`INSERT INTO alarms ( Module, Level, Aid, Text, createtime)
// 	VALUES ( ?, ?, ?, ?, ?);`, module, level, aid, text, createtime)
// 	if err != nil {
// 		return -1, err
// 	}
// 	id, err := r.LastInsertId()
// 	if err != nil {
// 		return -1, err
// 	}
// 	return id, nil

// }

/**
 * 返回告警列表
 */
func (a *AlarmsDao) GetAlarms(ctx context.Context) ([]models.Alarm, error) {
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	sql := "SELECT Sn, module, level, aid, cleartype, text, createtime, offtime FROM alarms WHERE createtime >= ? and offtime = 0"
	rows, err := a.db.Query(sql, sevenDaysAgo.Unix())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []models.Alarm{}
	for rows.Next() {
		entry := models.Alarm{}
		err := rows.Scan(&entry.Sn, &entry.Module, &entry.Level, &entry.Aid, &entry.ClearType, &entry.Text, &entry.CreateTime, &entry.OffTime)
		if err != nil {
			return nil, err
		}
		result = append(result, entry)
	}
	return result, nil
}

/**
 * 返回告警历史记录
 */
func (a *AlarmsDao) GetAlarmsHistory(ctx context.Context, start int64, end int64) ([]models.Alarm, error) {

	if end == 0 {
		end = time.Now().UnixNano()
	}
	if start == 0 {
		start = end - (7 * 24 * 60 * 60)
	}
	//sevenDaysAgo := time.Unix(end, 0).AddDate(0, 0, -7)
	//fmt.Println(sevenDaysAgo.Unix(), end, start)
	sql := "SELECT Sn, module, level, aid, cleartype, text, createtime, offtime FROM alarms WHERE createtime BETWEEN ? and ? ORDER BY createtime DESC"
	rows, err := a.db.Query(sql, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []models.Alarm{}
	for rows.Next() {
		entry := models.Alarm{}
		err := rows.Scan(&entry.Sn, &entry.Module, &entry.Level, &entry.Aid, &entry.ClearType, &entry.Text, &entry.CreateTime, &entry.OffTime)
		if err != nil {
			return nil, err
		}
		result = append(result, entry)
	}
	return result, nil
}

func (a *AlarmsDao) Insert(alarm *models.Alarm) error {
	stmt, err := a.db.Prepare("INSERT INTO alarms ( Module, Level, Aid, ClearType, Text, CreateTime, OffTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	r, err := stmt.Exec(alarm.Module, alarm.Level, alarm.Aid, alarm.ClearType, alarm.Text, alarm.CreateTime, alarm.OffTime)
	if err != nil {
		return err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return err
	}
	alarm.Sn = id
	return nil
}

func (a *AlarmsDao) Update(ctx context.Context, alarm models.Alarm) error {
	stmt, err := a.db.Prepare("UPDATE alarms SET Module=?, Level=?, Aid=?, ClearType=?, Text=?, CreateTime=?, OffTime=? WHERE Sn=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(alarm.Module, alarm.Level, alarm.Aid, alarm.ClearType, alarm.Text, alarm.CreateTime, alarm.OffTime, alarm.Sn)
	if err != nil {
		return err
	}
	return nil
}

func (a *AlarmsDao) GetByID(sn int) (*models.Alarm, error) {
	row := a.db.QueryRow("SELECT * FROM alarms WHERE Sn = ?", sn)
	alarm := &models.Alarm{}
	err := row.Scan(&alarm.Sn, &alarm.Module, &alarm.Level, &alarm.Aid, &alarm.ClearType, &alarm.Text, &alarm.CreateTime, &alarm.OffTime)
	if err != nil {
		return nil, err
	}
	return alarm, nil
}

func (a *AlarmsDao) GetList(from int64, to int64, parm map[string]any) ([]models.Alarm, error) {
	if to == 0 {
		to = time.Now().UnixNano()
	}
	if from == 0 {
		from = to - (7 * 24 * 60 * 60)
	}
	query := "SELECT * FROM alarms WHERE createtime BETWEEN ? and ?  "
	args := []interface{}{from, to}
	clauses := []string{}
	for key, value := range parm {
		clauses = append(clauses, key+"=?")
		args = append(args, value)
	}
	if len(clauses) > 0 {
		query += " AND " + strings.Join(clauses, " AND ")
	}
	rows, err := a.db.Query(query, args...)
	if err != nil {
		log.Debug("Error executing query:", err)
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	alarms := []models.Alarm{}
	for rows.Next() {
		alarm := models.Alarm{}
		err := rows.Scan(&alarm.Sn, &alarm.Module, &alarm.Level, &alarm.Aid, &alarm.ClearType, &alarm.Text, &alarm.CreateTime, &alarm.OffTime)
		if err != nil {
			return nil, err
		}
		alarms = append(alarms, alarm)
	}
	return alarms, nil
}
