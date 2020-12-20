package controllers

import (
	"github.com/gin-gonic/gin"
	"goacs/http/request"
	"goacs/http/response"
	"goacs/models/tasks"
	"goacs/repository"
	"goacs/repository/mysql"
	"gopkg.in/guregu/null.v4"
	"time"
)

type AddTaskRequest struct {
	Event  string `json:"event" validate:"required"`
	Task   string `json:"task" validate:"required"`
	Script string `json:"script"`
}

func AddTask(ctx *gin.Context) {
	var addTaskRequest AddTaskRequest
	_ = ctx.BindJSON(&addTaskRequest)

	validator := request.NewApiValidator(ctx, addTaskRequest)
	verr := validator.Validate()

	if verr != nil {
		response.ResponseValidationErrors(ctx, validator)
		return
	}

	cperepository := mysql.NewCPERepository(repository.GetConnection())
	taskrepository := mysql.NewTasksRepository(repository.GetConnection())
	cpeModel, err := getCPEFromContext(ctx, cperepository)

	if err != nil {
		return
	}

	task := tasks.Task{
		CpeUuid:   cpeModel.UUID,
		Event:     addTaskRequest.Event,
		NotBefore: time.Now(),
		Task:      addTaskRequest.Task,
		Script:    addTaskRequest.Script,
		Infinite:  false,
		CreatedAt: time.Now(),
		DoneAt:    null.Time{},
	}

	taskrepository.AddTaskForCPE(task)
	response.ResponseData(ctx, "")
}
