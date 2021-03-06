package service

import (
	"crypto/sha1" // nolint
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/p12s/furniture-store/account/internal/config"
	"github.com/p12s/furniture-store/account/internal/domain"
	"github.com/p12s/furniture-store/account/internal/repository"
)

var _ Accounter = (*AccountService)(nil)

// Accounter - service interface
type Accounter interface {
	CreateAccount(account domain.Account) error
	GetAccount(publicId string) (domain.Account, error)
	UpdateAccountInfo(input domain.UpdateAccountInput) error
	UpdateAccountRole(input domain.UpdateAccountRoleInput) error
	DeleteAccount(accountPublicId string) error
	GenerateTokenByCreds(email, password string) (string, error)
	ParseToken(token string) (string, error)
}

// AccountService - service
type AccountService struct {
	repo       repository.Accounter
	salt       string
	tokenTTL   time.Duration
	signingKey string
}

// NewAccountService - constructor
func NewAccountService(repo repository.Accounter, config *config.Auth) *AccountService {
	return &AccountService{
		repo:       repo,
		salt:       config.Salt,
		tokenTTL:   time.Duration(config.TokenTTL * 1000000),
		signingKey: config.SigningKey,
	}
}

// CreateAccount
func (s *AccountService) CreateAccount(account domain.Account) error {
	account.PublicId = uuid.New()
	account.Role = domain.ROLE_CUSTOMER
	passwordHash, err := s.generatePasswordHash(account.Password)
	if err != nil {
		return fmt.Errorf("generate password: %w", err)
	}
	account.Password = passwordHash
	return s.repo.CreateAccount(account)
}

// GetAccount
func (s *AccountService) GetAccount(publicId string) (domain.Account, error) {
	return s.repo.GetAccount(publicId)
}

// UpdateAccountInfo
func (s *AccountService) UpdateAccountInfo(input domain.UpdateAccountInput) error {
	if input.Password != nil {
		passwordHash, err := s.generatePasswordHash(*input.Password)
		if err != nil {
			return fmt.Errorf("generate password: %w", err)
		}
		*input.Password = passwordHash
	}
	return s.repo.UpdateAccountInfo(input)
}

// UpdateAccountRole
func (s *AccountService) UpdateAccountRole(input domain.UpdateAccountRoleInput) error {
	return s.repo.UpdateAccountRole(input)
}

// DeleteAccount
func (s *AccountService) DeleteAccount(accountPublicId string) error {
	return s.repo.DeleteAccount(accountPublicId)
}

// GenerateTokenByCreds
func (s *AccountService) GenerateTokenByCreds(email, password string) (string, error) {
	passwordHash, err := s.generatePasswordHash(password)
	if err != nil {
		return "", fmt.Errorf("generate password: %w", err)
	}
	account, err := s.repo.GetByCredentials(email, passwordHash)
	if err != nil {
		return "", fmt.Errorf("user creds wrong: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{ //nolint
		Subject:   account.PublicId.String(),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
	})

	return token.SignedString([]byte(s.signingKey))
}

// ParseToken
func (s *AccountService) ParseToken(accessToken string) (string, error) {
	t, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.signingKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("unexpected signing method: %w/n", err)
	}

	if !t.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("invalid subject")
	}

	return subject, nil
}

// generatePasswordHash
func (s *AccountService) generatePasswordHash(password string) (string, error) {
	hash := sha1.New() // #nosec
	if _, err := hash.Write([]byte(password)); err != nil {
		return "", fmt.Errorf("hash write: %w", err)
	}
	return fmt.Sprintf("%x", hash.Sum([]byte(s.salt))), nil
}
