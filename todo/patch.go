package todo

import "flow-todos/mysql"

type PatchBody struct {
	Name          *string `json:"name" validate:"omitempty"`
	Description   *string `json:"description" validate:"omitempty"`
	Date          *string `json:"date" validate:"omitempty,Y-M-D"`
	Time          *string `json:"time" validate:"omitempty,H:M"`
	ExecutionTime *uint   `json:"execution_time" validate:"omitempty"`
	TermId        *uint64 `json:"term_id" validate:"omitempty,gte=1"`
	ProjectId     *uint64 `json:"project_id" validate:"omitempty,gte=1"`
	Completed     *bool   `json:"completed" validate:"omitempty"`
}

func Patch(userId uint64, id uint64, new PatchBody) (_ Todo, notFound bool, err error) {
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
