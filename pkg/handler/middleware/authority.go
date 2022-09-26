package middleware

import (
	"github.com/gin-gonic/gin"
	apperrors "homeworkdeliverysystem/errors"
	"homeworkdeliverysystem/model"
	"log"
	"strings"
)

func Authority(authorities ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userFromContext, exists := ctx.Get("user")

		if !exists {
			log.Printf("Unable get user from context: %v\n", ctx)
			err := apperrors.NewForbidden()
			ctx.JSON(err.Status(), gin.H{
				"error": err,
			})
			ctx.Abort()
		}

		userRole := userFromContext.(*model.User).Role

		roles := ""
		for _, authority := range authorities {
			roles += authority
		}

		if !strings.Contains(roles, userRole) {
			log.Printf("Request forbidden because user role {%s} not contained in roles: {%s}\n", userRole, roles)
			err := apperrors.NewForbidden()
			ctx.JSON(err.Status(), gin.H{
				"error": err,
			})
			ctx.Abort()
		}

		ctx.Next()
	}
}
