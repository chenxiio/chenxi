package chenxi

import (
	"database/sql"
	"strings"
	"time"

	"github.com/chenxiio/chenxi/models"
	"github.com/pkg/errors"
)

type HistoriesDao struct {
	db       *sql.DB
	bulkstmt *sql.Stmt
}

func (h *HistoriesDao) createTable() error {
	sql := `CREATE TABLE IF NOT EXISTS histories (
		Parm       TEXT,
		Value      BLOB,
		CreateTime INTEGER
	);`
	sql += "CREATE INDEX IF NOT EXISTS idx_createtime_histories ON histories (createtime);"
	_, err := h.db.Exec(sql)
	if err != nil {
		return err
	}
	h.bulkstmt, _ = h.db.Prepare("INSERT INTO histories (Parm, Value, CreateTime) VALUES (?, ?, ?)")

	return nil
}
func (h *HistoriesDao) BulkInsert(histories []*models.His) error {
	// 开始事务
	tx, err := h.db.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	stmt := tx.Stmt(h.bulkstmt)

	// 准备插入语句
	// stmt, err := tx.Prepare("INSERT INTO histories (Parm, Value, CreateTime) VALUES (?, ?, ?)")
	// if err != nil {
	// 	tx.Rollback()
	// 	return errors.Wrap(err, "failed to prepare insert statement")
	// }
	// defer stmt.Close()

	// 执行批量插入
	for _, history := range histories {
		if history.CreateTime == 0 {
			history.CreateTime = time.Now().UnixNano()
		}
		_, err = stmt.Exec(history.Parm, history.Value, history.CreateTime)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, "failed to execute insert statement")
		}
	}
	// 提交事务
	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}
	return nil
}
func (h *HistoriesDao) Insert(history *models.His) error {
	history.CreateTime = time.Now().UnixNano()
	stmt, err := h.db.Prepare("INSERT INTO histories (Parm, Value, CreateTime) VALUES (?, ?, ?)")
	if err != nil {
		return errors.Wrap(err, "failed to prepare insert statement")
	}
	defer stmt.Close()
	_, err = stmt.Exec(history.Parm, history.Value, history.CreateTime)
	if err != nil {
		return errors.Wrap(err, "failed to execute insert statement")
	}
	return nil
}
func (h *HistoriesDao) Update(history *models.His) error {
	stmt, err := h.db.Prepare("UPDATE histories SET Value=?, CreateTime=? WHERE Parm=?")
	if err != nil {
		return errors.Wrap(err, "failed to prepare update statement")
	}
	defer stmt.Close()
	_, err = stmt.Exec(history.Value, history.CreateTime, history.Parm)
	if err != nil {
		return errors.Wrap(err, "failed to execute update statement")
	}
	return nil
}
func (h *HistoriesDao) GetList(from int64, to int64, parm map[string]interface{}) ([]models.His, error) {
	if to == 0 {
		to = time.Now().UnixNano()
	}
	if from == 0 {
		from = to - (7 * 24 * 60 * 60)
	}
	query := "SELECT * FROM histories WHERE CreateTime BETWEEN ? and ?"
	args := []interface{}{from, to}
	clauses := []string{}
	for key, value := range parm {
		clauses = append(clauses, key+"=?")
		args = append(args, value)
	}
	if len(clauses) > 0 {
		query += " AND " + strings.Join(clauses, " AND ")
	}
	rows, err := h.db.Query(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute query")
	}
	defer rows.Close()
	histories := []models.His{}
	for rows.Next() {
		history := models.His{}
		err := rows.Scan(&history.Parm, &history.Value, &history.CreateTime)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		histories = append(histories, history)
	}
	return histories, nil
}
