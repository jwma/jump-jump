package main

import (
	"flag"
	"fmt"
	"github.com/jwma/jump-jump/internal/app/models"
	"os"
)

func main() {
	username := flag.String("username", "", "username.")
	password := flag.String("password", "", "password.")
	role := flag.Int("role", models.RoleUser, "role, 1: normal user, 2: administrator.")
	flag.Parse()

	user := models.User{
		Username:    *username,
		RawPassword: *password,
		Role:        *role,
	}

	err := user.Save()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	_, _ = fmt.Fprintf(os.Stdout, "create user %s successfully\n", *username)
}
