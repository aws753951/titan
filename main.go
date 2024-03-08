package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)


type Person struct {
	Member_pk int `json:"member_pk"`
	Username string `json:"username"`
	Create_time time.Time `json:"create_time"`
}

type Transcations struct {
	Borrow_fee float64 `json:borrow_fee`
	Create_time time.Time `json:create_time`
	Type int `json:type`
	Times int `json:times` 
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbIP := os.Getenv("DB_IP")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	r := gin.Default()
	r.Use(corsMiddleware())

	r.GET("/user", func(c *gin.Context) {

		// 連接到MySQL數據庫
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbIP, dbPort, dbName))
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		username := c.Query("username")
		starttime := c.Query("starttime")
		endtime := c.Query("endtime")

		if starttime == "" {
			starttime = "DATE_SUB(NOW(), INTERVAL 1 YEAR)"
		}
		if endtime == "" {
			endtime = "NOW()"
		}


		// 取得客戶資料
		query := fmt.Sprintf("SELECT member_pk, username, create_time FROM member where username= '%s'", username)
		rows, err := db.Query(query)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// 定義一個切片來存放查詢結果
		var person Person
		// 遍歷查詢結果，將每一行資料轉換為Person結構體，並添加到切片中
		for rows.Next() {
			err := rows.Scan(&person.Member_pk, &person.Username, &person.Create_time)
			if err != nil {
				log.Fatal(err)
			}
			break
		}

		query = fmt.Sprintf("SELECT bf.borrow_fee, bf.create_time, td.type, td.times FROM trade_details td JOIN borrow_fee bf on td.trade_pk = bf.trade_pk WHERE td.member_pk = %d and create_time between %s and %s", person.Member_pk, starttime, endtime)
		rows, err = db.Query(query)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var transcations []Transcations

		for rows.Next() {
			var transcation Transcations
			err := rows.Scan(&transcation.Borrow_fee, &transcation.Create_time, &transcation.Type, &transcation.Times)
			if err != nil {
				log.Fatal(err)
			}
			transcations = append(transcations, transcation)
		}

		// 返回JSON格式的資料
		c.JSON(200, gin.H{
			"data": person,
			"trade": transcations,
		})
	})



	r.PUT("/user", func(c *gin.Context) {

		// 連接到MySQL數據庫
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbIP, dbPort, dbName))
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		var person Person
		if err := c.ShouldBindJSON(&person); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		_, err = db.Exec("UPDATE member SET username = ? where member_pk = ?", person.Username, person.Member_pk)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 返回成功响应，包含更新后的数据状态
		c.JSON(http.StatusOK, person)

	})


	r.POST("/user", func(c *gin.Context) {
		// 連接到MySQL數據庫
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbIP, dbPort, dbName))
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		var person Person
		if err := c.Bind(&person); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		person.Create_time = time.Now()

		result, err := db.Exec("INSERT INTO member (username, create_time) VALUES (?, ?)", person.Username, person.Create_time.Format("2006-01-02 15:04:05"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		lastInsertID, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"member_pk":   lastInsertID,
			"username":    person.Username,
			"create_time": person.Create_time,
		})
	})

	// 啟動伺服器，監聽預設端口 8080
	r.Run()
}


func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许的来源
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// 允许的请求方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT")
		// 允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// 如果是预检请求（OPTIONS 请求），直接返回
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		// 继续执行后续中间件和处理程序
		c.Next()
	}
}