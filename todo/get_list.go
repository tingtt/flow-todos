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
	queryStr :=
		`SELECT
			todo.id, todo.name, todo.description, todo.date, TIME_FORMAT(todo.time, '%H:%i') AS time, todo.execution_time, todo.sprint_id, todo.project_id, todo.completed,
			rpm.until, rpm.unit, rpm.every_other, rpm.date, rpd.day, TIME_FORMAT(rpd.time, '%H:%i') AS day_time
		FROM todos as todo
			LEFT JOIN repeat_models as rpm ON todo.repeat_model_id = rpm.id
			LEFT JOIN repeat_days as rpd ON rpm.id = rpd.repeat_model_id
		WHERE todo.user_id = ?`
	queryParams := []interface{}{userId}
	if q.Start != nil && q.End != nil {
		queryStr += " AND ADDTIME(CONVERT(todo.date,DATETIME),COALESCE(todo.time,0)) BETWEEN ? AND ?"
		queryParams = append(queryParams, q.Start.UTC(), q.End.UTC())
	} else if q.Start != nil {
		queryStr += " AND ADDTIME(CONVERT(todo.date,DATETIME),COALESCE(todo.time,0)) >= ?"
		queryParams = append(queryParams, q.Start)
	} else if q.End != nil {
		queryStr += " AND ADDTIME(CONVERT(todo.date,DATETIME),COALESCE(todo.time,0)) <= ?"
		queryParams = append(queryParams, q.End)
	}
	if q.ProjectId != nil {
		queryStr += " AND todo.project_id = ?"
		queryParams = append(queryParams, q.ProjectId)
	}
	if !q.WithCompleted {
		queryStr += " AND todo.completed = false"
	}
	queryStr += " ORDER BY todo.id, rpd.day, rpd.time"

	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(queryStr)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(queryParams...)
	if err != nil {
		return
	}

	var tmpTodo Todo
	for rows.Next() {
		t := Todo{}
		var repeatUnit *string
		repeatModel := Repeat{}
		var repeatDayNum *uint
		var repeatDayTime *string
		err = rows.Scan(
			&t.Id, &t.Name, &t.Description, &t.Date, &t.Time, &t.ExecutionTime, &t.SprintId, &t.ProjectId, &t.Completed,
			&repeatModel.Until, &repeatUnit, &repeatModel.EveryOther, &repeatModel.Date, &repeatDayNum, &repeatDayTime,
		)
		if err != nil {
			return
		}
		if repeatUnit != nil {
			repeatModel.Unit = *repeatUnit
			if repeatModel.Unit == "week" && repeatDayNum != nil {
				repeatModel.Days = []RepeatDay{{*repeatDayNum, repeatDayTime}}
			}
			t.Repeat = &repeatModel
		}
		if t.Id == tmpTodo.Id {
			tmpTodo.Repeat.Days = append(tmpTodo.Repeat.Days, t.Repeat.Days...)
		} else {
			if tmpTodo.Id != 0 {
				todos = append(todos, tmpTodo)
			}
			tmpTodo = t
		}
	}
	if tmpTodo.Id != 0 {
		todos = append(todos, tmpTodo)
	}

	return
}
