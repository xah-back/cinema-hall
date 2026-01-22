package services

import (
	"errors"
	"log/slog"
	"user-service/internal/dto"
	"user-service/internal/models"
	"user-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Create(req dto.CreateUserRequest) (*models.User, error)
	Get(id uint) (*models.User, error)
	List() ([]models.User, error)
	Update(id uint, req dto.UpdateUserRequest) (*models.User, error)
	Delete(id uint) error
}

type userService struct {
	repo repository.UserRepository
	log  *slog.Logger
}

func NewUserService(repo repository.UserRepository, log *slog.Logger) UserService {
	return &userService{repo: repo, log: log}
}

func (s *userService) Create(req dto.CreateUserRequest) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password), bcrypt.DefaultCost,
	)
	if err != nil {
		s.log.Error("failed to hash password", "err", err)
		return nil, err
	}

	user := &models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Role:     req.Role,
	}

	if user.Role == "" {
		user.Role = "user"
	}

	if err := s.repo.Create(user); err != nil {
		s.log.Error("failed to create user", "email", user.Email, "err", err)
		return nil, err
	}
	return user, nil
}

func (s *userService) Get(id uint) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Warn("user not found", "id", id)
			return nil, gorm.ErrRecordNotFound
		}
		s.log.Error("failed to get user", "id", id, "err", err)
		return nil, err
	}

	return user, nil
}

func (s *userService) List() ([]models.User, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		s.log.Error("failed to list users", "err", err)
		return nil, err
	}

	return users, nil
}

func (s *userService) Update(id uint, req dto.UpdateUserRequest) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Error("user not found for update", "id", id)
			return nil, gorm.ErrRecordNotFound
		}
		s.log.Error("failed to get user for update", "id", id, "err", err)
		return nil, err
	}

	if req.Email != nil {
		user.Email = *req.Email
	}

	if req.Role != nil {
		user.Role = *req.Role
	}

	if err := s.repo.Update(user); err != nil {
		s.log.Error("failed to update user", "id", id, "err", err)
		return nil, err
	}

	return user, nil
}

func (s *userService) Delete(id uint) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Warn("user not found for delete", "id", id)
			return err
		}
		s.log.Error("failed to get user for delete", "id", id, "err", err)
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		s.log.Error("failed to delete user", "id", id, "err", err)
		return err
	}

	return nil
}
