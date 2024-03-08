package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
	"sort"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Member struct {
	Pk         int
	Username   string
	CreateTime time.Time
}

func main() {
	// 讀取環境變數
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("無法讀取 .env 檔案: %v", err)
	}

	// 取得 MySQL 連線資訊
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbIP := os.Getenv("DB_IP")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// 連線到 MySQL 資料庫
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/", dbUser, dbPass, dbIP, dbPort))
	if err != nil {
		log.Fatalf("無法連線到 MySQL 資料庫: %v", err)
	}
	defer db.Close()

    _, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
    if err != nil {
        log.Fatal(err)
    }

	_, err = db.Exec(fmt.Sprintf("USE %s", dbName))
    if err != nil {
        log.Fatal(err)
    }


	// 建立 member 表格
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS member (
			member_pk int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '用戶pk',
			username varchar(16) NOT NULL COMMENT '登入帳號',
			create_time datetime NOT NULL COMMENT '註冊時間',
			PRIMARY KEY (member_pk)
		)
	`)
	if err != nil {
		log.Fatalf("建立表格失敗: %v", err)
	}

	// 建立 borrow_fee 表格
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS borrow_fee (
			trade_pk int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '交易pk',
			borrow_fee decimal(12,2) NOT NULL COMMENT '交易金額',
			create_time datetime NOT NULL COMMENT '發生時間',
			PRIMARY KEY (trade_pk)
		)
	`)
	if err != nil {
		log.Fatalf("建立表格失敗: %v", err)
	}

	// 建立 trade_details 表格
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS trade_details (
			trade_pk int(10) unsigned NOT NULL COMMENT '交易pk',
			member_pk int(10) unsigned NOT NULL COMMENT '用戶pk',
			type int NOT NULL COMMENT '1代表收款方，2代表付款方',
			times int NOT NULL COMMENT '該用戶第幾次交易',
			PRIMARY KEY (trade_pk, member_pk)
		)
	`)
	if err != nil {
		log.Fatalf("建立表格失敗: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("無法開始事務: %v", err)
	}
	
	// 產生並插入客戶資料
	for i := 0; i < 1000; i++ {
		member := generateRandomMember()
		_, err := tx.Exec("INSERT INTO member (username, create_time) VALUES (?, ?)", member.Username, member.CreateTime)
		if err != nil {
			tx.Rollback()
			log.Fatalf("插入數據失敗: %v", err)
		}
	}
	// // 提交事務
	// if err = tx.Commit(); err != nil {
	// 	tx.Rollback()
	// 	log.Fatalf("提交事務失敗: %v", err)
	// }
	
	// fmt.Println("客戶數據插入成功")

	// tx, err = db.Begin()

	var timeArray [] time.Time
	for i := 0; i < 5000; {
		randomTime := generateRandomTime()
		query := fmt.Sprintf("SELECT count(*) FROM member WHERE create_time < '%s'", randomTime.Format("2006-01-02 15:04:05"))

		var rowCount int
		err := tx.QueryRow(query).Scan(&rowCount)
		if err != nil {
			tx.Rollback()
			log.Fatalf("查詢member失敗: %v", err)
		}
		
		if rowCount < 2{
			fmt.Println("不足兩筆資料，重做", randomTime)
			continue
		}

		timeArray = append(timeArray, randomTime)
		i++
	}


	sort.Slice(timeArray, func(i, j int) bool {
		return timeArray[i].Before(timeArray[j])
	})

	for i := 0; i < len(timeArray); i++ {
		query := fmt.Sprintf("SELECT member_pk FROM member WHERE create_time < '%s' ORDER BY RAND() LIMIT 2", timeArray[i].Format("2006-01-02 15:04:05"))
		rows, err := tx.Query(query)
		if err != nil {
			tx.Rollback()
			log.Fatalf("查詢member失敗: %v", err)
		}
		defer rows.Close()

		var member_pks []int

		for rows.Next() {
			var member_pk int
			if err := rows.Scan(&member_pk); err != nil {
				log.Fatal(err)
			}
			member_pks = append(member_pks, member_pk)
		}

		randomBorrow_fee := generateRandomBorrow_fee()
		_, err = tx.Exec("INSERT INTO borrow_fee (borrow_fee, create_time) VALUES (?, ?)", randomBorrow_fee, timeArray[i])
		if err != nil {
			tx.Rollback()
			log.Fatalf("插入borrow_fee數據失敗: %v", err)
		}

		var trade_pk int
		err = tx.QueryRow("SELECT MAX(trade_pk) FROM borrow_fee").Scan(&trade_pk)
		if err != nil {
			tx.Rollback()
			log.Fatalf("查詢borrow_fee數據失敗: %v", err)
		}

		for index, member_pk := range member_pks {
			var times sql.NullInt64
			query := "SELECT MAX(times) FROM trade_details WHERE member_pk = ?"
			err = tx.QueryRow(query, member_pk).Scan(&times)
			if err != nil {
				tx.Rollback()
				log.Fatalf("查詢次數數據失敗: %v", err)
			}

			var maxTimes int
			if times.Valid {
				maxTimes = int(times.Int64)
			} else {
				maxTimes = 0
			}

			_, err := tx.Exec("INSERT INTO trade_details (trade_pk, member_pk, type, times) VALUES (?, ?, ?, ?)", trade_pk, member_pk, index + 1, maxTimes + 1)
			if err != nil {
				tx.Rollback()
				log.Fatalf("插入trade_details數據失敗: %v", err)
			}
		}

		i++ 

	}

	// 提交事務
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		log.Fatalf("提交事務失敗: %v", err)
	}

	fmt.Println("交易數據插入成功")
}

// 產生隨機會員資料
func generateRandomMember() Member {
	username := generateRandomString(16)
	createTime := generateRandomTime()
	return Member{Username: username, CreateTime: createTime}
}

// 產生指定長度的隨機字串
func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result string
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charSet))
		result += string(charSet[randomIndex])
	}
	return result
}

// 產生隨機時間
func generateRandomTime() time.Time {
	rand.Seed(time.Now().UnixNano())
	min := time.Now().AddDate(0, -18, 0).Unix()
	max := time.Now().Unix()
	delta := max - min
	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

// 產生隨機交易金額
func generateRandomBorrow_fee() float64 {
	rand.Seed(time.Now().UnixNano())
	
	return float64(rand.Intn(1000000000000)) / 100
}