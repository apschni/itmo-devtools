package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"homeworkdeliverysystem/dto"
	apperrors "homeworkdeliverysystem/errors"
	"homeworkdeliverysystem/model"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func (h *Handler) createTask(ctx *gin.Context) {
	userFromContext, _ := ctx.Get("user")
	user := userFromContext.(*model.User)

	var req dto.CreateTaskReq

	if ok := bindData(ctx, &req); !ok {
		return
	}

	now := time.Now()

	task := &model.Task{
		Label:      req.Label,
		Subject:    req.Subject,
		Text:       req.Text,
		Deadline:   req.Deadline,
		Points:     req.Points,
		Closed:     false,
		TeacherId:  user.Id,
		StudentId:  req.StudentId,
		FileName:   "",
		CreatedAt:  now,
		UpdatedAt:  now,
		IsKeyPoint: req.IsKeyPoint,
	}
	c := ctx.Request.Context()

	id, err := h.services.Task.Create(c, task)
	if err != nil {
		log.Printf("Failed to create task: %v\n", err.Error())
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	taskId, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"task_id": taskId,
	})
}

func (h *Handler) GetAllTasks(ctx *gin.Context) {
	userFromContext, _ := ctx.Get("user")
	user := userFromContext.(*model.User)

	c := ctx.Request.Context()

	tasks, err := h.services.Task.GetByUserId(c, user.Id)
	if err != nil {
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
	})
}

func (h *Handler) UpdateMultipleWithFile(ctx *gin.Context) {
	req := &dto.UpdateMultipleWithFileReq{}

	formFile, err := ctx.FormFile("file")
	if err != nil {
		log.Printf("Failed to get file from form: %v\n", err.Error())
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}
	req.File = formFile

	fileNameSlice := strings.Split(req.File.Filename, ".")

	if fileNameSlice[len(fileNameSlice)-1] != "pdf" {
		log.Printf("Failed to get file from form: unsupported media type.")
		err := apperrors.NewUnsupportedMediaType("Only pdf file format supported.")
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	ids := ctx.QueryMap("ids")
	idsSlice := make([]string, 0)
	for _, id := range ids {
		idsSlice = append(idsSlice, id)
	}

	req.Ids = idsSlice

	fileId, _ := uuid.NewRandom()

	req.File.Filename = fileId.String() + "_" + req.File.Filename

	c := ctx.Request.Context()

	err = h.services.Task.UpdateMultipleWithFile(c, req)
	if err != nil {
		log.Printf("Failed to update database with new filename: %v\n", err.Error())
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	dst := path.Join("./files/", req.File.Filename)
	err = ctx.SaveUploadedFile(req.File, dst)
	if err != nil {
		log.Printf("Failed save file: %v\n", err.Error())
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}
	ctx.String(http.StatusOK, "File uploaded successfully")
}

func (h *Handler) GetFile(ctx *gin.Context) {
	id := ctx.Param("id")

	c := ctx.Request.Context()

	fileName, err := h.services.Task.GetFileNameById(c, id)
	if err != nil {
		log.Printf("Failed to get task file: %v\n", err.Error())
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	targetPath := filepath.Join("./files/", fileName)

	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.File(targetPath)
}

func (h *Handler) OpenTask(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Printf("Failed to parse id from request: %v\n", err.Error())
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c := ctx.Request.Context()

	err = h.services.Task.Open(c, id)
	if err != nil {
		log.Printf("Failed to open task: %v\n", err.Error())
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) CloseTask(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Printf("Failed to parse id from request: %v\n", err.Error())
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c := ctx.Request.Context()

	err = h.services.Task.Close(c, id)
	if err != nil {
		log.Printf("Failed to close task: %v\n", err.Error())
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	ctx.Status(http.StatusOK)
}
