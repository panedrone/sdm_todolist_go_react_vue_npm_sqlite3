package dbal

import (
	"context"
	models2 "sdm_demo_todolist/internal/gorm/dbal/models"

	"gorm.io/gorm"
)

// Hand coded additions

// 1. "Task" as "Model" and "TaskLi" as Result, no "Select" --> "SELECT t_id, t_date, t_subject, t_priority FROM ...

//      ----- this is the best one so far -----

func (dao *TasksDao) ReadProjectTasks(ctx context.Context, pId int64) (res []*models2.TaskLi, err error) {
	queryModel := &models2.Task{PID: pId}
	err = dao.ds.Session(ctx).Model(queryModel).
		// Select("t_id", "t_date", "t_subject", "t_priority").
		Where(queryModel). // https://gist.github.com/WangYihang/7d43d70db432ff8f3a0a88425bfca7f2
		Order("t_date, t_id").Find(&res).Error

	return
}

// 2. "TaskLi" for both "Model" and Result, requires "Select" --> SELECT * FROM `tasks` WHERE p_id = 2 ORDER BY t_date, t_id

func (dao *TasksDao) _ReadProjectTasks2(ctx context.Context, pId int64) (res []*models2.TaskLi, err error) {
	queryModel := &models2.TaskLi{PID: pId}
	err = dao.ds.Session(ctx).Model(queryModel).
		// Select("t_id", "t_date", "t_subject", "t_priority").
		Where(queryModel).
		Order("t_date, t_id").Find(&res).Error

	return
}

// 3. Using "Table", requires "Select" --> SELECT * FROM `tasks` WHERE p_id = 2 ORDER BY t_date, t_id

func (dao *TasksDao) _ReadProjectTasks3(ctx context.Context, pId int64) (res []*models2.TaskLi, err error) {
	err = dao.ds.Session(ctx).Table("tasks").
		// Select("t_id", "t_date", "t_subject", "t_priority").
		Where(&models2.TaskLi{PID: pId}).
		Order("t_date, t_id").Find(&res).Error

	return
}

// 4. The case "direct TaskLi": no "Table", no "Model". Requires "Select" --> SELECT * FROM `tasks` WHERE p_id = 2 ORDER BY t_date, t_id

func (dao *TasksDao) _ReadProjectTasks4(ctx context.Context, pId int64) (res []*models2.TaskLi, err error) {
	err = dao.ds.Session(ctx).
		// Select("t_id", "t_date", "t_subject", "t_priority").
		Where(&models2.TaskLi{PID: pId}).
		Order("t_date, t_id").Find(&res).Error

	return
}

// 5. Using "Preload" for "educational purposes", requires "Select".

func (dao *TasksDao) _ReadProjectTasks5(ctx context.Context, pId int64) (res []*models2.TaskLi, err error) {
	var queryModel = &models2.ProjectWithTasks{
		Project: models2.Project{PID: pId},
	}
	err = dao.ds.Session(ctx).Model(queryModel).Preload(models2.RefProjectTasks,
		func(db *gorm.DB) *gorm.DB {
			// Use "Select" because "Preload" default issues "SELECT * FROM ..."
			return db.Select("t_id", "p_id", "t_date", "t_subject", "t_priority").
				Order("t_date, t_id")
		}).
		Where(queryModel).Take(queryModel).Error

	if err == nil {
		res = queryModel.RefTasks
	}

	return
}
