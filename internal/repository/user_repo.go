package repository

import (
	"gorm.io/gorm"
	"user_service/internal/model"
)

// UserRepositoryImpl is the implementation of UserRepository
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepositoryImpl
func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

// CreateUser creates a new user in the database.
func (r *UserRepositoryImpl) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

// GetUserByID fetches a user by their ID.
func (r *UserRepositoryImpl) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername fetches a user by their username.
func (r *UserRepositoryImpl) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates the user record in the database.
func (r *UserRepositoryImpl) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

// DeleteUser deletes a user from the database.
func (r *UserRepositoryImpl) DeleteUser(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

// GetUsersPaginated retrieves users with pagination for infinite scrolling.
func (r *UserRepositoryImpl) GetUsersPaginated(page, pageSize int) ([]model.User, error) {
	var users []model.User
	offset := (page - 1) * pageSize

	if err := r.db.Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
