package todo

import (
	"flow-todos/mysql"
	"strings"
)

type PatchBody struct {
	Name          *string `json:"name" validate:"omitempty"`
	Description   *string `json:"description" validate:"omitempty"`
	Date          *string `json:"date" validate:"omitempty,Y-M-D"`
	Time          *string `json:"time" validate:"omitempty,H:M"`
	ExecutionTime *uint   `json:"execution_time" validate:"omitempty"`
	SprintId      *uint64 `json:"sprint_id" validate:"omitempty,gte=1"`
	ProjectId     *uint64 `json:"project_id" validate:"omitempty,gte=1"`
	Completed     *bool   `json:"completed" validate:"omitempty"`
}

func Patch(userId uint64, id uint64, new PatchBody) (t Todo, notFound bool, err error) {
	// Get old
	t, notFound, err = Get(userId, id)
	if err != nil {
		return
	}
	if notFound {
		return
	}

	// Generate query
	queryStr := "UPDATE schemes SET"
	var queryParams []interface{}
	// Set no update values
	if new.Name != nil {
		queryStr += " name = ?,"
		queryParams = append(queryParams, new.Name)
		t.Name = *new.Name
	}
	if new.Description != nil {
		queryStr += " description = ?,"
		queryParams = append(queryParams, new.Name)
		t.Description = new.Description
	}
	if new.Date != nil {
		queryStr += " date = ?,"
		queryParams = append(queryParams, new.Name)
		t.Date = new.Date
	}
	if new.Time != nil {
		queryStr += " time = ?,"
		queryParams = append(queryParams, new.Name)
		t.Time = new.Time
	}
	if new.ExecutionTime != nil {
		queryStr += " execution_time = ?,"
		queryParams = append(queryParams, new.Name)
		t.ExecutionTime = new.ExecutionTime
	}
	if new.SprintId != nil {
		queryStr += " sprint_id = ?,"
		queryParams = append(queryParams, new.Name)
		t.SprintId = new.SprintId
	}
	if new.ProjectId != nil {
		queryStr += " project_id = ?,"
		queryParams = append(queryParams, new.Name)
		t.ProjectId = new.ProjectId
	}
	if new.Completed != nil {
		queryStr += " completed = ?"
		queryParams = append(queryParams, new.Name)
		t.Completed = *new.Completed
	}
	queryStr = strings.TrimRight(queryStr, ",")
	queryStr += " WHERE user_id = ? AND id = ?"
	queryParams = append(queryParams, userId, id)

	// Update row
	db, err := mysql.Open()
	if err != nil {
		return Todo{}, false, err
	}
	defer db.Close()
	stmtIns, err := db.Prepare(queryStr)
	if err != nil {
		return Todo{}, false, err
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(queryParams...)
	if err != nil {
		return Todo{}, false, err
	}

	return
}
