package handlers

import (
	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/mail"
	"os"
	"strings"
	"time"
)

type User struct {
	Id             string
	Email          string
	HashedPassword string
}

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignupHandler(c *fiber.Ctx) error {
	db, err := gorm.Open(sqlite.Open("./db/notesapp.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	contentType := c.Get("Content-Type")
	if contentType != "application/json" {
		return fiber.ErrUnsupportedMediaType
	}
	u := new(SignUpRequest)
	if err := c.BodyParser(u); err != nil {
		return fiber.ErrBadRequest
	}
	u.Email = strings.Trim(u.Email, " ")
	u.Password = strings.Trim(u.Password, " ")
	if u.Email == "" || u.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email or Password is empty")
	}
	if emailValidator(u.Email) == false {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Email")
	}

	if len(u.Password) < 8 {
		return fiber.NewError(fiber.StatusBadRequest, "Password must be at least 8 characters")
	}
	result := db.First(&u, "email = ?", u.Email)
	if result.RowsAffected != 0 {
		return fiber.NewError(fiber.StatusConflict, "Email already exists")
	}
	entryUuid := uuid.Must(uuid.NewRandom()).String()
	entryHashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err)
		return fiber.ErrInternalServerError
	}
	databaseEntry := User{Id: entryUuid, Email: u.Email, HashedPassword: string(entryHashedPassword)}
	result = db.Select("id", "email", "hashed_password").Create(&databaseEntry)
	if result.Error != nil {
		log.Error(err)
		return fiber.ErrInternalServerError
	}
	token, err := jwtGenerator(databaseEntry)
	if err != nil {
		log.Error(err)
		return fiber.ErrInternalServerError
	}
	return c.Status(201).JSON(fiber.Map{
		"token": token,
	})
}

func emailValidator(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func jwtGenerator(userEntry User) (string, error) {
	err := godotenv.Load(".env")
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
