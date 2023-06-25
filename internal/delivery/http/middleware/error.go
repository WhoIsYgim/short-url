package middleware

import (
	"github.com/gin-gonic/gin"
	"short-link/pkg/errs"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) == 0 {
			return
		}

		e := ctx.Errors[0].Unwrap()
		appErr, ok := e.(*errs.AppError)

		var err error
		if ok {
			err = appErr.Unwrap()
		} else {
			err = errs.InternalServerError
		}

		ctx.JSON(errs.Errors[err].Code, gin.H{
			"status":  errs.Errors[err].Code,
			"message": errs.Errors[err].Message,
		})
	}
}
