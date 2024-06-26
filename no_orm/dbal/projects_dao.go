package dbal

// Code generated by a tool. DO NOT EDIT.
// https://sqldalmaker.sourceforge.net/

import (
	"context"
	"sdm_demo_todolist/no_orm/dbal/dto"
)

type ProjectsDao struct {
	ds DataStore
}

// (C)RUD: projects
// Generated/AI values are passed to DTO/model.

func (dao *ProjectsDao) CreateProject(ctx context.Context, item *dto.Project) error {
	sql := `insert into projects (p_name) values (?)`
	row, err := dao.ds.Insert(ctx, sql, "p_id", item.PName)
	if err == nil {
		err = SetRes(&item.PID, row)
	}
	return err
}

// C(R)UD: projects

func (dao *ProjectsDao) ReadProject(ctx context.Context, pid int64) (*dto.Project, error) {
	sql := `select * from projects where p_id=?`
	row, err := dao.ds.QueryRow(ctx, sql, pid)
	if err != nil {
		return nil, err
	}
	item := dto.Project{}
	errMap := make(map[string]int)
	SetInt64(&item.PID, row, "p_id", errMap)
	SetString(&item.PName, row, "p_name", errMap)
	err = ErrMapToErr(errMap)
	return &item, err
}

// CR(U)D: projects

func (dao *ProjectsDao) UpdateProject(ctx context.Context, item *dto.Project) (rowsAffected int64, err error) {
	sql := `update projects set p_name=? where p_id=?`
	rowsAffected, err = dao.ds.Exec(ctx, sql, item.PName, item.PID)
	return
}

// CRU(D): projects

func (dao *ProjectsDao) DeleteProject(ctx context.Context, item *dto.Project) (rowsAffected int64, err error) {
	sql := `delete from projects where p_id=?`
	rowsAffected, err = dao.ds.Exec(ctx, sql, item.PID)
	return
}

func (dao *ProjectsDao) ReadAll(ctx context.Context) (res []*dto.ProjectLi, err error) {
	sql := `select p.*, 
		(select count(*) from tasks where p_id=p.p_id) as p_tasks_count 
		from projects p 
		order by p.p_id`
	errMap := make(map[string]int)
	_onRow := func(row map[string]interface{}) {
		item := dto.ProjectLi{}
		SetInt64(&item.PID, row, "p_id", errMap)
		SetString(&item.PName, row, "p_name", errMap)
		SetInt64(&item.PTasksCount, row, "p_tasks_count", errMap)
		res = append(res, &item)
	}
	err = dao.ds.QueryAllRows(ctx, sql, _onRow)
	if err == nil {
		err = ErrMapToErr(errMap)
	}
	return
}

func (dao *ProjectsDao) GetProjectIds(ctx context.Context) (res []int64, err error) {
	sql := `select p.*, 
		(select count(*) from tasks where p_id=p.p_id) as p_tasks_count 
		from projects p 
		order by p.p_id`
	err = dao.ds.QueryAll(ctx, sql, &res)
	return
}

func (dao *ProjectsDao) GetProjectID(ctx context.Context) (res []int64, err error) {
	sql := `select p.*, 
		(select count(*) from tasks where p_id=p.p_id) as p_tasks_count 
		from projects p 
		order by p.p_id`
	err = dao.ds.QueryAll(ctx, sql, &res)
	return
}

func (dao *ProjectsDao) PID(ctx context.Context) (res int64, err error) {
	sql := `select * from projects`
	err = dao.ds.Query(ctx, sql, &res)
	return
}

func (dao *ProjectsDao) ProjectID(ctx context.Context, pID string) (res int64, err error) {
	sql := `select * from projects where p_id=?`
	err = dao.ds.Query(ctx, sql, &res, pID)
	return
}

func (dao *ProjectsDao) GetPIds(ctx context.Context, pID int64) (res []int64, err error) {
	sql := `select * from projects where p_id=?`
	err = dao.ds.QueryAll(ctx, sql, &res, pID)
	return
}

func (dao *ProjectsDao) GetPID(ctx context.Context, pID int64) (res int64, err error) {
	sql := `select * from projects where p_id=?`
	err = dao.ds.Query(ctx, sql, &res, pID)
	return
}
