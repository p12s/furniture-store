package domain

import (
	"time"

	"github.com/google/uuid"
)

// Role
type Role int

const (
	ROLE_CUSTOMER Role = iota
	ROLE_ADMIN
	ROLE_DELIVERY
	ROLE_DEALER
)

// Account
type Account struct {
	Id        int        `json:"id,omitempty" db:"id"`
	PublicId  uuid.UUID  `json:"public_id" db:"public_id"`
	Name      string     `json:"name" db:"name" binding:"required"`
	Username  string     `json:"username" db:"username" binding:"required"`
	Password  string     `json:"password,omitempty" db:"password_hash" binding:"required"`
	Email     string     `json:"email" db:"email" binding:"required"`
	Address   string     `json:"address" db:"address" binding:"required"` // TODO should be different columns - Country, City, Street, etc.
	Token     string     `json:"token,omitempty"`
	Role      Role       `json:"role" db:"role"`
	CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at"` // nolint
}

// SignInInput
type SignInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateAccountInput
type UpdateAccountInput struct {
	PublicId uuid.UUID `json:"public_id" db:"public_id" binding:"required"`
	Name     *string   `json:"name" db:"role"`
	Username *string   `json:"username" db:"username"`
	Password *string   `json:"password,omitempty" db:"role"`
	Email    *string   `json:"email" db:"email"`
	Address  *string   `json:"address" db:"address"`
}

// UpdateAccountRoleInput
type UpdateAccountRoleInput struct {
	PublicId uuid.UUID `json:"public_id" db:"public_id" binding:"required"`
	Role     Role      `json:"role" db:"role" binding:"required"`
}

// AccountToken
type AccountToken struct {
	PublicId uuid.UUID `json:"public_id" db:"public_id" binding:"required"`
	Token    string    `json:"token" db:"token" binding:"token"`
}

// DeleteAccountInput
type DeleteAccountInput struct {
	PublicId string `json:"public_id" db:"public_id" binding:"required"`
}

// EventType
type EventType string

const (
	EVENT_ACCOUNT_CREATED       EventType = "auth.created"
	EVENT_ACCOUNT_INFO_UPDATED  EventType = "auth.info_updated"
	EVENT_ACCOUNT_ROLE_UPDATED  EventType = "auth.role_updated"
	EVENT_ACCOUNT_DELETED       EventType = "auth.deleted"
	EVENT_ACCOUNT_TOKEN_UPDATED EventType = "auth.token_updated" // nolint
)

// Event
type Event struct {
	Type  EventType
	Value interface{}
}
