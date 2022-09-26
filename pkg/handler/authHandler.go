package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"homeworkdeliverysystem/dto"
	apperrors "homeworkdeliverysystem/errors"
	"homeworkdeliverysystem/model"
	"log"
	"net/http"
)

func (h *Handler) SignUp(ctx *gin.Context) {
	var req dto.SignUpReq

	if ok := bindData(ctx, &req); !ok {
		return
	}

	u := &model.User{
		FullName:    req.FullName,
		GroupNumber: req.GroupNumber,
		Username:    req.Username,
		Password:    req.Password,
		Role:        "user",
	}

	c := ctx.Request.Context()
	id, err := h.services.Authorization.SignUp(c, u)
	if err != nil {
		log.Printf("Failed to sign up user: %v\n", err.Error())
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}
	userId, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
	}
	u.Id = userId

	tokens, err := h.services.Token.NewPairFromUser(c, u, "")

	if err != nil {
		log.Printf("Failed to create tokens for user: %v\n", err.Error())
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"tokens": tokens,
	})
}

func (h *Handler) SignIn(ctx *gin.Context) {
	var req dto.SignInReq

	if ok := bindData(ctx, &req); !ok {
		return
	}

	user := &model.User{
		Username: req.Username,
		Password: req.Password,
	}

	c := ctx.Request.Context()
	err := h.services.Authorization.SignIn(c, user)

	if err != nil {
		log.Printf("Failed to sign in user: %v\n", err.Error())
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	tokens, err := h.services.Token.NewPairFromUser(ctx, user, "")

	if err != nil {
		log.Printf("Failed to create tokens for user: %v\n", err.Error())

		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})
}

func (h *Handler) SignOut(ctx *gin.Context) {
	user := ctx.MustGet("user")

	c := ctx.Request.Context()
	if err := h.services.Authorization.SignOut(c, user.(*model.User).Id); err != nil {
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user signed out successfully!",
	})
}
