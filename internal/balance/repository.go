package balance

import (
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

func (r *Repository) GetUserBalance(userID uint) (current float64, withdrawn float64, err error) {
	err = r.db.Table("orders").
		Select("COALESCE(SUM(accrual), 0)").
		Where("user_id = ? AND status = ?", userID, "PROCESSED").
		Scan(&current).Error
	if err != nil {
		return
	}

	err = r.db.Table("withdrawals").
		Select("COALESCE(SUM(sum), 0)").
		Where("user_id = ?", userID).
		Scan(&withdrawn).Error
	return
}

func (r *Repository) CreateWithdrawal(userID uint, order string, sum float64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		withdrawal := db.Withdrawal{
			UserID:      userID,
			OrderNumber: order,
			Sum:         sum,
		}

		if err := tx.Create(&withdrawal).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE users SET balance = balance - ? WHERE id = ?", sum, userID).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *Repository) GetWithdrawals(userID uint) ([]db.Withdrawal, error) {
	var withdrawals []db.Withdrawal
	err := r.db.
		Where("user_id = ?", userID).
		Order("processed_at ASC").
		Find(&withdrawals).Error
	return withdrawals, err
}
