package service

import (
	"errors"
	"user_service/internal/model"
	"user_service/internal/transport/request"
	"user_service/internal/transport/response"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByID(id uint) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id uint) error
	GetUsersPaginated(page, pageSize int) ([]model.User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(req request.CreateUserRequest) (*response.UserResponseFull, error) {

	if req.Username == "" {
		return nil, errors.New("username cannot be empty")
	}
	if len(req.Username) < 3 {
		return nil, errors.New("username must be at least 3 characters long")
	}

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

func (s *UserService) UpdateUser(userID uint, req request.UpdateUserRequest, requestUserID uint) (*response.UserResponseFull, error) {

	if userID != requestUserID {
		return nil, errors.New("unauthorized: users can only update their own data")
	}

	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	if req.Username != "" {

		if len(req.Username) < 3 {
			return nil, errors.New("username must be at least 3 characters long")
		}

		if existingUser, err := s.repo.GetUserByID(userID); err == nil && existingUser.ID != user.ID {
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

func (s *UserService) DeleteUser(userID uint, requestUserID uint) error {

	if userID != requestUserID {
		return errors.New("unauthorized: users can only delete their own account")
	}

	return s.repo.DeleteUser(userID)
}

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

	return &response.PaginatedUsersResponse{
		Users: userResponses,
		Total: len(userResponses),
	}, nil
}
