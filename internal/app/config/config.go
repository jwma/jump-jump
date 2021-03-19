package config

import (
	"github.com/go-redis/redis"
	"github.com/jwma/reborn"
	"time"
)

var config *reborn.Reborn

const (
	ShortLinkNotFoundContentMode  = "content"
	ShortLinkNotFoundRedirectMode = "redirect"
)

type IdConfig struct {
	// ID 长度
	IdLength int `json:"idLength" format:"int" example:"6"`

	// 最小 ID 长度
	IdMinimumLength int `json:"idMinimumLength" format:"int" example:"2"`

	// 最大 ID 长度
	IdMaximumLength int `json:"idMaximumLength" format:"int" example:"10"`
} // @name IdConfig

type ShortLinkNotFoundConfig struct {
	// 模式
	Mode string `json:"mode" binding:"required" example:"content" enums:"content,redirect"`

	// 值
	Value string `json:"value" binding:"required" example:"page not found"`
} // @name ShortLinkNotFoundConfig

func (s *ShortLinkNotFoundConfig) ToMap() map[string]string {
	return map[string]string{
		"mode":  s.Mode,
		"value": s.Value,
	}
}

type SystemConfig struct {
	// 落地页 Host 列表
	LandingHosts []string `json:"landingHosts" format:"array" example:"https://a.com/,https://b.com/"`

	// ID 配置
	IdConfig *IdConfig `json:"idConfig"`

	// 短链接 404 配置
	ShortLinkNotFoundConfig *ShortLinkNotFoundConfig `json:"shortLinkNotFoundConfig"`
} // @name SystemConfig

func GetIdConfig() *IdConfig {
	return &IdConfig{
		IdLength:        config.GetIntValue("idLength", 6),
		IdMinimumLength: config.GetIntValue("idMinimumLength", 2),
		IdMaximumLength: config.GetIntValue("idMaximumLength", 10),
	}
}

func GetDefaultShortLinkNotFoundConfig() map[string]string {
	return map[string]string{
		"mode":  ShortLinkNotFoundContentMode,
		"value": "你访问的页面不存在哦",
	}
}

func GetShortLinkNotFoundConfig() *ShortLinkNotFoundConfig {
	c := config.GetStringStringMapValue("shortLinkNotFoundConfig", GetDefaultShortLinkNotFoundConfig())

	return &ShortLinkNotFoundConfig{
		Mode:  c["mode"],
		Value: c["value"],
	}
}

func GetSystemConfig() *SystemConfig {
	return &SystemConfig{
		LandingHosts:            config.GetStringSliceValue("landingHosts", make([]string, 0)),
		IdConfig:                GetIdConfig(),
		ShortLinkNotFoundConfig: GetShortLinkNotFoundConfig(),
	}
}

func UpdateLandingHosts(hosts []string) {
	config.SetValue("landingHosts", hosts)
	config.Persist()
}

func UpdateIdConfig(c *IdConfig) {
	config.SetValue("idMinimumLength", c.IdMinimumLength)
	config.SetValue("idLength", c.IdLength)
	config.SetValue("idMaximumLength", c.IdMaximumLength)
	config.Persist()
}

func UpdateShortLinkNotFoundConfig(s *ShortLinkNotFoundConfig) {
	config.SetValue("shortLinkNotFoundConfig", s.ToMap())
	config.Persist()
}

func getDefaultConfig() *reborn.Config {
	d := reborn.NewConfig()
	d.SetValue("landingHosts", []string{"http://127.0.0.1:8081/"})
	d.SetValue("idMinimumLength", 2)
	d.SetValue("idLength", 6)
	d.SetValue("idMaximumLength", 10)
	d.SetValue("shortLinkNotFoundConfig", GetDefaultShortLinkNotFoundConfig())

	return d
}

func GetConfig() *reborn.Reborn {
	return config
}

func SetupConfig(rdb *redis.Client) error {
	var err error
	config, err = reborn.NewWithDefaults(rdb, "j2config", getDefaultConfig())
	if err != nil {
		return err
	}
	config.SetAutoReloadDuration(time.Second * 30)
	config.StartAutoReload()

	return nil
}
