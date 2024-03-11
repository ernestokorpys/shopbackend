package util

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const SecretKey = "secret"

// Genera un token en base al ide del usuario para saber quien ingreso
// Este dura 12 hs issuer=emisor
func GenerateJwt(issuer string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    issuer,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //token expira luego de 20 hs
	})
	return claims.SignedString([]byte(SecretKey)) //se genera el token en base a la secret key
}

func Parsejwt(cookie string) (string, error) {
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil || !token.Valid {
		return "", err
	}
	claims := token.Claims.(*jwt.StandardClaims)
	return claims.Issuer, nil
}
