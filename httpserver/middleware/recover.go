package middleware

import (
	"fmt"
	"github.com/200Lab-Education/go-sdk/sdkcm"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Recover(sc ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := sc.Logger("service")
		lvLogger := logger.GetLevel()

		defer func() {
			if err := recover(); err != nil {
				c.Header("Content-Type", "application/json")
				didFireError := false

				if appErr, ok := err.(sdkcm.AppError); ok {
					appErr.RootCause = appErr.RootError()

					if appErr.RootCause != nil {
						appErr.Log = appErr.RootCause.Error()
					}

					logger.Errorln("App Error: ", appErr)

					c.AbortWithStatusJSON(appErr.StatusCode, appErr)

					if lvLogger == logrus.TraceLevel.String() {
						panic(err)
					}
				} else {
					var appErr sdkcm.AppError

					if e, ok := err.(error); ok {
						appErr = sdkcm.AppError{StatusCode: http.StatusInternalServerError, Message: "internal server error"}
						logger.Errorln(e.Error())

						c.AbortWithStatusJSON(appErr.StatusCode, appErr)
						didFireError = true
						panic(err)
					} else {
						appErr = sdkcm.AppError{StatusCode: http.StatusInternalServerError, Message: fmt.Sprintf("%s", err)}
						logger.Errorln(fmt.Sprintf("%s", err))

						c.AbortWithStatusJSON(appErr.StatusCode, appErr)
						didFireError = true
						panic(err)
					}
				}

				if lvLogger == logrus.TraceLevel.String() && !didFireError {
					panic(err)
				}
			}
		}()

		c.Next()
	}
}
