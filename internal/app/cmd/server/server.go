package server

import (
	"github.com/jwma/jump-jump/internal/app/config"
	_ "github.com/jwma/jump-jump/internal/app/config"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/routers"
)

func setupDB() error {
	c := db.GetRedisClient()
	pong := c.Ping()
	return pong.Err()
}

func Run(addr ...string) error {
	err := setupDB()
	if err != nil {
		return err
	}

	err = config.SetupConfig()
	if err != nil {
		return err
	}

	router := routers.SetupRouter()
	err = router.Run(addr...)
	return err
}

func RunLanding(addr ...string) error {
	err := setupDB()
	if err != nil {
		return err
	}

	err = config.SetupConfig()
	if err != nil {
		return err
	}

	router := routers.SetupLandingRouter()
	err = router.Run(addr...)
	return err
}
