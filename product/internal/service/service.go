package service

import (
	_ "github.com/golang/mock/mockgen/model"

	"github.com/p12s/furniture-store/product/internal/config"
	"github.com/p12s/furniture-store/product/internal/repository"
)

//go:generate mockgen -destination mocks/mock.go -package service github.com/p12s/furniture-store/product/internal/service Accounter

// Service - just service
type Service struct {
	Accounter
}

// NewService - constructor
func NewService(repos *repository.Repository, config *config.Auth) *Service {
	return &Service{
		Accounter: NewAccountService(repos.Accounter, config),
	}
}
