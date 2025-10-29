package orders

import (
	"errors"
	"gophermart/internal/db"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetUserOrders(userID uint) ([]db.Order, error) {
	return s.repo.FindByUserID(userID)
}

func (s *Service) CreateOrder(userID uint, number string) (string, error) {
	//if !ValidLuhn(number) {
	//	return "", errors.New("Invalid number")
	//}

	exists, err := s.repo.FindByNumber(number)
	if err != nil {
		return "", err
	}
	if exists != nil {
		if exists.UserID == userID {
			return "own", nil
		}
		return "conflict", errors.New("order belongs to another user")
	}

	order := &db.Order{
		Number:     number,
		UserID:     userID,
		Status:     "NEW",
		UploadedAt: time.Now(),
	}

	if err := s.repo.Create(order); err != nil {
		return "", err
	}

	return "new", nil
}
