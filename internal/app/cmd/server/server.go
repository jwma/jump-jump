package server

import (
	_ "github.com/jwma/jump-jump/internal/app/config"
	"github.com/jwma/jump-jump/internal/app/routers"
)

func Run(addr ...string) error {
	router := routers.SetupRouter()
	err := router.Run(addr...)
	return err
}

func RunLanding(addr ...string) error {
	router := routers.SetupLandingRouter()
	err := router.Run(addr...)
	return err
}
