package todo

import (
	"flow-todos/mysql"
)

func Get(userId uint64, id uint64) (t Todo, notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(
		`SELECT
			todo.name, todo.description, todo.date, TIME_FORMAT(todo.time, '%H:%i') AS time, todo.execution_time, todo.sprint_id, todo.project_id, todo.completed,
			rpm.until, rpm.unit, rpm.every_other, rpm.date, rpd.day, TIME_FORMAT(rpd.time, '%H:%i') AS day_time
		FROM todos as todo
			LEFT JOIN repeat_models as rpm ON todo.repeat_model_id = rpm.id
			LEFT JOIN repeat_days as rpd ON rpm.id = rpd.repeat_model_id
		WHERE todo.user_id = ? AND todo.id = ?`,
	)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId, id)
	if err != nil {
		return
	}
	defer rows.Close()

	if !rows.Next() {
		// Not found
		notFound = true
		return
	}
	var repeatUnit *string
	repeatModel := Repeat{}
	var repeatDays []RepeatDay
	var repeatDayNum *uint
	var repeatDayTime *string
	err = rows.Scan(
		&t.Name, &t.Description, &t.Date, &t.Time, &t.ExecutionTime, &t.SprintId, &t.ProjectId, &t.Completed,
		&repeatModel.Until, &repeatUnit, &repeatModel.EveryOther, &repeatModel.Date, &repeatDayNum, &repeatDayTime,
	)
	if err != nil {
		return
	}
	if repeatUnit != nil {
		repeatModel.Unit = *repeatUnit
		if repeatModel.Unit == "week" && repeatDayNum != nil {
			repeatDays = []RepeatDay{{*repeatDayNum, repeatDayTime}}
			for rows.Next() {
				var repeatDayNum2 *uint
				var repeatDayTime2 *string
				err = rows.Scan(
					&t.Name, &t.Description, &t.Date, &t.Time, &t.ExecutionTime, &t.SprintId, &t.ProjectId, &t.Completed,
					&repeatModel.Until, &repeatUnit, &repeatModel.EveryOther, &repeatModel.Date, &repeatDayNum2, &repeatDayTime2,
				)
				if err != nil {
					return
				}
				repeatDays = append(repeatDays, RepeatDay{*repeatDayNum2, repeatDayTime2})
			}
			repeatModel.Days = repeatDays
		}
		t.Repeat = &repeatModel
	}

	t.Id = id
	return
}
