package handler

import (
	"github.com/gin-gonic/gin"
	"homeworkdeliverysystem/dto"
	apperrors "homeworkdeliverysystem/errors"
	"log"
	"net/http"
)

func (h *Handler) Tokens(ctx *gin.Context) {
	var req dto.TokensReq

	if ok := bindData(ctx, &req); !ok {
		return
	}

	c := ctx.Request.Context()

	refreshToken, err := h.services.Token.ValidateRefreshToken(req.RefreshToken)

	if err != nil {
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	user, err := h.services.User.Get(c, refreshToken.UID)

	if err != nil {
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	tokens, err := h.services.Token.NewPairFromUser(c, user, refreshToken.ID.String())

	if err != nil {
		log.Printf("Failed to create tokens for user: %+v. Error: %v\n", user, err.Error())

		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})
}
