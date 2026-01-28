package dbal

import (
	"context"
	"sdm_demo_todolist/internal/gorm/dbal/models"
)

func (dao *ProjectsDao) ReadAll(ctx context.Context) (res []*models.ProjectLi, err error) {
	db := dao.ds.Session(ctx)
	subQuery := db.Table("tasks").Select("count(*)").Where("tasks.p_id = projects.p_id")
	err = db.Table("projects").Select("*, (?) as p_tasks_count", subQuery).Find(&res).Error
	return
}
