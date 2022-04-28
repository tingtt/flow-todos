package todo

import (
	"flow-todos/mysql"
	"time"
)

func Skip(userId uint64, id uint64) (t Todo, overUntil bool, notFound bool, repeatNotFound bool, dateNotFound bool, invalidUnit bool, err error) {
	// Generate query
	queryStr := "UPDATE todos SET"
	var queryParams []interface{}

	// Get old
	t, notFound, err = Get(userId, id)
	if err != nil {
		return
	}
	if notFound {
		return
	}

	// Repeat exists ?
	if t.Repeat == nil {
		repeatNotFound = true
		return
	}
	if t.Date == nil {
		dateNotFound = true
		return
	}

	// TODO: if out from sprint due

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
		return
	}
	queryStr += " date = ?"
	queryParams = append(queryParams, nextDate)
	if nextTime != nil {
		queryStr += ", time = ?"
		queryParams = append(queryParams, nextTime)
	}

	// Update row
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
	_, err = stmt.Exec(queryParams...)
	if err != nil {
		return
	}

	t.Date = &nextDate
	if nextTime != nil {
		t.Time = nextTime
	}

	return
}
