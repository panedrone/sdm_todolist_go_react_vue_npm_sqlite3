# About
A project demonstrating how to use [SQL DAL Maker](https://github.com/panedrone/sqldalmaker) + Golang.

[Part 1](internal/gorm): using "github.com/go-gorm/gorm"<br/>
[Part 2](internal/sqlx): using "github.com/jmoiron/sqlx"<br/>
[Part 3](internal/no_orm): using "database/sql" directly

Front-end is written in both React 16.14 and Vue 2.7. Switch in [handlers.go](shared/handlers.go).

![sdm-todo-app.png](sdm-todo-app.png)

# News

* mod: for sdm 1.321.2601 "mockable golang dao"
* add: [unit-tests](internal/sqlx/handlers) of Gin handlers (thank to "mockable golang dao", sqlx only so far) 

# Docker

<table>
<tr>
    <td>
Creating
    </td>
    <td>
sudo docker compose up -d
    </td>
</tr>
<tr>
    <td>
Stopping and Removing
    </td>
    <td>
sudo docker compose down
    </td>
<tr>
    <td>
Updating
    </td>
    <td>
sudo docker compose up --build -d
    </td>
</table>

sdm.xml:

```xml
<sdm>

    <dto-class name="Project" ref="projects"/>

    <dto-class name="ProjectLi" ref="projects">
        <field type="int64$" column="p_tasks_count"/>
    </dto-class>

    <dto-class name="Task" ref="tasks">
        <custom>
            TId // PK
            PId // FK ref. column -> Project
        </custom>
    </dto-class>

    <dto-class name="TaskLi" ref="tasks">
        <field type="-" column="p_id"/>
        <field type="-" column="t_comments"/>
    </dto-class>

    <dao-class name="ProjectsDao">
        <crud dto="Project"/>
        <query-dto-list method="ReadAll" dto="ProjectLi" ref="get_projects.sql"/>
    </dao-class>

    <dao-class name="TasksDao">
        <crud table="tasks" dto="Task"/>
        <query-dto-list method="ReadByProject(pId)" ref="get_project_tasks.sql" dto="TaskLi"/>
    </dao-class>

</sdm>
```

# Generated code in action

```go
func (h *projectHandlers) ProjectCreate(ctx *gin.Context) {
	var req request.Project
	if err := request.BindJSON(ctx, &req); err != nil {
		return
	}
	if err := dbal.CreateProject(ctx, &dto.Project{PName: req.PName}); err != nil {
		resp.Abort500(ctx, err)
		return
	}
	ctx.Status(http.StatusCreated)
}

func (h *projectHandlers) ProjectsReadAll(ctx *gin.Context) {
	all, err := dbal.ReadAll(ctx)
	if err != nil {
		resp.Abort500(ctx, err)
		return
	}
	resp.JSON(ctx, http.StatusOK, all)
}

func (h *projectHandlers) ProjectRead(ctx *gin.Context) {
	uri, err := request.BindProjectUri(ctx)
	if err != nil {
		return
	}
	pr, err := dbal.ReadProject(ctx, uri.PId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			resp.Abort404NotFound(ctx, err)
			return
		}
		resp.Abort500(ctx, err)
		return
	}
	resp.JSON(ctx, http.StatusOK, pr)

}

func (h *projectHandlers) ProjectUpdate(ctx *gin.Context) {
	uri, err := request.BindProjectUri(ctx)
	if err != nil {
		return
	}
	var req request.Project
	if err := request.BindJSON(ctx, &req); err != nil {
		return
	}
	pr := &dto.Project{PID: uri.PId, PName: req.PName}
	if _, err := dbal.UpdateProject(ctx, pr); err != nil {
		resp.Abort500(ctx, err)
	}
}

func (h *projectHandlers) ProjectDelete(ctx *gin.Context) {
	uri, err := request.BindProjectUri(ctx)
	if err != nil {
		return
	}
	if _, err := dbal.DeleteProject(ctx, &dto.Project{PID: uri.PId}); err != nil {
		resp.Abort500(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
```
