package orders

import (
	"errors"
	"gophermart/internal/db"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(order *db.Order) error {
	return r.db.Create(order).Error
}

func (r *Repository) FindByNumber(number string) (*db.Order, error) {
	var order db.Order
	if err := r.db.Where("number = ?", number).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (r *Repository) FindByUserID(userID uint) ([]db.Order, error) {
	var orders []db.Order
	if err := r.db.Where("user_id = ?", userID).Order("uploaded_at asc").Find(&orders).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return orders, nil
}
