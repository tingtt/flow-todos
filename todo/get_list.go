package todo

import (
	"flow-todos/mysql"
	"time"
)

type GetListQuery struct {
	Start         *time.Time `query:"start" validate:"omitempty"`
	End           *time.Time `query:"end" validate:"omitempty"`
	ProjectId     *uint64    `query:"project_id" validate:"omitempty,gte=1"`
	WithCompleted bool       `query:"with_completed" validate:"omitempty"`
}

func GetList(userId uint64, q GetListQuery) (todos []Todo, err error) {
	// Generate query
	queryStr := "SELECT id, name, description, date, TIME_FORMAT(time, '%H:%i') AS time, execution_time, sprint_id, project_id, completed FROM todos WHERE user_id = ?"
	queryParams := []interface{}{userId}
	if q.Start != nil && q.End != nil {
		queryStr += " AND ADDTIME(CONVERT(date,DATETIME),COALESCE(time,0)) BETWEEN ? AND ?"
		queryParams = append(queryParams, q.Start.UTC(), q.End.UTC())
	} else if q.Start != nil {
		queryStr += " AND ADDTIME(CONVERT(date,DATETIME),COALESCE(time,0)) >= ?"
		queryParams = append(queryParams, q.Start)
	} else if q.End != nil {
		queryStr += " AND ADDTIME(CONVERT(date,DATETIME),COALESCE(time,0)) <= ?"
		queryParams = append(queryParams, q.End)
	}
	if q.ProjectId != nil {
		queryStr += " AND project_id = ?"
		queryParams = append(queryParams, q.ProjectId)
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
		err = rows.Scan(&t.Id, &t.Name, &t.Description, &t.Date, &t.Time, &t.ExecutionTime, &t.SprintId, &t.ProjectId, &t.Completed)
		if err != nil {
			return
		}
		todos = append(todos, t)
	}

	return
}
