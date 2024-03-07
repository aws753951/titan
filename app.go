package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)


type Person struct {
	Username string `json:"username"`
	Create_time time.Time `json:"create_time"`
}

type Transcations struct {
	Borrow_fee float64 `json:borrow_fee`
	Create_time time.Time `json:create_time`
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

	// 初始化Gin引擎
	r := gin.Default()

	r.GET("/user", func(c *gin.Context) {

		// 連接到MySQL數據庫
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbIP, dbPort, dbName))
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		
		// 取得客戶資料
		query := fmt.Sprintf("SELECT username, create_time FROM member where member_pk=688")
		rows, err := db.Query(query)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// 定義一個切片來存放查詢結果
		var people []Person

		// 遍歷查詢結果，將每一行資料轉換為Person結構體，並添加到切片中
		for rows.Next() {
			var person Person
			err := rows.Scan(&person.Username, &person.Create_time)
			if err != nil {
				log.Fatal(err)
			}
			people = append(people, person)
		}

		query = "SELECT bf.borrow_fee, bf.create_time, td.times FROM trade_details td JOIN borrow_fee bf ON td.trade_pk = bf.trade_pk WHERE td.member_pk = 688;"
		rows, err = db.Query(query)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// 定義一個切片來存放查詢結果
		var transcations []Transcations

		for rows.Next() {
			var transcation Transcations
			err := rows.Scan(&transcation.Borrow_fee, &transcation.Create_time, &transcation.Times)
			if err != nil {
				log.Fatal(err)
			}
			transcations = append(transcations, transcation)
		}

		// 返回JSON格式的資料
		c.JSON(200, gin.H{
			"data": people,
			"trade": transcations,
		})
	})

	// 啟動伺服器，監聽預設端口 8080
	r.Run()


// 	// 連接MySQL數據庫
// 	db, err := sql.Open("mysql", "root:abc123456@tcp(127.0.0.1:3306)/test")
// 	if err != nil {
// 		fmt.Println("有問題", err)
// 		return
// 	}
// 	defer db.Close()

// 	// 確保與數據庫的連接是有效的
// 	err = db.Ping()
// 	if err != nil {
// 		fmt.Println("ping有問題", err)
// 		return
// 	}
// 	fmt.Println("Connected to the MySQL database")


// 	// // 創建表格
// 	// _, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
// 	// 	id INT AUTO_INCREMENT PRIMARY KEY,
// 	// 	name VARCHAR(50),
// 	// 	age INT
// 	// )`)
	
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// 	return
// 	// }
// 	// fmt.Println("Table 'users' created successfully")

// 	// 插入資料
// 	insertQuery := "INSERT INTO users (name, age) VALUES (?, ?)"
// 	insertStmt, err := db.Prepare(insertQuery)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer insertStmt.Close()

// 	_, err = insertStmt.Exec("Alice", 30)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println("Data inserted successfully")

// 	// 查詢資料
// 	rows, err := db.Query("SELECT id, name, age FROM users")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer rows.Close()

// 	fmt.Println("Users:")
// 	for rows.Next() {
// 		var id int
// 		var name string
// 		var age int
// 		err := rows.Scan(&id, &name, &age)
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 		fmt.Printf("ID: %d, Name: %s, Age: %d\n", id, name, age)
// 	}

// 	// 更新資料
// 	updateQuery := "UPDATE users SET age = ? WHERE name = ?"
// 	updateStmt, err := db.Prepare(updateQuery)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer updateStmt.Close()

// 	_, err = updateStmt.Exec(35, "Alice")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println("Data updated successfully")

// 	// 刪除資料
// 	deleteQuery := "DELETE FROM users WHERE name = ?"
// 	deleteStmt, err := db.Prepare(deleteQuery)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer deleteStmt.Close()

// 	_, err = deleteStmt.Exec("Alice")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println("Data deleted successfully")




}
