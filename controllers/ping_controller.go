package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	pong string = "pong"
)

var (
	PingController pingControllerInterface = &pingController{}
)

type pingControllerInterface interface {
	Ping(c *gin.Context)
}

type pingController struct{}

func (controller *pingController) Ping(c *gin.Context) {
	c.String(http.StatusOK, pong)
}