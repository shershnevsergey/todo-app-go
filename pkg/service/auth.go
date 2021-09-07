package service

import (
	"crypto/sha1"
	"fmt"
	"github.com/Shv-sergey70/todo-app-go"
	"github.com/Shv-sergey70/todo-app-go/pkg/repository"
)

const salt = "sdadwfawfaofmanfgoiqwfona"

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)

	return s.repo.CreateUser(user)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

