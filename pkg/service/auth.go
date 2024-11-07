package service

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/codeboris/todo"
	"github.com/codeboris/todo/pkg/repository"
	"github.com/dgrijalva/jwt-go"
)

const (
	salt       = "sdfsbpewur43q9tugq04gjv"
	signingKey = "sbapfewf00asfasdff3r34fsdfv4q"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserID int `json: "user_id"`
}

type AuthService struct {
	repos repository.Authorization
}

func NewAuthService(repos repository.Authorization) *AuthService {
	return &AuthService{repos: repos}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repos.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repos.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
	})

	return token.SignedString([]byte(signingKey))
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
