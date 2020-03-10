package main

import (
	"flag"
	"fmt"
	"github.com/jwma/jump-jump/internal/app/models"
	"os"
)

func main() {
	username := flag.String("username", "", "A admin user username.")
	password := flag.String("password", "", "A admin user password.")
	flag.Parse()

	user := models.User{
		Username:    *username,
		RawPassword: *password,
	}

	err := user.Save()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	_, _ = fmt.Fprintf(os.Stdout, "create user %s successfully\n", *username)
}
