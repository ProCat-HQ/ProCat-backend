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
	"github.com/procat-hq/procat-backend/internal/kzgov"
	"os"
	"strconv"
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

func generatePasswordHash(phoneNumber, password string) string {
	hash := sha1.New()
	hash.Write([]byte(phoneNumber + password))

	salt := os.Getenv("PASSWORD_SALT")
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *UserService) CreateUser(user model.SignUpInput) (int, error) {
	return s.repo.CreateUser(model.SignUpInput{
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		Password:    generatePasswordHash(user.PhoneNumber, user.Password),
	})
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
			// here i need to delete all sessions
			err := s.repo.WipeRefreshSessions(userId)
			if err != nil {
				return "", "", errors.New("suspicious activity detected, but: " + err.Error())
			}
			return "", "", errors.New("suspicious activity detected")
		}
		return "", "", err
	}
	if refreshSession.Fingerprint != fingerprint {
		// but here only one with certain refresh token
		_, err := s.repo.DeleteUserRefreshSession(refreshToken, userId)
		if err != nil {
			return "", "", errors.New("invalid refresh session: suspicious activity detected, but: " + err.Error())
		}
		return "", "", errors.New("invalid refresh session: suspicious activity detected")
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

func (s *UserService) GetAllUsers(limit, page, role, isConfirmed string) (int, []model.User, error) {
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return 0, nil, err
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return 0, nil, err
	}

	offset := limitInt * pageInt

	if isConfirmed != "" {
		_, err = strconv.ParseBool(isConfirmed)
		if err != nil {
			return 0, nil, err
		}
	}
	count, users, err := s.repo.GetAllUsers(limitInt, offset, role, isConfirmed)
	if err != nil {
		return 0, nil, err
	}
	return count, users, nil
}

func (s *UserService) GetUserById(userId int) (model.User, error) {
	user, err := s.repo.GetUserById(userId)
	return user, err
}

func (s *UserService) DeleteUserById(userId int) error {
	return s.repo.DeleteUserById(userId)
}

func (s *UserService) CheckPassword(password string, userId int) (bool, error) {
	user, err := s.repo.GetUserWithPasswordById(userId)
	if err != nil {
		return false, err
	}

	passwordHash := generatePasswordHash(user.PhoneNumber, password)
	return passwordHash == user.PasswordHash, nil
}

func (s *UserService) ChangeFullName(userId int, fullName string) error {
	return s.repo.ChangeFullName(userId, fullName)
}

func (s *UserService) ChangeIdentificationNumber(userId int, identificationNumber string) error {
	kzGovResponse, err := kzgov.GetArrear(identificationNumber)
	if err != nil {
		return err
	}

	user, err := s.repo.GetUserById(userId)
	if err != nil {
		return err
	}

	if !kzgov.CompareNames(kzGovResponse.NameKk, kzGovResponse.NameRu, user.FullName) {
		return errors.New("fullname from kz.gov service doesn't match with user's fullname")
	}

	return s.repo.ChangeIdentificationNumber(userId, identificationNumber)
}

func (s *UserService) ChangePassword(userId int, password string) error {
	user, err := s.repo.GetUserById(userId)
	if err != nil {
		return err
	}

	passwordHash := generatePasswordHash(user.PhoneNumber, password)
	return s.repo.ChangePassword(userId, passwordHash)
}

func (s *UserService) ChangePhoneNumber(userId int, phoneNumber, password string) error {
	passwordHash := generatePasswordHash(phoneNumber, password)
	return s.repo.ChangePhoneNumber(userId, phoneNumber, passwordHash)
}

func (s *UserService) ChangeEmail(userId int, email string) error {
	return s.repo.ChangeEmail(userId, email)
}

func getRolePriority(role string) (int, error) {
	switch role {
	case model.UserRole:
		return 1, nil
	case model.DeliverymanRole:
		return 2, nil
	case model.ModeratorRole:
		return 3, nil
	case model.AdminRole:
		return 4, nil
	default:
		return 0, errors.New("unknown role")
	}
}

func (s *UserService) ChangeUserRole(userId int, role string) error {
	if _, err := getRolePriority(role); err != nil {
		return err
	}

	return s.repo.ChangeUserRole(userId, role)
}
