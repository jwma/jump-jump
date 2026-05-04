package server

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/config"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/routers"
	"github.com/jwma/jump-jump/internal/app/workers"
)

func allowHostsChecking() error {
	if gin.Mode() == gin.ReleaseMode {
		allowedHosts := os.Getenv("ALLOWED_HOSTS")
		if slices.Contains([]string{"", "*"}, allowedHosts) {
			return fmt.Errorf("please set ALLOWED_HOSTS environment variable when GIN_MODE=release")
		}
	}
	return nil
}

func setupDB() error {
	if err := db.InitRedis(); err != nil {
		return fmt.Errorf("Redis init failed: %w", err)
	}
	if err := db.InitPostgres(); err != nil {
		return fmt.Errorf("PostgreSQL init failed: %w", err)
	}
	if err := db.RunMigrations(); err != nil {
		return fmt.Errorf("Migration failed: %w", err)
	}
	return nil
}

func Run(addr ...string) error {
	if err := allowHostsChecking(); err != nil {
		return err
	}

	if err := setupDB(); err != nil {
		return err
	}
	defer db.ClosePostgres()
	defer db.CloseRedis()

	if err := config.SetupConfig(db.GetPostgresPool(), db.GetRedisClient()); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go workers.NewHistoryFlushWorker(db.GetRedisClient(), db.GetPostgresPool()).Run(ctx)

	router := routers.SetupRouter()
	return router.Run(addr...)
}

func RunLanding(addr ...string) error {
	if err := setupDB(); err != nil {
		return err
	}
	defer db.ClosePostgres()
	defer db.CloseRedis()

	if err := config.SetupConfig(db.GetPostgresPool(), db.GetRedisClient()); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go workers.NewHistoryFlushWorker(db.GetRedisClient(), db.GetPostgresPool()).Run(ctx)

	router := routers.SetupLandingRouter()
	return router.Run(addr...)
}
