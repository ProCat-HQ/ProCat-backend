package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"github.com/procat-hq/procat-backend/internal/app/repository"
	"os"
	"time"
)

type UserService struct {
	repo repository.User
}

const (
	tokenTTL = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	model.TokenClaimsExtension
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user model.SignUpInput) (int, error) {
	user.Password = generatePasswordHash(user.PhoneNumber, user.Password)

	u := &model.User{
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		Password:    user.Password,
		IsConfirmed: false,
		Role:        "user",
	}
	return s.repo.CreateUser(*u)
}

func (s *UserService) GenerateToken(phoneNumber, password string) (string, error) {
	user, err := s.repo.GetUser(phoneNumber, generatePasswordHash(phoneNumber, password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		TokenClaimsExtension: model.TokenClaimsExtension{
			UserId:   user.Id,
			UserRole: user.Role,
		},
	})

	signingKey := os.Getenv("SIGNING_KEY")
	return token.SignedString([]byte(signingKey))
}

func (s *UserService) ParseToken(accessToken string) (*model.TokenClaimsExtension, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(os.Getenv("SIGNING_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*tokenClaims); !ok {
		return nil, errors.New("invalid token claims type")
	} else {
		return &claims.TokenClaimsExtension, nil
	}

}

func generatePasswordHash(phoneNumber, password string) string {
	hash := sha1.New()
	hash.Write([]byte(phoneNumber + password))

	salt := os.Getenv("PASSWORD_SALT")
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
