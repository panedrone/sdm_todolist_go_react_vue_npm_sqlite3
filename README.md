# About

```
┌───────────────────────────────────────────────────┐
│                  panedrone's sdm                  │
│               gin version: v1.12.0                │
│                 go sdk: go1.25.0                  │
│              http://127.0.0.1:3301/               │
│                     CPUs: 20                      │
└───────────────────────────────────────────────────┘
```


A project demonstrating how to use [SQL DAL Maker](https://github.com/panedrone/sqldalmaker) + Golang.

[Part 1](internal/gorm): using "github.com/go-gorm/gorm"<br/>
[Part 2](internal/sqlx): using "github.com/jmoiron/sqlx"<br/>
[Part 3](internal/no_orm): using "database/sql" directly

Front-end is written in both React 18.2.0 and Vue 3.5.32. Switch in [handlers.go](pkg/handlers.go).

![sdm-todo-app.png](sdm-todo-app.png)

# News

* ported to React 18.2.0
* mod: for sdm 1.321.2601 "mockable golang dao"
* add: [unit-tests](internal/sqlx/handlers) of Gin handlers for GORM, sqlx and No-ORM (thanks to "mockable golang dao") 
* add: [Swagger](http://localhost:3301/swagger/index.html) annotations for all handlers

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

# sdm.xml

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

# Generated code in action (example)

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
```
