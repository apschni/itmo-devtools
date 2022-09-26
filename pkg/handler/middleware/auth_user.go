package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	apperrors "homeworkdeliverysystem/errors"
	"homeworkdeliverysystem/pkg/service"
	"strings"
)

type authHeader struct {
	IDToken string `header:"Authorization"`
}

type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

func AuthUser(s service.Token) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h := authHeader{}

		if err := ctx.ShouldBindHeader(&h); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {
				var invalidArgs []invalidArgument

				for _, err := range errs {
					invalidArgs = append(invalidArgs, invalidArgument{
						err.Field(),
						err.Value().(string),
						err.Tag(),
						err.Param(),
					})
				}

				err := apperrors.NewBadRequest("Invalid request parameters. See invalidArgs")

				ctx.JSON(err.Status(), gin.H{
					"error":       err,
					"invalidArgs": invalidArgs,
				})
				ctx.Abort()
				return
			}

			err := apperrors.NewInternal()
			ctx.JSON(err.Status(), gin.H{
				"error": err,
			})
			ctx.Abort()
			return
		}

		idTokenHeader := strings.Split(h.IDToken, "Bearer ")

		if len(idTokenHeader) < 2 {
			err := apperrors.NewAuthorization("Must provide Authorization header with format `Bearer {token}`")

			ctx.JSON(err.Status(), gin.H{
				"error": err,
			})
			ctx.Abort()
			return
		}

		user, err := s.ValidateIdToken(idTokenHeader[1])

		if err != nil {
			err := apperrors.NewAuthorization("Provided token is invalid")
			ctx.JSON(err.Status(), gin.H{
				"error": err,
			})
			ctx.Abort()
			return
		}

		ctx.Set("user", user)

		ctx.Next()
	}
}
