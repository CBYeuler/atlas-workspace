package services

import (
	"errors"
	"time"

	"github.com/CBYeuler/atlas-workspace/backend/internal/config"
	"github.com/CBYeuler/atlas-workspace/backend/internal/database"
	"github.com/CBYeuler/atlas-workspace/backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService() *AuthService {
	return &AuthService{
		db: database.DB,
	}
}

// ========= Input / Output DTO =========

type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// ========= Public Methods =========

func (s *AuthService) Register(input RegisterInput) (*models.User, *TokenPair, error) {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		return nil, nil, errors.New("user already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, err
	}
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}
	user := &models.User{
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		FullName:     input.FullName,
	}
	if err := s.db.Create(user).Error; err != nil {
		return nil, nil, err
	}
	// Generate tokens
	tokens, err := s.generateAndStoreTokens(user.ID)
	if err != nil {
		return nil, nil, err
	}
	return user, tokens, nil

}

func (s *AuthService) Login(input LoginInput) (*models.User, *TokenPair, error) {
	var user models.User
	if err := s.db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("invalid credentials")
		}
		return nil, nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	tokens, err := s.generateAndStoreTokens(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return &user, tokens, nil
}

func (s *AuthService) Refresh(refreshToken string) (*TokenPair, error) {
	// 1) JWT verify
	claims, err := parseRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	userID := claims.Subject
	// 2) Check in DB
	var session models.Session
	if err := s.db.Where("refresh_token = ?", refreshToken).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("session not found")
		}
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, errors.New("session expired")
	}

	if session.UserID != userID {
		return nil, errors.New("invalid session owner")
	}
	// 3) Generate new tokens
	tokens, err := s.generateAndStoreTokens(userID)
	if err != nil {
		return nil, err
	}
	// optionally, delete old session (token rotation)
	_ = s.db.Delete(&session).Error

	return tokens, nil
}

// ========= Private Helpers =========

func (s *AuthService) generateAndStoreTokens(userID string) (*TokenPair, error) {
	accessToken, err := generateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExp, err := generateRefreshToken(userID)
	if err != nil {
		return nil, err
	}
	// Store session in DB
	session := &models.Session{
		UserID:       userID,
		RefreshToken: refreshToken,
		ExpiresAt:    refreshExp,
	}

	if err := s.db.Create(session).Error; err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ========= JWT Helpers =========

type CustomClaims struct {
	jwt.RegisteredClaims
}

func generateAccessToken(userID string) (string, error) {
	cfg := config.C

	dur, err := time.ParseDuration(cfg.JWTAccessExpires)
	if err != nil {
		dur = 15 * time.Minute
	}

	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(dur)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(cfg.JWTAccessSecret))
}

func generateRefreshToken(userID string) (string, time.Time, error) {
	cfg := config.C

	dur, err := time.ParseDuration(cfg.JWTRefreshExpires)
	if err != nil {
		dur = 30 * 24 * time.Hour // fallback 30 days
	}

	exp := time.Now().Add(dur)

	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(cfg.JWTRefreshSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, exp, nil
}
func parseRefreshToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	cfg := config.C

	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTRefreshSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return &claims.RegisteredClaims, nil
}
