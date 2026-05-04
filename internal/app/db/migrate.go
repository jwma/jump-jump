package db

import (
	"context"
	"log"
)

const schema = `
CREATE TABLE IF NOT EXISTS system_configs (
    id              INT PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    landing_hosts   TEXT[] NOT NULL DEFAULT '{}',
    id_length       SMALLINT NOT NULL DEFAULT 6,
    id_min_length   SMALLINT NOT NULL DEFAULT 2,
    id_max_length   SMALLINT NOT NULL DEFAULT 10,
    not_found_mode  VARCHAR(20) NOT NULL DEFAULT 'content',
    not_found_value TEXT NOT NULL DEFAULT '你访问的页面不存在哦'
);

CREATE TABLE IF NOT EXISTS users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username    VARCHAR(50) NOT NULL UNIQUE,
    password    BYTEA NOT NULL,
    salt        BYTEA NOT NULL,
    role        SMALLINT NOT NULL DEFAULT 1,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS short_links (
    id          VARCHAR(20) PRIMARY KEY,
    url         TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    is_enabled  BOOLEAN NOT NULL DEFAULT true,
    created_by  VARCHAR(50) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_short_links_created_by ON short_links(created_by);
CREATE INDEX IF NOT EXISTS idx_short_links_created_at ON short_links(created_at DESC);
`

func RunMigrations() error {
	_, err := GetPostgresPool().Exec(context.Background(), schema)
	if err != nil {
		log.Printf("Migration failed: %v", err)
		return err
	}

	// Ensure default config row exists
	_, err = GetPostgresPool().Exec(context.Background(),
		`INSERT INTO system_configs (id) VALUES (1) ON CONFLICT DO NOTHING`)
	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}
