package todo

import (
	"database/sql"
	"encoding/json"
	"flow-todos/mysql"
	"strings"
	"time"
)

type PatchBody struct {
	Name          *string                 `json:"name" validate:"omitempty"`
	Description   PatchNullJSONString     `json:"description" validate:"omitempty"`
	Date          PatchNullJSONDateString `json:"date" validate:"omitempty"`
	Time          PatchNullJSONDateString `json:"time" validate:"omitempty"`
	ExecutionTime PatchNullUint           `json:"execution_time" validate:"omitempty"`
	SprintId      PatchNullJSONUint64     `json:"sprint_id" validate:"omitempty"`
	ProjectId     PatchNullJSONUint64     `json:"project_id" validate:"omitempty"`
	Completed     *bool                   `json:"completed" validate:"omitempty"`
	Repeat        PatchNullJSONRepeat     `json:"repeat" validate:"omitempty"`
}

type PatchRepeatBody struct {
	Until      PatchNullJSONDateString `json:"until" validate:"omitempty"`
	Unit       *string                 `json:"unit" validate:"omitempty,oneof=day week month"`
	EveryOther PatchNullUint           `json:"every_other" validate:"omitempty"`
	Date       PatchNullDayOfMonth     `json:"date" validate:"omitempty"`
	Days       PatchNullSliceRepeatDay `json:"days" validate:"omitempty"`
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

type PatchNullDayOfMonth struct {
	UInt **uint `validate:"omitempty,gte=1,lte=31"`
}

type PatchNullJSONUint64 struct {
	UInt64 **uint64 `validate:"omitempty,gte=1"`
}

type PatchNullSliceRepeatDay struct {
	Slice **[]RepeatDay `validate:"omitempty,dive"`
}

type PatchNullJSONRepeat struct {
	Repeat **PatchRepeatBody `validate:"omitempty"`
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

func (p *PatchNullDayOfMonth) UnmarshalJSON(data []byte) error {
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

func (p *PatchNullSliceRepeatDay) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set.
	var valueP *[]RepeatDay = nil
	if string(data) == "null" {
		// key exists and value is null
		p.Slice = &valueP
		return nil
	}

	var tmp []RepeatDay
	tmpP := &tmp
	if err := json.Unmarshal(data, &tmp); err != nil {
		// invalid value type
		return err
	}
	// valid value
	p.Slice = &tmpP
	return nil
}

func (p *PatchNullJSONRepeat) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set.
	var valueP *PatchRepeatBody = nil
	if string(data) == "null" {
		// key exists and value is null
		p.Repeat = &valueP
		return nil
	}

	var tmp PatchRepeatBody
	tmpP := &tmp
	if err := json.Unmarshal(data, &tmp); err != nil {
		// invalid value type
		return err
	}
	// valid value
	p.Repeat = &tmpP
	return nil
}

func Patch(userId uint64, id uint64, new PatchBody) (t Todo, notFound bool, dateNotFound bool, dateOverUntil bool, noDaysWithWeekly bool, err error) {
	// Get old
	t, notFound, err = Get(userId, id)
	if err != nil {
		return
	}
	if notFound {
		return
	}

	// Open connection
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return
	}
	updated := t

	// Update repeat
	var idRepeatModel **uint64 = nil
	if new.Repeat.Repeat != nil {
		if t.Date == nil && (new.Date.String == nil || *new.Date.String == nil) ||
			new.Date.String != nil && *new.Date.String == nil {
			// Update repeat, but date not found
			dateNotFound = true
			return
		}
		if (*new.Repeat.Repeat).Until.String != nil && *(*new.Repeat.Repeat).Until.String != nil {
			var date, until time.Time
			if new.Date.String != nil && *new.Date.String != nil {
				date, err = time.Parse("2006-01-02", **new.Date.String)
			} else {
				date, err = time.Parse("2006-01-02", *t.Date)
			}
			if err != nil {
				return
			}
			until, err = time.Parse("2006-01-02", **(*new.Repeat.Repeat).Until.String)
			if err != nil {
				return
			}
			if date.After(until) {
				dateOverUntil = true
				return
			}
		}

		if t.Repeat != nil {
			// Repeat alredy exists

			// Repeat days
			if (*new.Repeat.Repeat).Unit == nil && t.Repeat.Unit == "week" ||
				(*new.Repeat.Repeat).Unit != nil && *(*new.Repeat.Repeat).Unit == "week" {
				// Weekly
				if (*new.Repeat.Repeat).Days.Slice != nil && (*(*new.Repeat.Repeat).Days.Slice == nil || len(**(*new.Repeat.Repeat).Days.Slice) == 0) ||
					(*new.Repeat.Repeat).Days.Slice == nil && len(t.Repeat.Days) == 0 {
					// Update `repeat.unit` to week, but `repeat.days` is null or update to null
					noDaysWithWeekly = true
					return
				}

				if (*new.Repeat.Repeat).Days.Slice != nil && *(*new.Repeat.Repeat).Days.Slice != nil {
					// Delete repeat days
					var stmtDeleteRepeatDays *sql.Stmt
					stmtDeleteRepeatDays, err = tx.Prepare("DELETE FROM repeat_days WHERE repeat_model_id = (SELECT repeat_model_id FROM todos WHERE user_id = ? AND id = ?)")
					if err != nil {
						return
					}
					defer stmtDeleteRepeatDays.Close()
					_, err = stmtDeleteRepeatDays.Exec(userId, id)
					if err != nil {
						return
					}

					// Insert repeat days
					queryStrRepeatDays := "INSERT INTO repeat_days (repeat_model_id, day, time)"
					var queryParams []interface{}
					for _, day := range **(*new.Repeat.Repeat).Days.Slice {
						queryStrRepeatDays += " SELECT repeat_model_id, ?, ? FROM todos WHERE id = ? UNION"
						queryParams = append(queryParams, day.Day, day.Time, id)
					}
					// Trim trailing ` UNION`
					queryStrRepeatDays = queryStrRepeatDays[:len(queryStrRepeatDays)-6]
					var stmtRepeatDays *sql.Stmt
					stmtRepeatDays, err = tx.Prepare(queryStrRepeatDays)
					if err != nil {
						if err2 := tx.Rollback(); err2 != nil {
							err = err2
						}
						return
					}
					defer stmtDeleteRepeatDays.Close()
					_, err = stmtRepeatDays.Exec(queryParams...)
					if err != nil {
						if err2 := tx.Rollback(); err2 != nil {
							err = err2
						}
						return
					}

					updated.Repeat.Days = **(*new.Repeat.Repeat).Days.Slice
				}
			} else if len(t.Repeat.Days) != 0 {
				// Not weekly and `repeat.days` in DB not empty

				// Delete repeat days
				var stmtDeleteRepeatDays *sql.Stmt
				stmtDeleteRepeatDays, err = tx.Prepare("DELETE FROM repeat_days WHERE repeat_model_id = (SELECT repeat_model_id FROM todos WHERE user_id = ? AND id = ?)")
				if err != nil {
					return
				}
				defer stmtDeleteRepeatDays.Close()
				_, err = stmtDeleteRepeatDays.Exec(userId, id)
				if err != nil {
					return
				}

				updated.Repeat.Days = nil
			}

			// Rpeat model
			queryStr := "UPDATE repeat_models SET"
			var queryParams []interface{}

			noUpdate := true
			if (*new.Repeat.Repeat).Until.String != nil && *(*new.Repeat.Repeat).Until.String != t.Repeat.Until {
				queryStr += " until = ?,"
				queryParams = append(queryParams, *(*new.Repeat.Repeat).Until.String)
				updated.Repeat.Until = *(*new.Repeat.Repeat).Until.String
				noUpdate = false
			}
			if (*new.Repeat.Repeat).Unit != nil && *(*new.Repeat.Repeat).Unit != t.Repeat.Unit {
				queryStr += " unit = ?,"
				queryParams = append(queryParams, *(*new.Repeat.Repeat).Unit)
				updated.Repeat.Unit = *(*new.Repeat.Repeat).Unit
				noUpdate = false
			}
			if (*new.Repeat.Repeat).EveryOther.UInt != nil && *(*new.Repeat.Repeat).EveryOther.UInt != t.Repeat.EveryOther {
				queryStr += " every_other = ?,"
				queryParams = append(queryParams, *(*new.Repeat.Repeat).EveryOther.UInt)
				updated.Repeat.EveryOther = *(*new.Repeat.Repeat).EveryOther.UInt
				noUpdate = false
			}
			if (*new.Repeat.Repeat).Date.UInt != nil && *(*new.Repeat.Repeat).Date.UInt != t.Repeat.Date {
				queryStr += " date = ?,"
				queryParams = append(queryParams, *(*new.Repeat.Repeat).Date.UInt)
				updated.Repeat.Date = *(*new.Repeat.Repeat).Date.UInt
				noUpdate = false
			}

			if !noUpdate {
				queryStr = strings.TrimRight(queryStr, ",")
				queryStr += " WHERE id = (SELECT repeat_model_id FROM todos WHERE user_id = ? AND id = ?)"
				queryParams = append(queryParams, userId, id)

				var stmtRepeatModel *sql.Stmt
				stmtRepeatModel, err = tx.Prepare(queryStr)
				if err != nil {
					if err2 := tx.Rollback(); err2 != nil {
						err = err2
					}
					return
				}
				defer stmtRepeatModel.Close()
				_, err = stmtRepeatModel.Exec(queryParams...)
				if err != nil {
					if err2 := tx.Rollback(); err2 != nil {
						err = err2
					}
					return
				}
			}
		} else {
			// Repeat not exists
			var stmtRepeatModel *sql.Stmt
			stmtRepeatModel, err = tx.Prepare("INSERT INTO repeat_models (user_id, until, unit, every_other, date) VALUES (?, ?, ?, ?, ?)")
			if err != nil {
				return
			}
			defer stmtRepeatModel.Close()
			var resultRepeatModel sql.Result
			resultRepeatModel, err = stmtRepeatModel.Exec(userId, (*new.Repeat.Repeat).Until.String, (*new.Repeat.Repeat).Unit, (*new.Repeat.Repeat).EveryOther.UInt, (*new.Repeat.Repeat).Date.UInt)
			if err != nil {
				return
			}
			var idRepeatModelTmp int64
			idRepeatModelTmp, err = resultRepeatModel.LastInsertId()
			if err != nil {
				return
			}
			idRepeatModelTmpUint := uint64(idRepeatModelTmp)
			idRepeatModelTmpUintP := &idRepeatModelTmpUint
			idRepeatModel = &idRepeatModelTmpUintP

			if (*new.Repeat.Repeat).Until.String != nil && *(*new.Repeat.Repeat).Until.String != t.Repeat.Until {
				updated.Repeat.Until = *(*new.Repeat.Repeat).Until.String
			}
			if (*new.Repeat.Repeat).Unit != nil && *(*new.Repeat.Repeat).Unit != t.Repeat.Unit {
				updated.Repeat.Unit = *(*new.Repeat.Repeat).Unit
			}
			if (*new.Repeat.Repeat).EveryOther.UInt != nil && *(*new.Repeat.Repeat).EveryOther.UInt != t.Repeat.EveryOther {
				updated.Repeat.EveryOther = *(*new.Repeat.Repeat).EveryOther.UInt
			}
			if (*new.Repeat.Repeat).Date.UInt != nil && *(*new.Repeat.Repeat).Date.UInt != t.Repeat.Date {
				updated.Repeat.Date = *(*new.Repeat.Repeat).Date.UInt
			}

			// Repeat days
			if (*(*new.Repeat.Repeat).Unit == "week" || t.Repeat.Unit == "week") && len(**(*new.Repeat.Repeat).Days.Slice) != 0 {
				queryStr := "INSERT INTO repeat_days (repeat_model_id, day, time) VALUES"
				var queryParams []interface{}
				for _, day := range **(*new.Repeat.Repeat).Days.Slice {
					queryStr += " (?, ?, ?),"
					queryParams = append(queryParams, idRepeatModel, day.Day, day.Time)
				}
				queryStr = strings.TrimRight(queryStr, ",")
				var stmtRepeatDays *sql.Stmt
				stmtRepeatDays, err = tx.Prepare(queryStr)
				if err != nil {
					if err2 := tx.Rollback(); err2 != nil {
						err = err2
					}
					return
				}
				_, err = stmtRepeatDays.Exec(queryParams...)
				if err != nil {
					if err2 := tx.Rollback(); err2 != nil {
						err = err2
					}
					return
				}

				updated.Repeat.Days = **(*new.Repeat.Repeat).Days.Slice
			}
		}
	} else if t.Repeat != nil && t.Repeat.Until != nil && new.Date.String != nil {
		// Repeat exists and no update
		// with update date
		if *new.Date.String == nil {
			dateNotFound = true
			return
		}
		var date, until time.Time
		date, err = time.Parse("2006-01-02", **new.Date.String)
		if err != nil {
			return
		}
		until, err = time.Parse("2006-01-02", *t.Repeat.Until)
		if err != nil {
			return
		}
		if date.After(until) {
			dateOverUntil = true
			return
		}
	}

	// Generate query
	queryStr := "UPDATE todos SET"
	var queryParams []interface{}
	// Set update values
	noUpdate := true
	if new.Name != nil {
		queryStr += " name = ?,"
		queryParams = append(queryParams, new.Name)
		updated.Name = *new.Name
		noUpdate = false
	}
	if new.Description.String != nil {
		if *new.Description.String != nil {
			queryStr += " description = ?,"
			queryParams = append(queryParams, **new.Description.String)
			updated.Description = *new.Description.String
		} else {
			queryStr += " description = ?,"
			queryParams = append(queryParams, nil)
			updated.Description = nil
		}
		noUpdate = false
	}
	if new.Date.String != nil {
		if *new.Date.String != nil {
			queryStr += " date = ?,"
			queryParams = append(queryParams, **new.Date.String)
			updated.Date = *new.Date.String
		} else {
			queryStr += " date = ?,"
			queryParams = append(queryParams, nil)
			updated.Date = nil
		}
		noUpdate = false
	}
	if new.Time.String != nil {
		if *new.Time.String != nil {
			queryStr += " time = ?,"
			queryParams = append(queryParams, **new.Time.String)
			updated.Time = *new.Time.String
		} else {
			queryStr += " time = ?,"
			queryParams = append(queryParams, nil)
			updated.Time = nil
		}
		noUpdate = false
	}
	if new.ExecutionTime.UInt != nil {
		if *new.ExecutionTime.UInt != nil {
			queryStr += " execution_time = ?,"
			queryParams = append(queryParams, **new.ExecutionTime.UInt)
			updated.ExecutionTime = *new.ExecutionTime.UInt
		} else {
			queryStr += " execution_time = ?,"
			queryParams = append(queryParams, nil)
			updated.ExecutionTime = nil
		}
		noUpdate = false
	}
	if new.SprintId.UInt64 != nil {
		if *new.SprintId.UInt64 != nil {
			queryStr += " sprint_id = ?,"
			queryParams = append(queryParams, **new.SprintId.UInt64)
			updated.SprintId = *new.SprintId.UInt64
		} else {
			queryStr += " sprint_id = ?,"
			queryParams = append(queryParams, nil)
			updated.SprintId = nil
		}
		noUpdate = false
	}
	if new.ProjectId.UInt64 != nil {
		if *new.ProjectId.UInt64 != nil {
			queryStr += " parent_id = ?,"
			queryParams = append(queryParams, **new.ProjectId.UInt64)
			updated.ProjectId = *new.ProjectId.UInt64
		} else {
			queryStr += " parent_id = ?,"
			queryParams = append(queryParams, nil)
			updated.ProjectId = nil
		}
		noUpdate = false
	}
	if new.Completed != nil {
		queryStr += " completed = ?,"
		queryParams = append(queryParams, new.Completed)
		updated.Completed = *new.Completed
		noUpdate = false
	}
	if idRepeatModel != nil {
		queryStr += "repeat_model_id = ?"
		queryParams = append(queryParams, *idRepeatModel)
		noUpdate = false
	}
	queryStr = strings.TrimRight(queryStr, ",")
	queryStr += " WHERE user_id = ? AND id = ?"
	queryParams = append(queryParams, userId, id)

	if noUpdate {
		t = updated
		err = tx.Commit()
		return
	}

	// Update row
	stmt, err := tx.Prepare(queryStr)
	if err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			err = err2
		}
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(queryParams...)
	if err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			err = err2
		}
		return
	}

	t = updated
	err = tx.Commit()
	return
}
