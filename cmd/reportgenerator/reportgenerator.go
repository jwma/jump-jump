package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/report"
)

func main() {
	if err := db.InitRedis(); err != nil {
		panic(err)
	}
	defer db.CloseRedis()

	if err := db.InitPostgres(); err != nil {
		panic(err)
	}
	defer db.ClosePostgres()

	rg := report.NewGenerator(db.GetRedisClient(), db.GetPostgresPool(), time.Second*30)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		_ = rg.Stop()
		os.Exit(1)
	}()

	if err := rg.Start(); err != nil {
		panic(err)
	}
}
