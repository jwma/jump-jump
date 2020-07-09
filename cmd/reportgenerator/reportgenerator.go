package main

import (
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/report"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 每 30 秒运行一次报表生成/更新
	rg := report.NewGenerator(db.GetRedisClient(), time.Second*30)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		_ = rg.Stop() // 如果程序中断，停止 report generator
		os.Exit(1)
	}()

	err := rg.Start()

	if err != nil {
		panic(err)
	}
}
