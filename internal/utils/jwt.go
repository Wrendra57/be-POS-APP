package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
	"time"
)

type Claimed struct {
	User_id uuid.UUID `json:"user_id"`
	Level   string    `json:"level"`
	jwt.RegisteredClaims
}

func GenerateJWT(uuid uuid.UUID, level string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": uuid,
		"level":   level,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	secret := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseJWT(tokenString string) (*Claimed, error) {
	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, &Claimed{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println(err.Error())

		return &Claimed{}, err
	}
	if !token.Valid {
		fmt.Println("invalid token")
		return &Claimed{}, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*Claimed)
	if !ok {
		fmt.Println("invalisd token")

		return &Claimed{}, fmt.Errorf("failed to parse claims")
	}

	return claims, nil

}
