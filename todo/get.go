package todo

import (
	"database/sql"
	"flow-todos/mysql"
)

func Get(userId uint64, id uint64) (t Todo, notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return Todo{}, false, err
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT name, description, date, TIME_FORMAT(time, '%H:%i') AS time, execution_time, term_id, project_id, completed FROM todos WHERE user_id = ? AND id = ?")
	if err != nil {
		return Todo{}, false, err
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(userId, id)
	if err != nil {
		return Todo{}, false, err
	}

	// TODO: uint64に対応
	var (
		name          string
		description   sql.NullString
		date          sql.NullString
		time          sql.NullString
		executionTime sql.NullInt16
		termId        sql.NullInt64
		projectId     sql.NullInt64
		completed     bool
	)
	if !rows.Next() {
		// Not found
		return Todo{}, true, nil
	}
	err = rows.Scan(&name, &description, &date, &time, &executionTime, &termId, &projectId, &completed)
	if err != nil {
		return Todo{}, false, err
	}

	t.Id = id
	t.Name = name
	if description.Valid {
		t.Description = &description.String
	}
	if date.Valid {
		t.Date = &date.String
	}
	if time.Valid {
		t.Time = &time.String
	}
	if executionTime.Valid {
		executionTimeTmp := uint(executionTime.Int16)
		t.ExecutionTime = &executionTimeTmp
	}
	if termId.Valid {
		termIdTmp := uint64(termId.Int64)
		t.TermId = &termIdTmp
	}
	if projectId.Valid {
		projectIdTmp := uint64(projectId.Int64)
		t.ProjectId = &projectIdTmp
	}
	t.Completed = completed

	return
}
