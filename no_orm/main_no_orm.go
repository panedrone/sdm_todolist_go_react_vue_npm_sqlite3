package main

import (
	"log"
	"sdm_demo_todolist/no_orm/dbal"
	"sdm_demo_todolist/no_orm/handlers"
	"sdm_demo_todolist/shared"

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

	shared.AssignHandlers(myRouter, "Go, database/sql, SQLite3", projectHandlers, taskHandlers)

	log.Fatal(myRouter.Run(":8080"))
}
