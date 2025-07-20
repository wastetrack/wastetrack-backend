package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	Repository[entity.RefreshToken]
	Log *logrus.Logger
}

func NewRefreshTokenRepository(log *logrus.Logger) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		Log: log,
	}
}

func (r *RefreshTokenRepository) Create(db *gorm.DB, token *entity.RefreshToken) error {
	return db.Create(token).Error
}

func (r *RefreshTokenRepository) FindByToken(db *gorm.DB, token *entity.RefreshToken, tokenStr string) error {
	return db.Where("token = ? AND is_revoked = false AND expires_at > ?", tokenStr, time.Now()).First(token).Error
}

func (r *RefreshTokenRepository) RevokeToken(db *gorm.DB, tokenStr string) error {
	return db.Model(&entity.RefreshToken{}).
		Where("token = ?", tokenStr).
		Update("is_revoked", true).Error
}

func (r *RefreshTokenRepository) RevokeAllUserTokens(db *gorm.DB, userID uuid.UUID) error {
	return db.Model(&entity.RefreshToken{}).
		Where("user_id = ? AND is_revoked = false", userID).
		Update("is_revoked", true).Error
}

func (r *RefreshTokenRepository) RevokeOldestTokenByUser(db *gorm.DB, userID uuid.UUID) error {
	var oldestToken entity.RefreshToken
	if err := db.Where("user_id = ? AND is_revoked = false AND expires_at > ?", userID, time.Now()).
		Order("created_at ASC").
		First(&oldestToken).Error; err != nil {
		return err
	}

	return db.Model(&oldestToken).
		Update("is_revoked", true).Error
}

func (r *RefreshTokenRepository) DeleteExpiredTokens(db *gorm.DB) error {
	return db.Where("expires_at < ? OR is_revoked = true", time.Now()).
		Delete(&entity.RefreshToken{}).Error
}

func (r *RefreshTokenRepository) CountActiveTokensByUser(db *gorm.DB, userID uuid.UUID) (int64, error) {
	var count int64
	err := db.Model(&entity.RefreshToken{}).
		Where("user_id = ? AND is_revoked = false AND expires_at > ?", userID, time.Now()).
		Count(&count).Error
	return count, err
}
