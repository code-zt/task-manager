package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

// CheckPassword проверяет, соответствует ли введённый пароль хешу
func ComparePassword(hashedPassword, password string) error {
	// Сравниваем хеш пароля с введённым паролем
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
