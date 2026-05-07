package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// This script migrates the database from the old schema (users with tenant_id, role)
// to the new schema (globally unique users + tenant_members + tenant_invitations).
//
// Usage:
//
//	DATABASE_URL=postgres://user:pass@host:5432/dbname go run ./cmd/migraterefactor
//
// What it does:
//  1. Creates tenant_members and tenant_invitations tables (if not exist).
//  2. Adds is_active column to users (if not exist).
//  3. Reads all existing user rows (id, tenant_id, username, password, salt, role).
//  4. Groups by username:
//     - Unique username: keeps the row, creates a tenant_members entry.
//     - Duplicate username across tenants: keeps the first row, creates tenant_members
//     entries for all tenants, updates foreign keys (user_preferences), and deletes
//     the duplicate rows. Logs the merge for manual review.
//  5. Drops the old UNIQUE(tenant_id, username) constraint, tenant_id and role columns.
//  6. Adds UNIQUE(username) constraint.

type oldUser struct {
	id       string
	tenantID string
	username string
	password []byte
	salt     []byte
	role     int
}

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	ctx := context.Background()

	// Step 1: Create new tables
	log.Println("Step 1: Creating tenant_members and tenant_invitations tables...")
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS tenant_members (
		    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
		    user_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		    role      VARCHAR(20) NOT NULL DEFAULT 'member',
		    joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
		    PRIMARY KEY (tenant_id, user_id)
		)`)
	if err != nil {
		log.Fatalf("Failed to create tenant_members: %v", err)
	}

	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS tenant_invitations (
		    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		    tenant_id  UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
		    inviter_id UUID NOT NULL REFERENCES users(id),
		    invitee_id UUID NOT NULL REFERENCES users(id),
		    status     VARCHAR(20) NOT NULL DEFAULT 'pending',
		    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
		    UNIQUE(tenant_id, invitee_id, status)
		)`)
	if err != nil {
		log.Fatalf("Failed to create tenant_invitations: %v", err)
	}

	// Step 2: Add is_active column if not exists
	log.Println("Step 2: Adding is_active column to users...")
	_, err = pool.Exec(ctx, `
		ALTER TABLE users ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT true`)
	if err != nil {
		log.Fatalf("Failed to add is_active column: %v", err)
	}

	// Step 3: Read all existing users
	log.Println("Step 3: Reading existing users...")
	rows, err := pool.Query(ctx, `
		SELECT id, tenant_id, username, password, salt, role
		FROM users ORDER BY created_at`)
	if err != nil {
		log.Fatalf("Failed to read users: %v", err)
	}

	usersByTenant := make(map[string][]oldUser) // tenant_id -> list of users
	allUsers := make([]oldUser, 0)
	for rows.Next() {
		u := oldUser{}
		rows.Scan(&u.id, &u.tenantID, &u.username, &u.password, &u.salt, &u.role)
		usersByTenant[u.tenantID] = append(usersByTenant[u.tenantID], u)
		allUsers = append(allUsers, u)
	}
	rows.Close()

	log.Printf("Found %d user records across %d tenants\n", len(allUsers), len(usersByTenant))

	// Group by username to detect duplicates
	usersByUsername := make(map[string][]oldUser)
	for _, u := range allUsers {
		usersByUsername[u.username] = append(usersByUsername[u.username], u)
	}

	// Step 4: Migrate data
	log.Println("Step 4: Migrating data...")
	mergedCount := 0
	for username, users := range usersByUsername {
		// Pick the first user as the "canonical" one
		canonical := users[0]

		for _, u := range users {
			role := "member"
			if u.role == 2 {
				role = "admin"
			}

			// Create tenant_members entry using canonical user's ID
			_, err := pool.Exec(ctx, `
				INSERT INTO tenant_members (tenant_id, user_id, role, joined_at)
				VALUES ($1, $2, $3, now())
				ON CONFLICT (tenant_id, user_id) DO NOTHING`,
				u.tenantID, canonical.id, role)
			if err != nil {
				log.Fatalf("Failed to create membership for user %s in tenant %s: %v", u.username, u.tenantID, err)
			}

			// If this is a duplicate (different id from canonical), migrate FKs and delete
			if u.id != canonical.id {
				mergedCount++
				log.Printf("MERGE: username '%s' — keeping id %s, merging id %s (tenant %s)",
					username, canonical.id, u.id, u.tenantID)

				// Migrate user_preferences
				_, err := pool.Exec(ctx, `
					UPDATE user_preferences SET user_id = $1 WHERE user_id = $2`,
					canonical.id, u.id)
				if err != nil {
					log.Printf("Warning: failed to migrate preferences for user %s: %v", u.id, err)
				}

				// Update short_links created_by to canonical username (already the same)
				// No action needed since created_by stores username, not user_id

				// Delete the duplicate user row
				_, err = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, u.id)
				if err != nil {
					log.Fatalf("Failed to delete duplicate user %s: %v", u.id, err)
				}
			}
		}
	}

	log.Printf("Merged %d duplicate user(s)\n", mergedCount)
	if mergedCount > 0 {
		log.Println("WARNING: Duplicate usernames were found and merged. Please review the logs above for conflicts.")
		fmt.Println("\n⚠ ATTENTION: Some usernames existed in multiple tenants and were auto-merged.")
		fmt.Println("The first occurrence's credentials were kept. Please verify with affected users.")
	}

	// Step 5: Drop old columns and constraints, add new ones
	log.Println("Step 5: Altering users table...")

	// Drop the old unique constraint on (tenant_id, username)
	_, err = pool.Exec(ctx, `
		ALTER TABLE users DROP CONSTRAINT IF EXISTS users_tenant_id_username_key`)
	if err != nil {
		log.Printf("Warning: could not drop old unique constraint: %v", err)
	}

	// Also try the auto-generated name pattern
	_, _ = pool.Exec(ctx, `
		ALTER TABLE users DROP CONSTRAINT IF EXISTS users_tenant_id_username_uniq`)

	// Drop tenant_id column
	_, err = pool.Exec(ctx, `ALTER TABLE users DROP COLUMN IF EXISTS tenant_id`)
	if err != nil {
		log.Fatalf("Failed to drop tenant_id column: %v", err)
	}

	// Drop role column
	_, err = pool.Exec(ctx, `ALTER TABLE users DROP COLUMN IF EXISTS role`)
	if err != nil {
		log.Fatalf("Failed to drop role column: %v", err)
	}

	// Add UNIQUE constraint on username
	_, err = pool.Exec(ctx, `
		ALTER TABLE users ADD CONSTRAINT users_username_key UNIQUE (username)`)
	if err != nil {
		// May already exist
		log.Printf("Note: username unique constraint: %v (may already exist)", err)
	}

	log.Println("Migration completed successfully!")
}
