package todo

import (
	"database/sql"
	"flow-todos/mysql"
	"strings"
	"time"

	"github.com/go-playground/validator"
)

type PostBody struct {
	Name          string  `json:"name" validate:"required"`
	Description   *string `json:"description" validate:"omitempty"`
	Date          *string `json:"date" validate:"omitempty,Y-M-D"`
	Time          *string `json:"time" validate:"omitempty,H:M"`
	ExecutionTime *uint   `json:"execution_time" validate:"omitempty"`
	SprintId      *uint64 `json:"sprint_id" validate:"omitempty,gte=1"`
	ProjectId     *uint64 `json:"project_id" validate:"omitempty,gte=1"`
	Completed     *bool   `json:"completed" validate:"omitempty"`
	Repeat        *Repeat `json:"repeat" validate:"omitempty,dive"`
}

func DateStrValidation(fl validator.FieldLevel) bool {
	// `yyyy-mm-dd`
	_, err := time.Parse("2006-1-2", fl.Field().String())
	return err == nil
}

func HMTimeStrValidation(fl validator.FieldLevel) bool {
	// `hh:mm`
	_, err := time.Parse("15:4", fl.Field().String())
	return err == nil
}

func Post(userId uint64, post PostBody) (p Todo, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	// Repeat model
	var idRepeatModel *int64
	if post.Repeat != nil {
		var stmtRepeatModel *sql.Stmt
		stmtRepeatModel, err = db.Prepare("INSERT INTO repeat_models (user_id, until, unit, every_other, date) VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			return
		}
		defer stmtRepeatModel.Close()
		var resultRepeatModel sql.Result
		resultRepeatModel, err = stmtRepeatModel.Exec(userId, post.Repeat.Until, post.Repeat.Unit, post.Repeat.EveryOther, post.Repeat.Date)
		if err != nil {
			return
		}
		var idRepeatModelTmp int64
		idRepeatModelTmp, err = resultRepeatModel.LastInsertId()
		if err != nil {
			return
		}
		idRepeatModel = &idRepeatModelTmp

		// Repeat days
		if post.Repeat.Unit == "week" && len(post.Repeat.Days) != 0 {
			queryStr := "INSERT INTO repeat_days (repeat_model_id, day, time) VALUES"
			var queryParams []interface{}
			for _, day := range post.Repeat.Days {
				queryStr += " (?, ?, ?),"
				queryParams = append(queryParams, idRepeatModel, day.Day, day.Time)
			}
			queryStr = strings.TrimRight(queryStr, ",")
			var stmtRepeatDays *sql.Stmt
			stmtRepeatDays, err = db.Prepare(queryStr)
			if err != nil {
				return
			}
			_, err = stmtRepeatDays.Exec(queryParams...)
			if err != nil {
				return
			}
		}
	}

	// Set defualt value
	if post.Completed == nil {
		completed := false
		post.Completed = &completed
	}

	// Insert DB
	stmt, err := db.Prepare("INSERT INTO todos (user_id, name, description, date, time, execution_time, sprint_id, project_id, completed, repeat_model_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return Todo{}, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(userId, post.Name, post.Description, post.Date, post.Time, post.ExecutionTime, post.SprintId, post.ProjectId, post.Completed, idRepeatModel)
	if err != nil {
		return Todo{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Todo{}, err
	}

	p.Id = uint64(id)
	p.Name = post.Name
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
	if post.SprintId != nil {
		p.SprintId = post.SprintId
	}
	if post.ProjectId != nil {
		p.ProjectId = post.ProjectId
	}
	if post.Completed != nil {
		p.Completed = *post.Completed
	}
	if post.Repeat != nil {
		p.Repeat = post.Repeat
	}

	return
}
