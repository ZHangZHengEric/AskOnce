package middleware

import (
	"dispatch/components"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func TaskTimeout(time time.Duration) gin.HandlerFunc {
	return timeout.New(timeout.WithHandler(func(c *gin.Context) { c.Next() }), timeout.WithTimeout(time), timeout.WithResponse(timeoutResponse))
}

func timeoutResponse(c *gin.Context) {
	c.JSON(http.StatusGatewayTimeout, components.ErrorTaskTimeOut)
}
