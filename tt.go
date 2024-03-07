package main

import (
	"fmt"
	"math/rand"
	"time"
)

func generateRandomFloat() float64 {
	// 設定隨機種子
	rand.Seed(time.Now().UnixNano())

	// 產生隨機浮點數
	return float64(rand.Intn(100000000000000)) / 100
}

func main() {
	// 呼叫函式並輸出結果
	fmt.Println(generateRandomFloat())
}