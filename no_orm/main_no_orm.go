package main

import (
	"sdm_demo_todolist/no_orm/dbal"
	"sdm_demo_todolist/no_orm/handlers"
	"sdm_demo_todolist/pkg"

	"github.com/gin-gonic/gin"
)

func main() {
	err := dbal.OpenDB()
	if err != nil {
		println(err.Error())
		return
	}
	defer func() {
		_ = dbal.CloseDB()
	}()

	gin.SetMode(gin.ReleaseMode)

	myRouter := gin.New()

	projectHandlers := handlers.NewProjectHandlers()
	taskHandlers := handlers.NewTaskHandlers()

	etc.AssignHandlers(myRouter, "Go, database/sql, SQLite3", projectHandlers, taskHandlers)

	etc.Listen(myRouter)
}
