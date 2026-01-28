A project demonstrating how to use [SQL DAL Maker](https://github.com/panedrone/sqldalmaker) + Golang.

[Part 1](internal/gorm): using "github.com/go-gorm/gorm"<br/>
[Part 2](internal/sqlx): using "github.com/jmoiron/sqlx"<br/>
[Part 3](internal/no_orm): using "database/sql" directly

Front-end is written in both React 16.14 and Vue 2.7. Switch in [handlers.go](shared/handlers.go).

![sdm-todo-app.png](sdm-todo-app.png)

Docker:

<table>
<tr>
    <td>
Creating
    </td>
    <td>
docker-compose up -d
    </td>
</tr>
<tr>
    <td>
Stopping and Removing
    </td>
    <td>
docker-compose down
    </td>
<tr>
    <td>
Updating
    </td>
    <td>
docker-compose up --build -d
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

Generated code in action:

```go
type projectHandlers struct {
	dao *dbal.ProjectsDao
}

func NewProjectHandlers() shared.ProjectHandlers {
	return &projectHandlers{
		dao: dbal.NewProjectsDao(),
	}
}

func (h *projectHandlers) ProjectCreate(ctx *gin.Context) {
	var req request.Project
	if err := request.BindJSON(ctx, &req); err != nil {
		return
	}
	if err := h.dao.CreateProject(ctx, &m.Project{PName: req.PName}); err != nil {
		resp.Abort500(ctx, err)
		return
	}
	ctx.Status(http.StatusCreated)
}

func (h *projectHandlers) ProjectsReadAll(ctx *gin.Context) {
	all, err := h.dao.ReadAll(ctx)
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
	pr, err := h.dao.ReadProject(ctx, uri.PId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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
	pr := &m.Project{PId: uri.PId, PName: req.PName}
	if _, err := h.dao.UpdateProject(ctx, pr); err != nil {
		resp.Abort500(ctx, err)
	}
}

func (h *projectHandlers) ProjectDelete(ctx *gin.Context) {
	uri, err := request.BindProjectUri(ctx)
	if err != nil {
		return
	}
	if _, err := h.dao.DeleteProject(ctx, &m.Project{PId: uri.PId}); err != nil {
		resp.Abort500(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
```
