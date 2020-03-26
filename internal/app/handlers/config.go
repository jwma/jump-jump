package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jwma/jump-jump/internal/app/config"
	"net/http"
)

func GetConfig(c *gin.Context) {
	cfg := config.GetConfig()
	landingHosts := cfg.GetStringSliceValue("landingHosts", make([]string, 0))

	c.JSON(http.StatusOK, gin.H{
		"msg":  "ok",
		"code": 0,
		"data": gin.H{
			"config": gin.H{"landingHosts": landingHosts},
		},
	})
}
