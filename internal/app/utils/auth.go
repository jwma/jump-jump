package utils

import (
	"crypto/rand"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/scrypt"
)

var SecretKey = os.Getenv("SECRET_KEY")

func RandomSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	return salt, err
}

func EncodePassword(password []byte, salt []byte) ([]byte, error) {
	return scrypt.Key(password, salt, 1<<15, 8, 1, 32)
}

func GenerateJWT(username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"identifier": username,
		"iat":        time.Now().Unix(),
		"exp":        time.Now().Add(time.Hour * 2).Unix(),
	})
	t, _ := token.SignedString([]byte(SecretKey))
	return t
}
