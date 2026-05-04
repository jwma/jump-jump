package config

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/jwma/jump-jump/internal/app/utils"
	"time"
)

const (
	ShortLinkNotFoundContentMode  = "content"
	ShortLinkNotFoundRedirectMode = "redirect"
	DefaultIdLength               = 6
	DefaultIdMinimumLength        = 2
	DefaultIdMaximumLength        = 10
)

type IdConfig struct {
	IdLength        int `json:"idLength"`
	IdMinimumLength int `json:"idMinimumLength"`
	IdMaximumLength int `json:"idMaximumLength"`
}

type ShortLinkNotFoundConfig struct {
	Mode  string `json:"mode"`
	Value string `json:"value"`
}

type TenantConfig struct {
	IdConfig                *IdConfig                `json:"idConfig"`
	ShortLinkNotFoundConfig *ShortLinkNotFoundConfig `json:"shortLinkNotFoundConfig"`
}

var (
	pool *pgxpool.Pool
	rdb  *redis.Client
	mu   sync.RWMutex
	cache map[string]*TenantConfig
)

func SetupConfig(p *pgxpool.Pool, r *redis.Client) error {
	pool = p
	rdb = r
	cache = make(map[string]*TenantConfig)
	return nil
}

func GetTenantConfig(tenantID string) *TenantConfig {
	mu.RLock()
	if c, ok := cache[tenantID]; ok {
		mu.RUnlock()
		return c
	}
	mu.RUnlock()

	// Cache miss, load from PG
	c := loadTenantConfig(tenantID)

	mu.Lock()
	cache[tenantID] = c
	mu.Unlock()
	return c
}

func loadTenantConfig(tenantID string) *TenantConfig {
	c := &TenantConfig{
		IdConfig: &IdConfig{
			IdLength:        DefaultIdLength,
			IdMinimumLength: DefaultIdMinimumLength,
			IdMaximumLength: DefaultIdMaximumLength,
		},
		ShortLinkNotFoundConfig: &ShortLinkNotFoundConfig{
			Mode:  ShortLinkNotFoundContentMode,
			Value: "你访问的页面不存在哦",
		},
	}

	_ = pool.QueryRow(context.Background(),
		`SELECT id_length, id_min_length, id_max_length, not_found_mode, not_found_value
		 FROM tenant_configs WHERE tenant_id = $1`, tenantID).Scan(
		&c.IdConfig.IdLength, &c.IdConfig.IdMinimumLength, &c.IdConfig.IdMaximumLength,
		&c.ShortLinkNotFoundConfig.Mode, &c.ShortLinkNotFoundConfig.Value)

	return c
}

func GetIdConfig(tenantID string) *IdConfig {
	return GetTenantConfig(tenantID).IdConfig
}

func GetShortLinkNotFoundConfig(tenantID string) *ShortLinkNotFoundConfig {
	return GetTenantConfig(tenantID).ShortLinkNotFoundConfig
}

func UpdateIdConfig(tenantID string, c *IdConfig) {
	pool.Exec(context.Background(),
		`INSERT INTO tenant_configs (tenant_id, id_length, id_min_length, id_max_length)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (tenant_id) DO UPDATE SET id_length = $2, id_min_length = $3, id_max_length = $4, updated_at = now()`,
		tenantID, c.IdLength, c.IdMinimumLength, c.IdMaximumLength)
	invalidateCache(tenantID)
}

func UpdateShortLinkNotFoundConfig(tenantID string, s *ShortLinkNotFoundConfig) {
	pool.Exec(context.Background(),
		`INSERT INTO tenant_configs (tenant_id, not_found_mode, not_found_value)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (tenant_id) DO UPDATE SET not_found_mode = $2, not_found_value = $3, updated_at = now()`,
		tenantID, s.Mode, s.Value)
	invalidateCache(tenantID)
}

func invalidateCache(tenantID string) {
	mu.Lock()
	delete(cache, tenantID)
	mu.Unlock()
	if rdb != nil {
		rdb.Del(context.Background(), utils.GetTenantConfigCacheKey(tenantID))
	}
}

// ResolveTenantID maps a domain to a tenant_id, using Redis cache
func ResolveTenantID(domain string) (string, error) {
	// Try cache
	if rdb != nil {
		val, err := rdb.Get(context.Background(), utils.GetDomainCacheKey(domain)).Result()
		if err == nil && val != "" {
			return val, nil
		}
	}

	// Query PG
	var tenantID string
	err := pool.QueryRow(context.Background(),
		`SELECT tenant_id FROM tenant_domains WHERE domain = $1`, domain).Scan(&tenantID)
	if err != nil {
		return "", err
	}

	// Cache for 30 min
	if rdb != nil {
		rdb.Set(context.Background(), utils.GetDomainCacheKey(domain), tenantID, 30*time.Minute)
	}

	return tenantID, nil
}
