package todo

import (
	"fmt"
	"sort"
	"time"
)

type Todo struct {
	Id            uint64  `json:"id,omitempty"`
	OriginalId    uint64  `json:"original_id,omitempty"`
	Name          string  `json:"name"`
	Description   *string `json:"description,omitempty"`
	Date          *string `json:"date,omitempty"`
	Time          *string `json:"time,omitempty"`
	ExecutionTime *uint   `json:"execution_time,omitempty"`
	SprintId      *uint64 `json:"sprint_id,omitempty"`
	ProjectId     *uint64 `json:"project_id,omitempty"`
	Completed     bool    `json:"completed"`
	Repeat        *Repeat `json:"repeat,omitempty"`
}

type Repeat struct {
	Until      *string     `json:"until,omitempty" validate:"omitempty,Y-M-D"`
	Unit       string      `json:"unit" validate:"required,oneof=day week month"`
	EveryOther *uint       `json:"every_other,omitempty" validate:"omitempty,gte=1"`
	Date       *uint       `json:"date,omitempty" validate:"omitempty,min=0,max=31"`
	Days       []RepeatDay `json:"days,omitempty" validate:"omitempty,dive"`
}

type RepeatDay struct {
	Day  uint    `json:"day" validate:"required,min=0,max=6"`
	Time *string `json:"time,omitempty" validate:"omitempty,H:M"`
}

func (r *Repeat) GetNext(year int, month time.Month, day int) (nextDate string, nextTime *string, overUntil bool, invalidUnit bool, err error) {
	date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	switch r.Unit {
	case "day":
		// next day
		date = date.AddDate(0, 0, 1)
		if r.EveryOther != nil {
			// every others
			date = date.AddDate(0, 0, int(*r.EveryOther))
		}

		nextDate = date.Format("2006-01-02")

	case "week":
		currentDay := date.Weekday()
		// Sort asc by day of week
		sort.Slice(r.Days, func(i, j int) bool {
			return r.Days[i].Day < r.Days[j].Day
		})
		// Find next
		nextDay := time.Weekday(r.Days[0].Day)
		for _, rd := range r.Days {
			if currentDay < time.Weekday(rd.Day) {
				nextDay = time.Weekday(rd.Day)
				nextTime = rd.Time
				break
			}
		}
		// every other ?
		if r.EveryOther != nil && nextDay < currentDay {
			date = date.AddDate(0, 0, 7)
		}
		// Create time.Time
		date = date.AddDate(0, 0, 1)
		for date.Weekday() != nextDay {
			date = date.AddDate(0, 0, 1)
		}

		nextDate = date.Format("2006-01-02")

	case "month":
		// next month
		currentMonth := date.Month()
		targetMonth := currentMonth + 1
		if r.EveryOther == nil {
			date = date.AddDate(0, 1, 0)
		} else {
			// every other
			targetMonth += time.Month(*r.EveryOther)
			date = date.AddDate(0, 1+int(*r.EveryOther), 0)
		}
		for targetMonth != date.Month() {
			date = date.AddDate(0, 0, -1)
		}
		if r.Date != nil && date.Day() < int(*r.Date) {
			for date.Day() != int(*r.Date) && date.Month() == targetMonth {
				date = date.AddDate(0, 0, 1)
			}
			if date.Month() != targetMonth {
				date = date.AddDate(0, 0, -1)
			}
		}
		nextDate = date.Format("2006-01-02")

	default:
		invalidUnit = true
	}

	if r.Until != nil {
		var until time.Time
		until, err = time.Parse("2006-1-2", *r.Until)
		if err != nil {
			return
		}
		// Over until ?
		if date.After(until) {
			overUntil = true
		}
	}
	return
}

func (t *Todo) GetScheduledRepeats(start *time.Time, end time.Time) (todos []Todo, noRepeat bool, noDate bool, invalidUnit bool, err error) {
	if t.Repeat == nil {
		noRepeat = true
		return
	}

	if t.Date == nil {
		noDate = true
		return
	}

	// Create time.Time by todo.Date and todo.Time
	var datetime time.Time
	if t.Time == nil {
		datetime, err = time.Parse("2006-1-2", *t.Date)
	} else {
		datetime, err = time.Parse("2006-1-2T15:4", *t.Date+"T"+*t.Time)
	}
	if err != nil {
		return
	}

	var current time.Time
	if start == nil {
		current = datetime
	} else if !datetime.Before(*start) {
		current = datetime
		// Add this Todo
		todos = append(todos, *t)
	} else {
		current = *start
	}

	for !current.After(end) {
		var nextDate string
		var nextTime *string
		var overUntil bool
		nextDate, nextTime, overUntil, invalidUnit, err = t.Repeat.GetNext(current.Year(), current.Month(), current.Day())
		if err != nil {
			return
		}
		if invalidUnit {
			return
		}
		if overUntil {
			break
		}

		if t.Time == nil || nextTime == nil || *nextTime == "" {
			current, err = time.Parse("2006-1-2", nextDate)
		} else {
			current, err = time.Parse("2006-1-2T15:4", nextDate+"T"+*nextTime)
		}
		if err != nil {
			return
		}

		fmt.Printf("current: %v\n", current)
		fmt.Printf("Add: %v\n", !current.After(end))
		if !current.After(end) {
			nextTodo := *t
			nextTodo.OriginalId = nextTodo.Id
			nextTodo.Id = 0
			nextTodo.Date = &nextDate
			nextTodo.Repeat = nil
			if nextTime != nil {
				nextTodo.Time = nextTime
			}
			todos = append(todos, nextTodo)
		}
	}

	return
}
