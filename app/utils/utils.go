package utils

import (
	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/scrypt"
	"math/rand"
	"time"
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
		"username": username,                             // Payload 部分可以加入用户名的记录，最终可通过解密 Token 得到
		"iat":      time.Now().Unix(),                    // 设置 Token 的签发时间
		"exp":      time.Now().Add(time.Hour * 2).Unix(), // 设置 Token 过期时间
	})

	// 使用一个密钥字符串对 Token 进行签名，只要密钥没有泄露，就没有人能篡改 Token 的数据
	jwt, _ := token.SignedString([]byte(beego.AppConfig.String("secret_key")))
	return jwt
}
