package auth

import (
	"errors"
	"gophermart/internal/db"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) Create(user *db.User) error {
	return r.DB.Create(user).Error
}

func (r *Repository) FindByLogin(login string) (*db.User, error) {
	var user db.User
	if err := r.DB.Where("login = ?", login).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
