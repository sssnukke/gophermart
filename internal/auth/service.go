package auth

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gophermart/internal/db"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Login(login, password string) (*db.User, error) {
	user, err := s.repo.FindByLogin(login)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("Invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("Invalid credentials")
	}

	return user, nil
}

func (s *Service) Register(login, password string) (*db.User, error) {
	exists, err := s.repo.FindByLogin(login)
	if err != nil {
		return nil, err
	}
	if exists != nil {
		return nil, errors.New("User already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &db.User{
		Login:        login,
		PasswordHash: string(hash),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
