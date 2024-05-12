package service

import (
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"github.com/procat-hq/procat-backend/internal/app/repository"
	"os"
	"time"
)

type UserService struct {
	repo repository.User
}

const (
	maxSessionsNumber = 5
	accessTTL         = 30 * time.Minute
	refreshTTL        = 60 * 24 * time.Hour
)

type accessTokenClaims struct {
	jwt.StandardClaims
	model.AccessTokenClaimsExtension
}

type refreshTokenClaims struct {
	jwt.StandardClaims
	model.RefreshTokenClaimsExtension
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

func (s *UserService) GetUserByCredentials(phoneNumber, password string) (model.User, error) {
	user, err := s.repo.GetUser(phoneNumber, generatePasswordHash(phoneNumber, password))
	return user, err
}

func generateTokensPair(user model.User) (string, string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &accessTokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		AccessTokenClaimsExtension: model.AccessTokenClaimsExtension{
			UserId:   user.Id,
			UserRole: user.Role,
		},
	})

	uuidVar := uuid.New().String()
	if uuidVar == "" {
		return "", "", errors.New("can't generate uuid for jti")
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &refreshTokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(refreshTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        uuidVar,
		},
		RefreshTokenClaimsExtension: model.RefreshTokenClaimsExtension{
			UserId: user.Id,
		},
	})

	accessSigningKey := os.Getenv("ACCESS_SIGNING_KEY")
	signedAccessToken, err := accessToken.SignedString([]byte(accessSigningKey))
	if err != nil {
		return "", "", err
	}

	refreshSigningKey := os.Getenv("REFRESH_SIGNING_KEY")
	signedRefreshToken, err := refreshToken.SignedString([]byte(refreshSigningKey))
	if err != nil {
		return "", "", err
	}

	return signedAccessToken, signedRefreshToken, nil
}

func (s *UserService) GenerateTokens(user model.User, fingerprint string) (string, string, error) {
	signedAccessToken, signedRefreshToken, err := generateTokensPair(user)
	if err != nil {
		return "", "", err
	}

	err = s.repo.SaveSessionData(signedRefreshToken, fingerprint, user.Id)
	if err != nil {
		return "", "", err
	}

	refreshSessions, err := s.repo.GetRefreshSessions(user.Id)
	if len(refreshSessions) > maxSessionsNumber {
		err = s.repo.WipeRefreshSessionsWithFingerprint(fingerprint, user.Id)
		if err != nil {
			return "", "", err
		}
	}

	return signedAccessToken, signedRefreshToken, nil
}

func (s *UserService) ParseAccessToken(accessToken string) (*model.AccessTokenClaimsExtension, error) {
	token, err := jwt.ParseWithClaims(accessToken, &accessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(os.Getenv("ACCESS_SIGNING_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*accessTokenClaims); !ok {
		return nil, errors.New("invalid access token claims type")
	} else {
		return &claims.AccessTokenClaimsExtension, nil
	}
}

func (s *UserService) ParseRefreshToken(refreshToken string) (*model.RefreshTokenClaimsExtension, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &refreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(os.Getenv("REFRESH_SIGNING_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*refreshTokenClaims); !ok {
		return nil, errors.New("invalid refresh token claims type")
	} else {
		return &claims.RefreshTokenClaimsExtension, nil
	}
}

func (s *UserService) LogoutUser(refreshToken string, userId int) (int, error) {
	status, err := s.repo.DeleteUserRefreshSession(refreshToken, userId)
	return status, err
}

func (s *UserService) RegenerateTokens(userId int, refreshToken, fingerprint string) (string, string, error) {
	refreshSession, err := s.repo.GetRefreshSession(refreshToken, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := s.repo.WipeRefreshSessions(userId)
			if err != nil {
				return "", "", errors.New("suspicious activity detected, but: " + err.Error())
			}
			return "", "", errors.New("suspicious activity detected")
		}
		return "", "", err
	}
	if refreshSession.Fingerprint != fingerprint {
		return "", "", errors.New("invalid refresh session")
	}

	user, err := s.repo.GetUserById(userId)
	if err != nil {
		return "", "", err
	}

	signedAccessToken, signedRefreshToken, err := generateTokensPair(user)
	if err != nil {
		return "", "", err
	}

	err = s.repo.SaveSessionData(signedRefreshToken, fingerprint, user.Id)
	if err != nil {
		return "", "", err
	}

	return signedAccessToken, signedRefreshToken, nil
}

func generatePasswordHash(phoneNumber, password string) string {
	hash := sha1.New()
	hash.Write([]byte(phoneNumber + password))

	salt := os.Getenv("PASSWORD_SALT")
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
