package repository

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/p12s/furniture-store/product/internal/domain"
)

type Accounter interface {
	CreateAccount(account domain.Account) error
	UpdateAccountInfo(input domain.UpdateAccountInput) error
	UpdateAccountRole(input domain.UpdateAccountRoleInput) error
	DeleteAccount(accountPublicId uuid.UUID) error
	GetByCredentials(email, password string) (domain.Account, error)
}

type Account struct {
	db *sqlx.DB
}

func NewAccount(db *sqlx.DB) *Account {
	return &Account{db: db}
}

func (r *Account) CreateAccount(account domain.Account) error {
	query := fmt.Sprintf(`INSERT INTO %s (public_id, name, username, password_hash, email, address, role)
		values ($1, $2, $3, $4, $5, $6, $7)`, accountTable)
	_, err := r.db.Exec(query, account.PublicId, account.Name,
		account.Username, account.Password, account.Email, account.Address, domain.ROLE_CUSTOMER)
	return err
}

func (r *Account) UpdateAccountInfo(input domain.UpdateAccountInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.Username != nil {
		setValues = append(setValues, fmt.Sprintf("username=$%d", argId))
		args = append(args, *input.Username)
		argId++
	}

	if input.Password != nil {
		setValues = append(setValues, fmt.Sprintf("password_hash=$%d", argId))
		args = append(args, *input.Password)
		argId++
	}

	if input.Email != nil {
		setValues = append(setValues, fmt.Sprintf("email=$%d", argId))
		args = append(args, *input.Email)
		argId++
	}

	if input.Address != nil {
		setValues = append(setValues, fmt.Sprintf("address=$%d", argId))
		args = append(args, *input.Address)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE public_id = $%d`,
		accountTable, setQuery, argId)
	args = append(args, input.PublicId.String())

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *Account) UpdateAccountRole(input domain.UpdateAccountRoleInput) error {
	query := fmt.Sprintf(`UPDATE %s SET role=$1 WHERE public_id = $2`, accountTable)
	_, err := r.db.Exec(query, input.Role, input.PublicId)
	return err
}

func (r *Account) DeleteAccount(accountPublicId uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE public_id = $1`, accountTable)
	_, err := r.db.Exec(query, accountPublicId.String())
	return err
}

func (r *Account) GetByCredentials(email, password string) (domain.Account, error) {
	var account domain.Account

	query := fmt.Sprintf(`SELECT * FROM %s WHERE email=$1 AND password_hash=$2`, accountTable)
	err := r.db.Get(&account, query, email, password)
	if err != nil {
		return account, fmt.Errorf("get account: %w", err)
	}

	return account, err
}
