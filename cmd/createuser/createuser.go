package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/repository"
)

func main() {
	tenantID := flag.String("tenant-id", "", "tenant ID (UUID).")
	username := flag.String("username", "", "username.")
	password := flag.String("password", "", "password.")
	role := flag.String("role", models.RoleMember, "role: 'admin' or 'member'.")
	flag.Parse()

	if *tenantID == "" {
		fmt.Fprintf(os.Stderr, "tenant-id is required\n")
		os.Exit(1)
	}
	if *role != models.RoleAdmin && *role != models.RoleMember {
		fmt.Fprintf(os.Stderr, "role must be 'admin' or 'member'\n")
		os.Exit(1)
	}

	if err := db.InitPostgres(); err != nil {
		fmt.Fprintf(os.Stderr, "DB init failed: %v\n", err)
		os.Exit(1)
	}
	defer db.ClosePostgres()

	if err := db.RunMigrations(); err != nil {
		fmt.Fprintf(os.Stderr, "Migration failed: %v\n", err)
		os.Exit(1)
	}

	user := &models.User{
		Username:    strings.TrimSpace(*username),
		RawPassword: strings.TrimSpace(*password),
	}

	userRepo := repository.GetUserRepo(db.GetPostgresPool())
	err := userRepo.Save(user)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	member := &models.TenantMember{
		TenantID: strings.TrimSpace(*tenantID),
		UserID:   user.ID,
		Role:     strings.TrimSpace(*role),
	}

	memberRepo := repository.GetTenantMemberRepo(db.GetPostgresPool())
	err = memberRepo.Save(member)
	if err != nil {
		fmt.Fprintf(os.Stderr, "user created but membership failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "create user %s (id: %s) with role %s in tenant %s successfully\n",
		*username, user.ID, *role, *tenantID)
}
