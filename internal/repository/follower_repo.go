package repository

import (
	"gorm.io/gorm"
	"user_service/internal/model" // update this import path based on your project structure
)

type FollowerRelationRepository interface {
	Create(relation *model.FollowerRelation) error
	UpdateStatus(id uint, status string) error
	GetByID(id uint) (*model.FollowerRelation, error)
	ListFollowers(userID uint) ([]model.FollowerRelation, error)
	ListFollowing(followerID uint) ([]model.FollowerRelation, error)
	Delete(id uint) error
}

type followerRelationRepository struct {
	db *gorm.DB
}

func NewFollowerRelationRepository(db *gorm.DB) FollowerRelationRepository {
	return &followerRelationRepository{db: db}
}

func (r *followerRelationRepository) Create(relation *model.FollowerRelation) error {
	return r.db.Create(relation).Error
}

func (r *followerRelationRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&model.FollowerRelation{}).Where("id = ?", id).Update("status", status).Error
}

func (r *followerRelationRepository) GetByID(id uint) (*model.FollowerRelation, error) {
	var relation model.FollowerRelation
	err := r.db.First(&relation, id).Error
	if err != nil {
		return nil, err
	}
	return &relation, nil
}

func (r *followerRelationRepository) ListFollowers(userID uint) ([]model.FollowerRelation, error) {
	var relations []model.FollowerRelation
	err := r.db.Where("user_id = ?", userID).Find(&relations).Error
	return relations, err
}

func (r *followerRelationRepository) ListFollowing(followerID uint) ([]model.FollowerRelation, error) {
	var relations []model.FollowerRelation
	err := r.db.Where("follower_id = ?", followerID).Find(&relations).Error
	return relations, err
}

func (r *followerRelationRepository) Delete(id uint) error {
	return r.db.Delete(&model.FollowerRelation{}, id).Error
}
