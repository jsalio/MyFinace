package controllers

import "github.com/gin-gonic/gin"

type Controller interface {
	RegisterRoutes(router *gin.RouterGroup)
}

type BaseController struct {
	Path string
}

func NewBaseController(path string) *BaseController {
	return &BaseController{
		Path: path,
	}
}
