package todo

import "flow-todos/mysql"

func DeleteAll(userId uint64) (err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()
	stmt, err := db.Prepare(
		`DELETE todos, repeat_models
			FROM todos LEFT JOIN repeat_models ON todos.repeat_model_id = repeat_models.id
			WHERE todos.user_id = ?`,
	)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(userId)
	if err != nil {
		return
	}

	return
}
