package balance

import (
	"errors"
	"gophermart/internal/db"
)

var (
	ErrInsufficientFunds = errors.New("not enough funds")
	ErrInvalidOrder      = errors.New("invalid order number")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

type BalanceReponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func (s *Service) GetBalance(userID uint) (*BalanceReponse, error) {
	current, withdrawn, err := s.repo.GetUserBalance(userID)
	if err != nil {
		return nil, err
	}

	return &BalanceReponse{
		Current:   current - withdrawn,
		Withdrawn: withdrawn,
	}, nil
}

func (s *Service) Withdraw(userID uint, order string, sum float64) error {
	current, withdrawn, err := s.repo.GetUserBalance(userID)
	if err != nil {
		return err
	}

	if (current - withdrawn) < sum {
		return ErrInsufficientFunds
	}

	return s.repo.CreateWithdrawal(userID, order, sum)
}

func (s *Service) GetWithdrawals(userID uint) ([]db.Withdrawal, error) {
	return s.repo.GetWithdrawals(userID)
}
