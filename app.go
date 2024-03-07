package main



FOREIGN KEY (trade_pk) REFERENCES borrow_fee(trade_pk),
FOREIGN KEY (member_pk) REFERENCES member(member_pk)



























import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}


	// 初始化Gin引擎
	r := gin.Default()

	// 設置路由
	r.GET("/", func(c *gin.Context) {

		// 連接到MySQL數據庫
		db, err := sql.Open("mysql", "root:"+os.Getenv("DB_PASSWORD") +"@tcp("+ os.Getenv("dbHost")+":3306)/test")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// 查詢資料庫取得資料
		rows, err := db.Query("SELECT name, age FROM users")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// 定義一個切片來存放查詢結果
		var people []Person

		// 遍歷查詢結果，將每一行資料轉換為Person結構體，並添加到切片中
		for rows.Next() {
			var person Person
			err := rows.Scan(&person.Name, &person.Age)
			if err != nil {
				log.Fatal(err)
			}
			people = append(people, person)
		}

		// 返回JSON格式的資料
		c.JSON(200, gin.H{
			"data": people,
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
