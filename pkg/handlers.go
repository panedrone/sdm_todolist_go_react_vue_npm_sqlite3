package etc

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

type ProjectHandlers interface {
	ProjectsReadAll(context *gin.Context)
	ProjectCreate(context *gin.Context)
	ProjectRead(context *gin.Context)
	ProjectUpdate(context *gin.Context)
	ProjectDelete(context *gin.Context)
}

type TaskHandlers interface {
	TasksReadByProject(context *gin.Context)
	TaskCreate(context *gin.Context)
	TaskRead(context *gin.Context)
	TaskUpdate(context *gin.Context)
	TaskDelete(context *gin.Context)
}

func AssignHandlers(
	myRouter *gin.Engine,
	whoIam string,
	projectHandlers ProjectHandlers,
	taskHandlers TaskHandlers) {

	//  https://hoohoo.top/blog/20210530112304-golang-tutorial-introduction-gin-html-template-and-how-integration-with-bootstrap/

	// === panedrone: type "http://localhost:8080" to render index.html

	myRouter.Static("/static", "./front-react/static")
	myRouter.StaticFile("/", "./front-react/index.html")

	//myRouter.Static("/static", "./front-vue/static")
	//myRouter.StaticFile("/", "./front-vue/index.html")

	/////////////////////////////////////////

	groupApi := myRouter.Group("/api")

	groupApi.GET("/whoiam", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, whoIam)
	})

	{
		groupProjects := groupApi.Group("/projects")
		groupProjects.GET("/", projectHandlers.ProjectsReadAll)
		groupProjects.POST("/", projectHandlers.ProjectCreate)
		{
			groupProject := groupProjects.Group("/:p_id")
			groupProject.GET("/", projectHandlers.ProjectRead)
			groupProject.PUT("/", projectHandlers.ProjectUpdate)
			groupProject.DELETE("/", projectHandlers.ProjectDelete)
			{
				groupProjectTasks := groupProject.Group("/tasks")
				groupProjectTasks.GET("/", taskHandlers.TasksReadByProject)
				groupProjectTasks.POST("/", taskHandlers.TaskCreate)
			}
		}
	}

	groupTask := groupApi.Group("/tasks/:t_id")
	groupTask.GET("", taskHandlers.TaskRead)
	groupTask.PUT("", taskHandlers.TaskUpdate)
	groupTask.DELETE("", taskHandlers.TaskDelete)
}

func Listen(myRouter *gin.Engine) {
	getAppListen := func() string {
		listen := os.Getenv("APP_LISTEN")
		if listen == "" {
			// fallback: bind explicitly to localhost with default port
			listen = "127.0.0.1:3301"
		}
		return listen
	}

	// Сбор данных
	title := "panedrone's sdm"
	ginVer := "gin version: v1.11.0" // или gin.Version
	goVer := "go sdk: " + runtime.Version()
	appUrl := fmt.Sprintf("http://%s/", getAppListen())
	cpuCount := fmt.Sprintf("CPUs: %d", runtime.NumCPU())

	// Функция для идеального центрирования
	center := func(s string, width int) string {
		padding := width - len(s)
		if padding < 0 {
			return s
		}
		left := padding / 2
		right := padding - left
		return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
	}

	const frameWidth = 51
	logo := fmt.Sprintf(`
 ┌───────────────────────────────────────────────────┐
 │%s│
 │%s│
 │%s│
 │%s│
 │%s│
 └───────────────────────────────────────────────────┘
`,
		center(title, frameWidth),
		center(ginVer, frameWidth),
		center(goVer, frameWidth),
		center(appUrl, frameWidth),
		center(cpuCount, frameWidth),
	)

	fmt.Print(logo)

	log.Fatal(myRouter.Run(getAppListen()))
}
