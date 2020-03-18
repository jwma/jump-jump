package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// 获取指定 URL Query 参数的值（整数类型）
func GetIntQueryValue(c *gin.Context, key string, defaults int) int {
	if c.Query(key) == "" {
		return defaults
	}
	v, err := strconv.Atoi(c.Query(key))
	if err != nil {
		v = defaults
	}
	return v
}
