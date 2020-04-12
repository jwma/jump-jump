package main

import (
	"flag"
	"fmt"
	"github.com/jwma/jump-jump/internal/app/db"
	"github.com/jwma/jump-jump/internal/app/models"
	"github.com/jwma/jump-jump/internal/app/repository"
	"os"
	"strings"
)

func main() {
	username := flag.String("username", "", "username.")
	password := flag.String("password", "", "password.")
	role := flag.Int("role", models.RoleUser, "role, 1: normal user, 2: administrator.")
	flag.Parse()

	user := &models.User{
		Username:    strings.TrimSpace(*username),
		RawPassword: strings.TrimSpace(*password),
		Role:        *role,
	}

	repo := repository.GetUserRepo(db.GetRedisClient())
	err := repo.Save(user)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	_, _ = fmt.Fprintf(os.Stdout, "create user %s successfully\n", *username)
}
