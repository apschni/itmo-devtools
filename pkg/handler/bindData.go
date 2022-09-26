package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	apperrors "homeworkdeliverysystem/errors"
	"log"
)

type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

func bindData(ctx *gin.Context, req interface{}) bool {
	if ctx.ContentType() != "application/json" {
		msg := fmt.Sprintf("%s only accepts Content-Type application/json", ctx.FullPath())

		err := apperrors.NewUnsupportedMediaType(msg)

		ctx.JSON(err.Status(), gin.H{
			"error": err,
		})
		return false
	}

	if err := ctx.ShouldBind(req); err != nil {
		log.Printf("Error binding data: %+v\n", err)

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
			return false
		}

		fallBack := apperrors.NewInternal()

		ctx.JSON(fallBack.Status(), gin.H{"error": fallBack})
		return false
	}

	return true
}
