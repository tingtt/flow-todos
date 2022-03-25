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

	stmt, err := db.Prepare("SELECT name, description, date, TIME_FORMAT(time, '%H:%i') AS time, execution_time, sprint_id, project_id, completed FROM todos WHERE user_id = ? AND id = ?")
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId, id)
	if err != nil {
		return
	}

	if !rows.Next() {
		// Not found
		notFound = true
		return
	}
	err = rows.Scan(&t.Name, &t.Description, &t.Date, &t.Time, &t.ExecutionTime, &t.SprintId, &t.ProjectId, &t.Completed)
	if err != nil {
		return
	}

	t.Id = id
	return
}
