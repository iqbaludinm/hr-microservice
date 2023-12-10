package helper

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/iqbaludinm/hr-microservice/auth-service/config"
	_ "github.com/joho/godotenv/autoload"
)

type Claims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func GenerateJwt(issuer, name, email, phone, password string) (string, error) {
	session, _ := strconv.Atoi(config.SessionLogin)
	claims := &Claims{
		Name: name,
		Email: email,
		Phone: phone,
		Password: password,
		StandardClaims: jwt.StandardClaims{
			Issuer:    issuer,
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(session)).Unix(),
		},
	}
	tokens := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return tokens.SignedString([]byte(config.SecretKey))
}

func ParseJwt(cookie string) (string, string, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(cookie, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return "", "", err
	}

	return claims.Issuer, "", nil
}

// refresh
func GenerateRefreshJwt(issuer, name, email ,phone string) (string, error) {
	session, _ := strconv.Atoi(config.SessionRefreshToken)
	claims := &Claims{
		Name: name,
		Email: email,
		Phone: phone,
		StandardClaims: jwt.StandardClaims{
			Issuer:    issuer,
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(session)).Unix(),
		},
	}
	tokens := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return tokens.SignedString([]byte(config.SecretKeyRefresh))
}

func ParseRefreshJwt(cookie string) (string, error) {
	var claims jwt.StandardClaims
	token, err := jwt.ParseWithClaims(cookie, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKeyRefresh), nil
	})
	if err != nil || !token.Valid {
		return "", err
	}

	return claims.Issuer, err
}
