package config

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	ShortLinkNotFoundContentMode  = "content"
	ShortLinkNotFoundRedirectMode = "redirect"
	DefaultIdLength               = 6
	DefaultIdMinimumLength        = 2
	DefaultIdMaximumLength        = 10
)

type IdConfig struct {
	IdLength        int `json:"idLength" format:"int" example:"6"`
	IdMinimumLength int `json:"idMinimumLength" format:"int" example:"2"`
	IdMaximumLength int `json:"idMaximumLength" format:"int" example:"10"`
}

type ShortLinkNotFoundConfig struct {
	Mode  string `json:"mode" binding:"required" example:"content" enums:"content,redirect"`
	Value string `json:"value" binding:"required" example:"page not found"`
}

func (s *ShortLinkNotFoundConfig) ToMap() map[string]string {
	return map[string]string{
		"mode":  s.Mode,
		"value": s.Value,
	}
}

type SystemConfig struct {
	LandingHosts            []string                `json:"landingHosts" format:"array" example:"https://a.com/,https://b.com/"`
	IdConfig                *IdConfig               `json:"idConfig"`
	ShortLinkNotFoundConfig *ShortLinkNotFoundConfig `json:"shortLinkNotFoundConfig"`
}

type dbConfig struct {
	LandingHosts  []string
	IdLength      int
	IdMinLength   int
	IdMaxLength   int
	NotFoundMode  string
	NotFoundValue string
}

var (
	pool   *pgxpool.Pool
	cached *dbConfig
	mu     sync.RWMutex
)

func SetupConfig(p *pgxpool.Pool) error {
	pool = p
	return reload()
}

func reload() error {
	c := &dbConfig{}
	err := pool.QueryRow(context.Background(),
		`SELECT landing_hosts, id_length, id_min_length, id_max_length, not_found_mode, not_found_value
		 FROM system_configs WHERE id = 1`).Scan(
		&c.LandingHosts, &c.IdLength, &c.IdMinLength, &c.IdMaxLength,
		&c.NotFoundMode, &c.NotFoundValue)
	if err != nil {
		return err
	}

	mu.Lock()
	cached = c
	mu.Unlock()
	return nil
}

func GetIdConfig() *IdConfig {
	mu.RLock()
	defer mu.RUnlock()
	return &IdConfig{
		IdLength:        cached.IdLength,
		IdMinimumLength: cached.IdMinLength,
		IdMaximumLength: cached.IdMaxLength,
	}
}

func GetShortLinkNotFoundConfig() *ShortLinkNotFoundConfig {
	mu.RLock()
	defer mu.RUnlock()
	return &ShortLinkNotFoundConfig{
		Mode:  cached.NotFoundMode,
		Value: cached.NotFoundValue,
	}
}

func GetSystemConfig() *SystemConfig {
	return &SystemConfig{
		LandingHosts:            getLandingHosts(),
		IdConfig:                GetIdConfig(),
		ShortLinkNotFoundConfig: GetShortLinkNotFoundConfig(),
	}
}

func getLandingHosts() []string {
	mu.RLock()
	defer mu.RUnlock()
	return cached.LandingHosts
}

func UpdateLandingHosts(hosts []string) {
	pool.Exec(context.Background(),
		`UPDATE system_configs SET landing_hosts = $1 WHERE id = 1`, hosts)
	reload()
}

func UpdateIdConfig(c *IdConfig) {
	pool.Exec(context.Background(),
		`UPDATE system_configs SET id_min_length = $1, id_length = $2, id_max_length = $3 WHERE id = 1`,
		c.IdMinimumLength, c.IdLength, c.IdMaximumLength)
	reload()
}

func UpdateShortLinkNotFoundConfig(s *ShortLinkNotFoundConfig) {
	pool.Exec(context.Background(),
		`UPDATE system_configs SET not_found_mode = $1, not_found_value = $2 WHERE id = 1`,
		s.Mode, s.Value)
	reload()
}
