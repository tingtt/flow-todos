package todo

import "flow-todos/mysql"

func Delete(userId uint64, id uint64) (notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return false, err
	}
	defer db.Close()
	stmt, err := db.Prepare(
		`DELETE todos, repeat_models
			FROM todos LEFT JOIN repeat_models ON todos.repeat_model_id = repeat_models.id
			WHERE todos.user_id = ? AND todos.id = ?`,
	)
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(userId, id)
	if err != nil {
		return false, err
	}
	affectedRowCount, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if affectedRowCount == 0 {
		// Not found
		return true, nil
	}

	return false, nil
}
