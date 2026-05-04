package db

import (
	"context"
	"log"
)

var migrationStmts = []string{
	`CREATE TABLE IF NOT EXISTS tenants (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL UNIQUE,
    slug        VARCHAR(50) NOT NULL UNIQUE,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
)`,

	`CREATE TABLE IF NOT EXISTS tenant_domains (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    domain      VARCHAR(255) NOT NULL UNIQUE,
    is_default  BOOLEAN NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
)`,

	`CREATE TABLE IF NOT EXISTS tenant_configs (
    tenant_id       UUID PRIMARY KEY REFERENCES tenants(id) ON DELETE CASCADE,
    id_length       SMALLINT NOT NULL DEFAULT 6,
    id_min_length   SMALLINT NOT NULL DEFAULT 2,
    id_max_length   SMALLINT NOT NULL DEFAULT 10,
    not_found_mode  VARCHAR(20) NOT NULL DEFAULT 'content',
    not_found_value TEXT NOT NULL DEFAULT '你访问的页面不存在哦',
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
)`,

	`CREATE TABLE IF NOT EXISTS users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    username    VARCHAR(50) NOT NULL,
    password    BYTEA NOT NULL,
    salt        BYTEA NOT NULL,
    role        SMALLINT NOT NULL DEFAULT 1,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tenant_id, username)
)`,

	`CREATE TABLE IF NOT EXISTS short_links (
    id          VARCHAR(20) PRIMARY KEY,
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    url         TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    is_enabled  BOOLEAN NOT NULL DEFAULT true,
    created_by  VARCHAR(50) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
)`,

	`CREATE INDEX IF NOT EXISTS idx_short_links_tenant ON short_links(tenant_id)`,
	`CREATE INDEX IF NOT EXISTS idx_short_links_tenant_created ON short_links(tenant_id, created_at DESC)`,
	`CREATE INDEX IF NOT EXISTS idx_short_links_created_by ON short_links(created_by)`,
	`CREATE INDEX IF NOT EXISTS idx_tenant_domains_domain ON tenant_domains(domain)`,

	`CREATE TABLE IF NOT EXISTS request_histories (
    id            BIGSERIAL PRIMARY KEY,
    short_link_id VARCHAR(20) NOT NULL REFERENCES short_links(id) ON DELETE CASCADE,
    tenant_id     UUID NOT NULL,
    url           TEXT NOT NULL DEFAULT '',
    ip            VARCHAR(45) NOT NULL DEFAULT '',
    ua            TEXT NOT NULL DEFAULT '',
    os            VARCHAR(50) NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
)`,

	`CREATE INDEX IF NOT EXISTS idx_request_histories_link_time ON request_histories(short_link_id, created_at DESC)`,
	`CREATE INDEX IF NOT EXISTS idx_request_histories_tenant ON request_histories(tenant_id)`,

	`CREATE TABLE IF NOT EXISTS user_preferences (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    key     VARCHAR(100) NOT NULL,
    value   TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (user_id, key)
)`,
}

func RunMigrations() error {
	for _, stmt := range migrationStmts {
		_, err := GetPostgresPool().Exec(context.Background(), stmt)
		if err != nil {
			log.Printf("Migration statement failed: %v\nStatement: %s", err, stmt)
			return err
		}
	}

	// Ensure a default tenant exists for development/testing
	_, err := GetPostgresPool().Exec(context.Background(), `
		INSERT INTO tenants (id, name, slug)
		VALUES ('00000000-0000-0000-0000-000000000001', 'Default', 'default')
		ON CONFLICT DO NOTHING`)
	if err != nil {
		log.Printf("Default tenant creation skipped: %v", err)
	}

	// Ensure default tenant has a config row
	_, err = GetPostgresPool().Exec(context.Background(), `
		INSERT INTO tenant_configs (tenant_id)
		VALUES ('00000000-0000-0000-0000-000000000001')
		ON CONFLICT DO NOTHING`)
	if err != nil {
		log.Printf("Default tenant config creation skipped: %v", err)
	}

	// Ensure default tenant has localhost domain for local development
	_, err = GetPostgresPool().Exec(context.Background(), `
		INSERT INTO tenant_domains (tenant_id, domain, is_default)
		VALUES ('00000000-0000-0000-0000-000000000001', 'localhost', true)
		ON CONFLICT DO NOTHING`)
	if err != nil {
		log.Printf("Default tenant domain creation skipped: %v", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}
