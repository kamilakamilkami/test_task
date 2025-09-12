package utils

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func GenerateUUID() uuid.UUID {
	return uuid.New()
}

// Проверка валидности email
func EmailInvalid(email string) bool {
	return !strings.Contains(email, "@") || !strings.Contains(email, ".") || len(email) < 6
}

// Хеширование пароля
func HashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hashed)
}
