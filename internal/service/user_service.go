package service

import (
	"errors"
	"user_service/internal/model"
	"user_service/internal/transport/request"
	"user_service/internal/transport/response"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByID(id uint) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id uint) error
	GetUsersPaginated(page, pageSize int) ([]model.User, error)
}

// UserService provides user-related operations
type UserService struct {
	repo UserRepository
}

// NewUserService creates a new UserService instance
func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser creates a new user after validation
func (s *UserService) CreateUser(req request.CreateUserRequest) (*response.UserResponseFull, error) {
	// Validate username
	if req.Username == "" {
		return nil, errors.New("username cannot be empty")
	}
	if len(req.Username) < 3 {
		return nil, errors.New("username must be at least 3 characters long")
	}

	// Check if username already exists
	if _, err := s.repo.GetUserByUsername(req.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	user := &model.User{
		Username: req.Username,
		Bio:      req.Bio,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return &response.UserResponseFull{
		ID:        user.ID,
		Username:  user.Username,
		AvatarURL: user.AvatarURL,
		Bio:       user.Bio,
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id uint) (*response.UserResponseFull, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return &response.UserResponseFull{
		ID:             user.ID,
		Username:       user.Username,
		AvatarURL:      user.AvatarURL,
		Bio:            user.Bio,
		FollowersCount: user.FollowersCount,
		FollowingCount: user.FollowingCount,
	}, nil
}

// UpdateUser updates user data with validation
func (s *UserService) UpdateUser(userID uint, req request.UpdateUserRequest, requestUserID uint) (*response.UserResponseFull, error) {
	// Check if the requesting user is the same as the user being updated
	if userID != requestUserID {
		return nil, errors.New("unauthorized: users can only update their own data")
	}

	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Update fields if they are provided in the request
	if req.Username != "" {
		// Validate username
		if len(req.Username) < 3 {
			return nil, errors.New("username must be at least 3 characters long")
		}

		// Check if new username is already taken
		if existingUser, err := s.repo.GetUserByUsername(req.Username); err == nil && existingUser.ID != user.ID {
			return nil, errors.New("username already taken")
		}
		user.Username = req.Username
	}

	if req.Bio != "" {
		user.Bio = req.Bio
	}

	if err := s.repo.UpdateUser(user); err != nil {
		return nil, err
	}

	return &response.UserResponseFull{
		ID:             user.ID,
		Username:       user.Username,
		AvatarURL:      user.AvatarURL,
		Bio:            user.Bio,
		FollowersCount: user.FollowersCount,
		FollowingCount: user.FollowingCount,
	}, nil
}

// DeleteUser deletes a user with validation
func (s *UserService) DeleteUser(userID uint, requestUserID uint) error {
	// Check if the requesting user is the same as the user being deleted
	if userID != requestUserID {
		return errors.New("unauthorized: users can only delete their own account")
	}

	return s.repo.DeleteUser(userID)
}

// GetUsersPaginated retrieves a paginated list of users
func (s *UserService) GetUsersPaginated(page, pageSize int) (*response.PaginatedUsersResponse, error) {
	users, err := s.repo.GetUsersPaginated(page, pageSize)
	if err != nil {
		return nil, err
	}

	var userResponses []response.UserResponseShort
	for _, user := range users {
		userResponses = append(userResponses, response.UserResponseShort{
			ID:        user.ID,
			Username:  user.Username,
			AvatarURL: user.AvatarURL,
		})
	}

	// Note: In a real implementation, you'd want to get the total count from the repository
	// For simplicity, we're just returning the count of the current page here
	return &response.PaginatedUsersResponse{
		Users: userResponses,
		Total: len(userResponses),
	}, nil
}
