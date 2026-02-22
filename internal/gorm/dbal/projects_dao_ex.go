package dbal

import (
	"context"
	"sdm_demo_todolist/internal/gorm/dbal/models"
)

var ReadAll = func(ctx context.Context) (res []*models.ProjectLi, err error) {
	subQuery := RootDb().Table("tasks").Select("count(*)").Where("tasks.p_id = projects.p_id")
	err = WithContext(ctx).Table("projects").Select("*, (?) as p_tasks_count", subQuery).Find(&res).Error
	return
}
