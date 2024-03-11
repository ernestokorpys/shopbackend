package util

import (
	"errors"
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

func VerifyJwt(tokenString string) (string, error) {
	// Parsear el token con las reclamaciones personalizadas que necesitas
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return "", err // Error al parsear el token
	}

	// Verificar si el token es válido
	if !token.Valid {
		return "", errors.New("Token inválido")
	}

	// Obtener las reclamaciones (claims) del token
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", errors.New("No se pudieron obtener las reclamaciones del token")
	}

	// Retornar el ID del emisor (issuer) desde las reclamaciones del token
	return claims.Issuer, nil
}
