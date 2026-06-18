package user

import (
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	logger    *slog.Logger
	repo      Repository
	jwtSecret string
}

type Service interface {
	Register(u User) (User, error)
	Login(email, password string) (accessToken string, err error)
	GetMe(userID int) (User, error)
}

func NewService(logger *slog.Logger, repo Repository, jwtSecret string) Service {
	return &service{
		logger:    logger,
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (s *service) Register(u User) (User, error) {
	if s.repo.ExistsByEmail(u.Email) {
		return User{}, errors.New("email already registered")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("bcrypt error", slog.Any("err", err))
		return User{}, errors.New("failed to process password")
	}
	u.Password = string(hashed)

	created, err := s.repo.Create(u)
	if err != nil {
		s.logger.Error("create user error", slog.Any("err", err))
		return User{}, errors.New("failed to create user")
	}

	return created, nil
}

func (s *service) Login(email, password string) (string, error) {
	u, err := s.repo.GetByEmail(email)
	if err != nil || u.ID == 0 {
		return "", errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", errors.New("wrong email or password")
	}

	token, err := s.generateToken(u.ID)
	if err != nil {
		s.logger.Error("generate token error", slog.Any("err", err))
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func (s *service) GetMe(userID int) (User, error) {
	u, err := s.repo.GetByID(userID)
	if err != nil {
		return User{}, errors.New("user not found")
	}
	return u, nil
}

func (s *service) generateToken(userID int) (string, error) {
	type jwtClaims struct {
		ID int `json:"id"`
		jwt.RegisteredClaims
	}

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims{
		ID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		},
	})

	return token.SignedString([]byte(s.jwtSecret))
}
