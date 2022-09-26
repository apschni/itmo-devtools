package handler

import (
	"github.com/gin-gonic/gin"
	"homeworkdeliverysystem/dto"
	apperrors "homeworkdeliverysystem/errors"
	"log"
	"net/http"
)

func (h *Handler) GetSubjects(ctx *gin.Context) {
	var req dto.GetSubjectsReq

	err := ctx.ShouldBindUri(&req)

	if err != nil {
		log.Printf("Failed bind uri path params: %v\n", err.Error())
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c := ctx.Request.Context()

	subjects, err := h.services.Group.GetSubjectsByNumber(c, req.Number)
	if err != nil {
		log.Printf("Failed to get group subjects: %v\n", err.Error())
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"subjects": subjects,
	})
}

func (h *Handler) GetStudents(ctx *gin.Context) {
	var req dto.GetSubjectsReq

	err := ctx.ShouldBindUri(&req)

	if err != nil {
		log.Printf("Failed bind uri path params: %v\n", err.Error())
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c := ctx.Request.Context()

	students, err2 := h.services.Group.GetStudentsByNumber(c, req.Number)
	if err2 != nil {
		log.Printf("Failed to get group students: %v\n", err2.Error())
		ctx.JSON(apperrors.Status(err2), gin.H{
			"error": err2,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"students": students,
	})
}
