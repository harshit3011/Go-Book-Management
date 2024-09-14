package utils

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))


func GenerateToken(userID string) (string, error) {
    claims := jwt.MapClaims{}
    claims["user_id"] = userID
    claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
    claims := jwt.MapClaims{}

    // Parse the token
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    // Check if token is valid
    if err != nil {
        return nil, err
    }

    // Check the token signing method
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, errors.New("invalid token signing method")
    }

    // Validate the token expiration time
    if !token.Valid {
        return nil, errors.New("invalid or expired token")
    }

    return claims, nil
}
