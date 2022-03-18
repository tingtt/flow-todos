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
}
