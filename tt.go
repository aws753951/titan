package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Username interface{} `form:"username" json:"username"`
}

func main() {
	r := gin.Default()

	r.POST("/user", func(c *gin.Context) {
		// 定义一个结构体用于绑定请求主体数据
		var user User

		// 使用 Bind 方法绑定请求主体数据到结构体上
		if err := c.Bind(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 在这里处理 user.Username
		// 这里仅打印 username
		fmt.Println("Received username:", user.Username)

		c.JSON(http.StatusOK, gin.H{
			"message": user.Username,
		})
	})

	r.Run(":8080")
}