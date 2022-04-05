package todo

type Todo struct {
	Id            uint64  `json:"id"`
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
	EveryOther *uint       `json:"every_other,omitempty" validate:"omitempty"`
	Date       *uint       `json:"date,omitempty" validate:"omitempty,min=0,max=31"`
	Days       []RepeatDay `json:"days,omitempty" validate:"omitempty,dive"`
}

type RepeatDay struct {
	Day  uint    `json:"day" validate:"required,min=0,max=6"`
	Time *string `json:"time,omitempty" validate:"omitempty,H:M"`
}
