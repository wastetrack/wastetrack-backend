package helper

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/repository"
	"gorm.io/gorm"
)

type JWTHelper struct {
	SecretKey              string
	RefreshSecretKey       string
	AccessTokenTTL         time.Duration
	RefreshTokenTTL        time.Duration
	RefreshTokenRepository *repository.RefreshTokenRepository
}

func NewJWTHelper(
	secretKey, refreshSecretKey string,
	accessTTL, refreshTTL time.Duration,
	refreshTokenRepo *repository.RefreshTokenRepository,
) *JWTHelper {
	return &JWTHelper{
		SecretKey:              secretKey,
		RefreshSecretKey:       refreshSecretKey,
		AccessTokenTTL:         accessTTL,
		RefreshTokenTTL:        refreshTTL,
		RefreshTokenRepository: refreshTokenRepo,
	}
}

// Generate random token string
func (j *JWTHelper) generateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (j *JWTHelper) GenerateAccessToken(userID, role string, isEmailVerified bool) (string, error) {
	claims := &model.JWTClaims{
		UserID:          userID,
		Role:            role,
		IsEmailVerified: isEmailVerified,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}

// Enhanced refresh token generation with database storage
func (j *JWTHelper) GenerateRefreshToken(db *gorm.DB, userID uuid.UUID) (string, error) {
	// Generate random token string
	tokenStr, err := j.generateRandomToken()
	if err != nil {
		return "", err
	}

	// Create refresh token record
	refreshToken := &entity.RefreshToken{
		UserID:    userID,
		Token:     tokenStr,
		IsRevoked: false,
		ExpiresAt: time.Now().Add(j.RefreshTokenTTL),
	}

	// Store in database
	if err := j.RefreshTokenRepository.Create(db, refreshToken); err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (j *JWTHelper) ValidateAccessToken(tokenString string) (*model.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// Enhanced refresh token validation with database check
func (j *JWTHelper) ValidateRefreshToken(db *gorm.DB, tokenStr string) (*entity.RefreshToken, error) {
	refreshToken := &entity.RefreshToken{}

	// Check if token exists in database and is valid
	if err := j.RefreshTokenRepository.FindByToken(db, refreshToken, tokenStr); err != nil {
		return nil, err
	}

	return refreshToken, nil
}

// Revoke specific refresh token
func (j *JWTHelper) RevokeRefreshToken(db *gorm.DB, tokenStr string) error {
	return j.RefreshTokenRepository.RevokeToken(db, tokenStr)
}

// Revoke all user refresh tokens (logout from all devices)
func (j *JWTHelper) RevokeAllUserTokens(db *gorm.DB, userID uuid.UUID) error {
	return j.RefreshTokenRepository.RevokeAllUserTokens(db, userID)
}

// Cleanup expired tokens (run as background job)
func (j *JWTHelper) CleanupExpiredTokens(db *gorm.DB) error {
	return j.RefreshTokenRepository.DeleteExpiredTokens(db)
}

// Check active session limit
func (j *JWTHelper) CheckSessionLimit(db *gorm.DB, userID uuid.UUID, maxSessions int) (bool, error) {
	count, err := j.RefreshTokenRepository.CountActiveTokensByUser(db, userID)
	if err != nil {
		return false, err
	}
	return count < int64(maxSessions), nil
}

// NEW: Revoke oldest refresh token for a user
func (j *JWTHelper) RevokeOldestToken(db *gorm.DB, userID uuid.UUID) error {
	return j.RefreshTokenRepository.RevokeOldestTokenByUser(db, userID)
}

// NEW: Enforce session limit by automatically revoking oldest token
func (j *JWTHelper) EnforceSessionLimit(db *gorm.DB, userID uuid.UUID, maxSessions int) error {
	count, err := j.RefreshTokenRepository.CountActiveTokensByUser(db, userID)
	if err != nil {
		return err
	}

	// If we're at the limit, revoke exactly 1 oldest token to make room for the new one
	if int(count) >= maxSessions {
		return j.RevokeOldestToken(db, userID)
	}

	return nil
}
