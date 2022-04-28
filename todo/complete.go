package todo

import (
	"database/sql"
	"flow-todos/mysql"
	"time"
)

func Complete(userId uint64, id uint64) (t Todo, new Todo, notFound bool, dateNotFound bool, invalidUnit bool, err error) {
	// Get old
	t, notFound, err = Get(userId, id)
	if err != nil {
		return
	}
	if notFound {
		return
	}

	var db *sql.DB
	db, err = mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	// No repeat
	if t.Repeat == nil {
		// Update row
		var stmt *sql.Stmt
		stmt, err = db.Prepare("UPDATE todos SET completed = ? WHERE user_id = ? AND id = ?")
		if err != nil {
			return
		}
		defer stmt.Close()
		_, err = stmt.Exec(true, userId, id)
		if err != nil {
			return
		}

		t.Completed = true
		return
	}

	// Repeat todo and no date
	if t.Date == nil {
		dateNotFound = true
		return
	}

	/**
	 * Create next repeat todo
	**/

	// Get next date
	var date time.Time
	date, err = time.Parse("2006-1-2", *t.Date)
	if err != nil {
		return
	}
	nextDate, nextTime, overUntil, invalidUnit, err := t.Repeat.GetNext(date.Year(), date.Month(), date.Day())
	if err != nil {
		return
	}
	if invalidUnit {
		return
	}
	if overUntil {
		// Update old
		var stmt *sql.Stmt
		stmt, err = db.Prepare("UPDATE todos SET completed = ? WHERE user_id = ? AND id = ?")
		if err != nil {
			return
		}
		defer stmt.Close()
		_, err = stmt.Exec(true, userId, id)
		if err != nil {
			return
		}

		t.Completed = true
		return
	}

	new = t
	new.Date = &nextDate
	if nextTime != nil {
		new.Time = nextTime
	}

	// TODO: if out from sprint due

	// Insert DB
	stmt2, err := db.Prepare(
		`INSERT INTO todos
			(user_id, name, description, date, time, execution_time, sprint_id, project_id, repeat_model_id)
		SELECT
			user_id, name, description, ?, ?, execution_time, sprint_id, project_id, repeat_model_id
		FROM todos
		WHERE user_id = ? AND id = ?`,
	)
	if err != nil {
		return
	}
	defer stmt2.Close()
	result, err := stmt2.Exec(new.Date, new.Time, userId, id)
	if err != nil {
		return
	}
	newId, err := result.LastInsertId()
	if err != nil {
		return
	}
	new.Id = uint64(newId)

	/**
	 * Update old
	**/

	// Update row
	var stmt *sql.Stmt
	stmt, err = db.Prepare("UPDATE todos SET completed = ?, repeat_model_id = ? WHERE user_id = ? AND id = ?")
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(true, nil, userId, id)
	if err != nil {
		return
	}

	t.Completed = true
	t.Repeat = nil

	return
}
