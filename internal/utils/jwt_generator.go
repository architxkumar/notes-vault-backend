package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"log"
	"notes-vault-backend/internal/model"
	"os"
	"time"
)

func JwtGenerator(userEntry model.User) (string, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
		return "", err
	}
	key := []byte(os.Getenv("JWT_SECRET"))
	var claims = jwt.MapClaims{"id": userEntry.Id, "email": userEntry.Email, "exp": time.Now().Add(time.Hour).Unix()}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
