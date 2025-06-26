package repository

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository[T any] struct {
	DB *gorm.DB
}

func (r *Repository[T]) Create(db *gorm.DB, entity *T) error {
	return db.Create(entity).Error
}
func (r *Repository[T]) CreateBatch(db *gorm.DB, entities []*T) error {
	// Handle empty slice case
	if len(entities) == 0 {
		return errors.New("empty slice found")
	}

	return db.Create(&entities).Error
}

func (r *Repository[T]) Update(db *gorm.DB, entity *T) error {
	return db.Save(entity).Error
}

func (r *Repository[T]) Delete(db *gorm.DB, entity *T) error {
	return db.Delete(entity).Error
}

func (r *Repository[T]) CountById(db *gorm.DB, id any) (int64, error) {
	var total int64
	err := db.Model(new(T)).Where("id = ?", id).Count(&total).Error
	return total, err
}

func (r *Repository[T]) FindById(db *gorm.DB, entity *T, id any) error {
	return db.Where("id = ?", id).Take(entity).Error
}

func (r *Repository[T]) Upsert(db *gorm.DB, entity *T, conflictColumns []clause.Column, updateColumns []string) error {
	return db.Clauses(clause.OnConflict{
		Columns:   conflictColumns,
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}).Create(entity).Error
}
