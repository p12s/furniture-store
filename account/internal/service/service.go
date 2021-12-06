package service

import (
	"time"

	"github.com/p12s/furniture-store/account/internal/repository"
)

// Service - just service
type Service struct {
	Accounter
}

type AccountConfig struct {
	Salt       string
	TokenTTL   time.Duration
	SigningKey string
}

// NewService - constructor
func NewService(repos *repository.Repository, config *AccountConfig) *Service {
	return &Service{
		Accounter: NewAccountService(repos.Accounter, config),
	}
}
