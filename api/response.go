package api

import "github.com/gin-gonic/gin"

type Response struct {
	Status string      `json:"status"`
	Code   int         `json:"code"`
	Data   interface{} `json:"data"`
}

func SendSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(code, Response{
		Status: "success",
		Code:   code,
		Data:   data,
	})
}

func SendError(c *gin.Context, code int, data interface{}) {
	c.JSON(code, Response{
		Status: "error",
		Code:   code,
		Data:   data,
	})
}
