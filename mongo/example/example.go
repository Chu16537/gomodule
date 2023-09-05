package example

import (
	"context"
	"fmt"
	"gomodule/mongo"
	"time"
)

// 範例
func Example() {
	// 创建一个上下文和配置
	ctx := context.Background()
	cfg := mongo.Config{
		Addr:     "mongodb://localhost:27017", // MongoDB 地址
		Database: "mydatabase",                // 数据库名称
	}

	var dbHandler *mongo.Handler
	for i := 0; i < 10; i++ {
		// 初始化 MongoDB 连接处理程序
		mgo, err := mongo.New(ctx, cfg)
		if err != nil {
			fmt.Println("Failed to initialize MongoDB:", err)
			continue
		}
		dbHandler = mgo
	}

	if dbHandler == nil {
		panic("mongo New fail")
	}

	defer dbHandler.Done() // 在程序结束时关闭连接

	go check(ctx, dbHandler)

	// 为了演示，这里只是等待一段时间然后退出
	fmt.Println("Running your application...")
	time.Sleep(10 * time.Second)
	fmt.Println("Application finished.")
}

func check(ctx context.Context, mgo *mongo.Handler) {
	// 設定檢查秒數
	checkInterval := 2 * time.Second
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return // 上下文已完成，退出 goroutine
		case <-ticker.C:
			if err := mgo.Check(); err != nil {
				fmt.Println(err)
			}
		}
	}
}
