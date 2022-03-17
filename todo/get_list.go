package todo

import (
	"flow-todos/mysql"
	"time"
)

func GetList(userId uint64, withCompleted bool, projectId *uint64) (todos []Todo, err error) {
	// Generate query
	queryStr := "SELECT id, name, description, date, TIME_FORMAT(time, '%H:%i') AS time, execution_time, term_id, project_id, completed FROM todos WHERE user_id = ?"
	queryParams := []interface{}{userId}
	if !withCompleted {
		queryStr += " AND completed = false"
	}
	if projectId != nil {
		queryStr += " AND project_id = ?"
		queryParams = append(queryParams, projectId)
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

	rows, err := stmtOut.Query(queryParams...)
	if err != nil {
		return
	}

	for rows.Next() {
		t := Todo{}
		err = rows.Scan(&t.Id, &t.Name, &t.Description, &t.Date, &t.Time, &t.ExecutionTime, &t.TermId, &t.ProjectId, &t.Completed)
		if err != nil {
			return
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
	queryStr := "SELECT id, name, description, date, TIME_FORMAT(time, '%H:%i') AS time, execution_time, term_id, project_id, completed FROM todos WHERE user_id = ?"
	queryParams := []interface{}{userId}
	if dateRange == nil {
		queryStr += " AND date = ?"
		queryParams = append(queryParams, dateStr)
	} else {
		queryStr += " AND date BETWEEN ? AND ?"
		queryParams = append(queryParams, dateStr, date.AddDate(0, 0, int(*dateRange)-1).Format("2006-1-2"))
	}
	if !withCompleted {
		queryStr += " AND completed = false"
	}
	if projectId != nil {
		queryStr += " AND project_id = ?"
		queryParams = append(queryParams, projectId)
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

	rows, err := stmtOut.Query(queryParams...)
	if err != nil {
		return
	}

	for rows.Next() {
		t := Todo{}
		err = rows.Scan(&t.Id, &t.Name, &t.Description, &t.Date, &t.Time, &t.ExecutionTime, &t.TermId, &t.ProjectId, &t.Completed)
		if err != nil {
			return
		}
		todos = append(todos, t)
	}

	return
}
