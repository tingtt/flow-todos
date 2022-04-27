package todo

import (
	"encoding/json"
	"flow-todos/mysql"
	"strings"
)

type PatchBody struct {
	Name          *string                 `json:"name" validate:"omitempty"`
	Description   PatchNullJSONString     `json:"description" validate:"omitempty"`
	Date          PatchNullJSONDateString `json:"date" validate:"omitempty,Y-M-D"`
	Time          PatchNullJSONDateString `json:"time" validate:"omitempty,H:M"`
	ExecutionTime PatchNullUint           `json:"execution_time" validate:"omitempty"`
	SprintId      PatchNullJSONUint64     `json:"sprint_id" validate:"dive"`
	ProjectId     PatchNullJSONUint64     `json:"project_id" validate:"dive"`
	Completed     *bool                   `json:"completed" validate:"omitempty"`
}

type PatchNullJSONString struct {
	String **string `validate:"omitempty,gte=1"`
}

type PatchNullJSONDateString struct {
	String **string `validate:"omitempty,Y-M-D"`
}

type PatchNullJSONTimeString struct {
	String **string `validate:"omitempty,H:M"`
}

type PatchNullUint struct {
	UInt **uint `validate:"omitempty,gte=1"`
}

type PatchNullJSONUint64 struct {
	UInt64 **uint64 `validate:"omitempty,gte=1"`
}

func (p *PatchNullJSONString) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set.
	var valueP *string = nil
	if string(data) == "null" {
		// key exists and value is null
		p.String = &valueP
		return nil
	}

	var tmp string
	tmpP := &tmp
	if err := json.Unmarshal(data, &tmp); err != nil {
		// invalid value type
		return err
	}
	// valid value
	p.String = &tmpP
	return nil
}

func (p *PatchNullJSONDateString) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set.
	var valueP *string = nil
	if string(data) == "null" {
		// key exists and value is null
		p.String = &valueP
		return nil
	}

	var tmp string
	tmpP := &tmp
	if err := json.Unmarshal(data, &tmp); err != nil {
		// invalid value type
		return err
	}
	// valid value
	p.String = &tmpP
	return nil
}

func (p *PatchNullJSONTimeString) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set.
	var valueP *string = nil
	if string(data) == "null" {
		// key exists and value is null
		p.String = &valueP
		return nil
	}

	var tmp string
	tmpP := &tmp
	if err := json.Unmarshal(data, &tmp); err != nil {
		// invalid value type
		return err
	}
	// valid value
	p.String = &tmpP
	return nil
}

func (p *PatchNullUint) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set.
	var valueP *uint = nil
	if string(data) == "null" {
		// key exists and value is null
		p.UInt = &valueP
		return nil
	}

	var tmp uint
	tmpP := &tmp
	if err := json.Unmarshal(data, &tmp); err != nil {
		// invalid value type
		return err
	}
	// valid value
	p.UInt = &tmpP
	return nil
}

func (p *PatchNullJSONUint64) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set.
	var valueP *uint64 = nil
	if string(data) == "null" {
		// key exists and value is null
		p.UInt64 = &valueP
		return nil
	}

	var tmp uint64
	tmpP := &tmp
	if err := json.Unmarshal(data, &tmp); err != nil {
		// invalid value type
		return err
	}
	// valid value
	p.UInt64 = &tmpP
	return nil
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
	queryStr := "UPDATE todos SET"
	var queryParams []interface{}
	// Set no update values
	if new.Name != nil {
		queryStr += " name = ?,"
		queryParams = append(queryParams, new.Name)
		t.Name = *new.Name
	}
	if new.Description.String != nil {
		if *new.Description.String != nil {
			queryStr += " description = ?,"
			queryParams = append(queryParams, **new.Description.String)
			t.Description = *new.Description.String
		} else {
			queryStr += " description = ?,"
			queryParams = append(queryParams, nil)
			t.Description = nil
		}
	}
	if new.Date.String != nil {
		if *new.Date.String != nil {
			queryStr += " date = ?,"
			queryParams = append(queryParams, **new.Date.String)
			t.Date = *new.Date.String
		} else {
			queryStr += " date = ?,"
			queryParams = append(queryParams, nil)
			t.Date = nil
		}
	}
	if new.Time.String != nil {
		if *new.Time.String != nil {
			queryStr += " time = ?,"
			queryParams = append(queryParams, **new.Time.String)
			t.Time = *new.Time.String
		} else {
			queryStr += " time = ?,"
			queryParams = append(queryParams, nil)
			t.Time = nil
		}
	}
	if new.ExecutionTime.UInt != nil {
		if *new.ExecutionTime.UInt != nil {
			queryStr += " execution_time = ?,"
			queryParams = append(queryParams, **new.ExecutionTime.UInt)
			t.ExecutionTime = *new.ExecutionTime.UInt
		} else {
			queryStr += " execution_time = ?,"
			queryParams = append(queryParams, nil)
			t.ExecutionTime = nil
		}
	}
	if new.SprintId.UInt64 != nil {
		if *new.SprintId.UInt64 != nil {
			queryStr += " sprint_id = ?,"
			queryParams = append(queryParams, **new.SprintId.UInt64)
			t.SprintId = *new.SprintId.UInt64
		} else {
			queryStr += " sprint_id = ?,"
			queryParams = append(queryParams, nil)
			t.SprintId = nil
		}
	}
	if new.ProjectId.UInt64 != nil {
		if *new.ProjectId.UInt64 != nil {
			queryStr += " parent_id = ?,"
			queryParams = append(queryParams, **new.ProjectId.UInt64)
			t.ProjectId = *new.ProjectId.UInt64
		} else {
			queryStr += " parent_id = ?,"
			queryParams = append(queryParams, nil)
			t.ProjectId = nil
		}
	}
	if new.Completed != nil {
		queryStr += " completed = ?"
		queryParams = append(queryParams, new.Completed)
		t.Completed = *new.Completed
	}
	queryStr = strings.TrimRight(queryStr, ",")
	queryStr += " WHERE user_id = ? AND id = ?"
	queryParams = append(queryParams, userId, id)

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

	return
}
