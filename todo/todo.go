package todo

import (
	"database/sql"
	"flow-todos/mysql"
)

type Todo struct {
	Id            uint64
	Name          string  `json:"name"`
	Description   *string `json:"description,omitempty"`
	Date          *string `json:"date,omitempty"`
	Time          *string `json:"time,omitempty"`
	ExecutionTime *uint   `json:"execution_time,omitempty"`
	TermId        *uint64 `json:"term_id,omitempty"`
	ProjectId     *uint64 `json:"project_id,omitempty"`
	Completed     bool    `json:"completed"`
}

type Post struct {
	Name          *string `json:"name" validate:"required"`
	Description   *string `json:"description" validate:"omitempty"`
	Date          *string `json:"date" validate:"omitempty"`
	Time          *string `json:"time" validate:"omitempty"`
	ExecutionTime *uint   `json:"execution_time" validate:"omitempty"`
	TermId        *uint64 `json:"term_id" validate:"omitempty"`
	ProjectId     *uint64 `json:"project_id" validate:"omitempty"`
	Completed     *bool   `json:"completed" validate:"omitempty"`
}

type Patch struct {
	Name          *string `json:"name" validate:"omitempty"`
	Description   *string `json:"description" validate:"omitempty"`
	Date          *string `json:"date" validate:"omitempty"`
	Time          *string `json:"time" validate:"omitempty"`
	ExecutionTime *uint   `json:"execution_time" validate:"omitempty"`
	TermId        *uint64 `json:"term_id" validate:"omitempty"`
	ProjectId     *uint64 `json:"project_id" validate:"omitempty"`
	Completed     *bool   `json:"completed" validate:"omitempty"`
}

func Get(userId uint64, id uint64) (t Todo, notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return Todo{}, false, err
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT name, description, date, time, execution_time, term_id, project_id, completed FROM todos WHERE user_id = ? AND id = ?")
	if err != nil {
		return Todo{}, false, err
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(userId, id)
	if err != nil {
		return Todo{}, false, err
	}

	// TODO: uint64に対応
	var (
		name          string
		description   sql.NullString
		date          sql.NullString
		time          sql.NullString
		executionTime sql.NullInt16
		termId        sql.NullInt64
		projectId     sql.NullInt64
		completed     bool
	)
	if !rows.Next() {
		// Not found
		return Todo{}, true, nil
	}
	err = rows.Scan(&name, &description, &date, &time, &executionTime, &termId, &projectId, &completed)
	if err != nil {
		return Todo{}, false, err
	}

	t.Id = id
	t.Name = name
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
	t.Completed = completed

	return
}

func Insert(userId uint64, post Post) (p Todo, err error) {
	// Set defualt value
	if post.Completed == nil {
		completed := false
		post.Completed = &completed
	}

	// Insert DB
	db, err := mysql.Open()
	if err != nil {
		return Todo{}, err
	}
	defer db.Close()
	stmtIns, err := db.Prepare("INSERT INTO todos (user_id, name, description, date, time, execution_time, term_id, project_id, completed) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return Todo{}, err
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(userId, post.Name, post.Description, post.Date, post.Time, post.ExecutionTime, post.TermId, post.ProjectId, post.Completed)
	if err != nil {
		return Todo{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Todo{}, err
	}

	p.Id = uint64(id)
	p.Name = *post.Name
	if post.Description != nil {
		p.Description = post.Description
	}
	if post.Date != nil {
		p.Date = post.Date
	}
	if post.Time != nil {
		p.Time = post.Time
	}
	if post.ExecutionTime != nil {
		p.ExecutionTime = post.ExecutionTime
	}
	if post.TermId != nil {
		p.TermId = post.TermId
	}
	if post.ProjectId != nil {
		p.ProjectId = post.ProjectId
	}
	if post.Completed != nil {
		p.Completed = *post.Completed
	}

	return
}

func Update(userId uint64, id uint64, new Patch) (_ Todo, notFound bool, err error) {
	// Get old
	old, notFound, err := Get(userId, id)
	if err != nil {
		return Todo{}, false, err
	}
	if notFound {
		return Todo{}, true, nil
	}

	// Set no update values
	if new.Name == nil {
		new.Name = &old.Name
	}
	if new.Description == nil {
		new.Description = old.Description
	}
	if new.Date == nil {
		new.Date = old.Date
	}
	if new.Time == nil {
		new.Time = old.Time
	}
	if new.ExecutionTime == nil {
		new.ExecutionTime = old.ExecutionTime
	}
	if new.TermId == nil {
		new.TermId = old.TermId
	}
	if new.ProjectId == nil {
		new.ProjectId = old.ProjectId
	}
	if new.Completed == nil {
		new.Completed = &old.Completed
	}

	// Update row
	db, err := mysql.Open()
	if err != nil {
		return Todo{}, false, err
	}
	defer db.Close()
	stmtIns, err := db.Prepare("UPDATE todos SET name = ?, description = ?, date = ?, time = ?, execution_time = ?, term_id = ?, project_id = ?, completed = ? WHERE user_id = ? AND id = ?")
	if err != nil {
		return Todo{}, false, err
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(new.Name, new.Description, new.Date, new.Time, new.ExecutionTime, new.TermId, new.ProjectId, new.Completed, userId, id)
	if err != nil {
		return Todo{}, false, err
	}

	return Todo{id, *new.Name, new.Description, new.Date, new.Time, new.ExecutionTime, new.TermId, new.ProjectId, *new.Completed}, false, nil
}

func Delete(userId uint64, id uint64) (notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return false, err
	}
	defer db.Close()
	stmtIns, err := db.Prepare("DELETE FROM todos WHERE user_id = ? AND id = ?")
	if err != nil {
		return false, err
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(userId, id)
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

func GetList(userId uint64, withCompleted bool, projectId *uint64) (projects []Todo, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	// TODO: sort
	queryStr := "SELECT id, name, description, date, time, execution_time, term_id, project_id, completed FROM todos WHERE user_id = ?"
	if !withCompleted {
		queryStr += " AND completed = false"
	}
	if projectId != nil {
		queryStr += " AND project_id = ?"
	}
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

		projects = append(projects, t)
	}

	return
}
