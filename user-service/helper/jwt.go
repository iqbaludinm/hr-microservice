package helper

import (
	"github.com/golang-jwt/jwt"
	"github.com/iqbaludinm/hr-microservice/user-service/utils"
	_ "github.com/joho/godotenv/autoload"
)

type Claims struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	jwt.StandardClaims
}

var SecretKey = utils.GetEnv("SECRET_KEY")

// The return value of this function is as follows:
//
// 1. Issuer/Actor
//
// 2. Role ID
//
// 3. Permissions
//
// 4. Error
func ParseJwt(cookie string) (string, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(cookie, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil || !token.Valid {
		return "", err
	}

	return claims.Issuer, nil
}
