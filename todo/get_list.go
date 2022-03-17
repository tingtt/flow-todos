package todo

import (
	"database/sql"
	"flow-todos/mysql"
	"time"
)

func GetList(userId uint64, withCompleted bool, projectId *uint64) (todos []Todo, err error) {
	// Generate query
	queryStr := "SELECT id, name, description, date, TIME_FORMAT(time, '%H:%i') AS time, execution_time, term_id, project_id, completed FROM todos WHERE user_id = ?"
	if !withCompleted {
		queryStr += " AND completed = false"
	}
	if projectId != nil {
		queryStr += " AND project_id = ?"
	}
	queryStr += " ORDER BY date, time"

	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	stmtOut, err := db.Prepare(queryStr)
	if err != nil {
		return
	}
	defer stmtOut.Close()

	var rows *sql.Rows
	if projectId == nil {
		rows, err = stmtOut.Query(userId)
	} else {
		rows, err = stmtOut.Query(userId, *projectId)
	}
	if err != nil {
		return
	}

	for rows.Next() {
		// TODO: uint64に対応
		var (
			id            uint64
			name          string
			description   sql.NullString
			date          sql.NullString
			time          sql.NullString
			executionTime sql.NullInt16
			termId        sql.NullInt64
			projectId     sql.NullInt64
			completed     bool
		)
		err = rows.Scan(&id, &name, &description, &date, &time, &executionTime, &termId, &projectId, &completed)
		if err != nil {
			return
		}

		t := Todo{Id: id, Name: name, Completed: completed}
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

		todos = append(todos, t)
	}

	return
}

func GetListDate(userId uint64, dateStr string, dateRange *uint, withCompleted bool, projectId *uint64) (todos []Todo, invalidDateStr bool, invalidRange bool, err error) {
	// Validate params
	date, err := time.Parse("20060102", dateStr)
	if err != nil {
		date, err = time.Parse("2006-1-2", dateStr)
		if err != nil {
			err = nil
			invalidDateStr = true
			return
		}
	}
	if dateRange != nil && *dateRange <= 1 {
		invalidRange = true
		return
	}

	// Generate query
	queryStr := ""
	if dateRange == nil {
		queryStr = "SELECT id, name, description, date, TIME_FORMAT(time, '%H:%i') AS time, execution_time, term_id, project_id, completed FROM todos WHERE user_id = ? AND date = ?"
	} else {
		queryStr = "SELECT id, name, description, date, TIME_FORMAT(time, '%H:%i') AS time, execution_time, term_id, project_id, completed FROM todos WHERE user_id = ? AND date BETWEEN ? AND ?"
	}
	if !withCompleted {
		queryStr += " AND completed = false"
	}
	if projectId != nil {
		queryStr += " AND project_id = ?"
	}
	queryStr += " ORDER BY date, time"

	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	stmtOut, err := db.Prepare(queryStr)
	if err != nil {
		return
	}
	defer stmtOut.Close()

	var rows *sql.Rows
	if dateRange == nil {
		if projectId == nil {
			rows, err = stmtOut.Query(userId, dateStr)
		} else {
			rows, err = stmtOut.Query(userId, dateStr, *projectId)
		}
	} else {
		dateEnd := date.AddDate(0, 0, int(*dateRange)-1)
		if projectId == nil {
			rows, err = stmtOut.Query(userId, dateStr, dateEnd.Format("2006-1-2"))
		} else {
			rows, err = stmtOut.Query(userId, dateStr, dateEnd.Format("2006-1-2"), *projectId)
		}
	}
	if err != nil {
		return
	}

	for rows.Next() {
		// TODO: uint64に対応
		var (
			id            uint64
			name          string
			description   sql.NullString
			date          sql.NullString
			time          sql.NullString
			executionTime sql.NullInt16
			termId        sql.NullInt64
			projectId     sql.NullInt64
			completed     bool
		)
		err = rows.Scan(&id, &name, &description, &date, &time, &executionTime, &termId, &projectId, &completed)
		if err != nil {
			return
		}

		t := Todo{Id: id, Name: name, Completed: completed}
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

		todos = append(todos, t)
	}

	return
}
