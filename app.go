package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// 初始化Gin引擎
	router := gin.Default()

	// 設置路由
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	// 啟動伺服器，監聽預設端口 8080
	router.Run()
}