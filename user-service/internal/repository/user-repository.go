package repository

import (
	"log/slog"
	"user-service/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetAll() ([]models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
}

type userRepository struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewUserRepository(db *gorm.DB, log *slog.Logger) UserRepository {
	return &userRepository{db: db, log: log}
}

func (r *userRepository) Create(u *models.User) error {
	if err := r.db.Create(u).Error; err != nil {
		r.log.Error("failed to create user", "err", err)
		return err
	}
	return nil
}

func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var u models.User
	if err := r.db.First(&u, id).Error; err != nil {
		r.log.Error("failed to get user by id", "id", id, "err", err)
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		r.log.Error("failed to get user by email", "email", email, "err", err)
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetAll() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		r.log.Error("failed to get all users", "err", err)
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Update(u *models.User) error {
	if err := r.db.Save(u).Error; err != nil {
		r.log.Error("failed to update user", "id", u.ID, "err", err)
		return err
	}
	return nil
}

func (r *userRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.User{}, id).Error; err != nil {
		r.log.Error("failed to delete user", "id", id, "err", err)
		return err
	}
	return nil
}
