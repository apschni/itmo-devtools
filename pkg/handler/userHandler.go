package handler

import (
	"github.com/gin-gonic/gin"
	apperrors "homeworkdeliverysystem/errors"
	"homeworkdeliverysystem/model"
	"log"
	"net/http"
)

func (h *Handler) Me(c *gin.Context) {
	user, exists := c.Get("user")

	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		err := apperrors.NewInternal()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})

		return
	}

	uid := user.(*model.User).Id

	ctx := c.Request.Context()
	u, err := h.services.User.Get(ctx, uid)

	if err != nil {
		log.Printf("Unable to find user: %v\n%v", uid, err)
		e := apperrors.NewNotFound("user", uid.String())

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": u,
	})
}
