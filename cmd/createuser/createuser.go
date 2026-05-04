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
	username := flag.String("username", "", "username.")
	password := flag.String("password", "", "password.")
	role := flag.Int("role", models.RoleUser, "role, 1: normal user, 2: administrator.")
	flag.Parse()

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
		Role:        *role,
	}

	repo := repository.GetUserRepo(db.GetPostgresPool())
	err := repo.Save(user)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "create user %s successfully\n", *username)
}
