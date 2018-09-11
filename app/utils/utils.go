package utils

import (
	"math/rand"
	"time"
	"golang.org/x/crypto/scrypt"
	"github.com/dgrijalva/jwt-go"
	"github.com/astaxie/beego"
)

var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyz")
var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[seededRand.Intn(len(letterRunes))]
	}
	return string(b)
}

// 生成随机盐值
func RandomSalt(size int) ([]byte, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	return salt, err
}

// 加盐密码哈希
func EncodePassword(password []byte, salt []byte) ([]byte, error) {
	return scrypt.Key(password, salt, 1<<15, 8, 1, 32)
}

func GenerateJWT(username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour * 2).Unix(),
	})

	jwt, _ := token.SignedString([]byte(beego.AppConfig.String("secret_key")))
	return jwt
}
