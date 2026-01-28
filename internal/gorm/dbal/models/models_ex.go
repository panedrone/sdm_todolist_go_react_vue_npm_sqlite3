package models

// Hand coded additions

func (t *TaskLi) TableName() string {
	return "tasks"
}

const (
	RefProjectTasks = "RefTasks"
)

type ProjectWithTasks struct {
	Project
	RefTasks []*TaskLi `gorm:"ForeignKey:PId;references:PId"`
}
