package handler

import (
	"github.com/gin-gonic/gin"
	"homeworkdeliverysystem/errors"
	"homeworkdeliverysystem/pkg/handler/middleware"
	"homeworkdeliverysystem/pkg/service"
	"os"
	"strconv"
	"time"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	//router.MaxMultipartMemory = 8 << 20

	app := router.Group("/")

	handlerTimeout := os.Getenv("HANDLER_TIMEOUT")
	ht, _ := strconv.ParseInt(handlerTimeout, 0, 64)

	app.Use(middleware.Timeout(time.Duration(ht)*time.Second, errors.NewServiceUnavailable()))

	auth := app.Group("/auth")
	{
		auth.POST("/sign-up", h.SignUp)                                          //registration
		auth.POST("/sign-in", h.SignIn)                                          //authentication
		auth.POST("/tokens", h.Tokens)                                           //refresh token
		auth.POST("/sign-out", middleware.AuthUser(h.services.Token), h.SignOut) //signing out
	}

	api := app.Group("/api", middleware.AuthUser(h.services.Token))
	{
		user := api.Group("/user")
		{
			user.GET("/me", h.Me) //get user that currently logged in
			userTasks := user.Group("/tasks")
			{
				userTasks.GET("/", h.GetAllTasks)     //get all tasks ordered by deadline
				userTasks.GET("/:id/file", h.GetFile) //get file from task
			}
		}

		group := api.Group("/group")
		{
			group.GET("/:number/subjects", h.GetSubjects) //get subjects by group number
			group.GET("/:number/students", h.GetStudents) //get students by group number
		}

		tasks := api.Group("/tasks", middleware.Authority(middleware.Teacher, middleware.Admin))
		{
			tasks.POST("/", h.createTask)                                      //create task
			tasks.POST("/update-multiple-with-file", h.UpdateMultipleWithFile) //update tasks with file
			//tasks.POST("/:id/answer", h.answerTask)  //прикрепить ответ на таску

			task := tasks.Group("/:id")
			{
				task.POST("/open", h.OpenTask)   //open task
				task.POST("/close", h.CloseTask) //close task
				//task.POST("/approve", h.approveTask)  //аппрувнуть ответ и закрыть таску
				//task.POST("/rate", h.rateTask)        //оценить таску поинтами
				//task.GET("/answers", h.getAllAnswers) //получить все ответы на таску (в порядке их создания)
			}
		}
	}

	return router
}
