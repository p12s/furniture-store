package repository

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/p12s/furniture-store/account/internal/domain"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestAccount_CreateAccount(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	assert.Equal(t, nil, err)
	defer db.Close()

	repo := NewAccount(db)

	publicId := "265cee57-2ff9-4ed3-85e1-d3373fa2a1a5"
	uuidPublicId, err := uuid.Parse(publicId)
	assert.Equal(t, nil, err)

	type args struct {
		account domain.Account
	}
	type mockBehavior func(args args)

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantErr      bool
	}{
		{
			name: "Can create account with right input",
			mockBehavior: func(args args) {
				mock.ExpectExec("INSERT INTO "+accountTable).WithArgs(args.account.PublicId,
					args.account.Name, args.account.Username, args.account.Password,
					args.account.Email, args.account.Address, domain.ROLE_CUSTOMER).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			args: args{
				account: domain.Account{
					PublicId: uuidPublicId,
					Name:     "Ivan",
					Username: "ivan",
					Password: "qwerty",
					Email:    "test@test.ru",
					Address:  "Some-city, some-street, some-hause",
				},
			},
		},
		{
			name: "Can't create account without required name-field",
			mockBehavior: func(args args) {
				mock.ExpectExec("INSERT INTO "+accountTable).WithArgs(args.account.PublicId,
					args.account.Name, args.account.Username, args.account.Password,
					args.account.Email, args.account.Address, domain.ROLE_CUSTOMER).
					WillReturnError(errors.New("some error"))
			},
			args: args{
				account: domain.Account{
					PublicId: uuidPublicId,
					Username: "ivan",
					Password: "qwerty",
					Email:    "test@test.ru",
					Address:  "Some-city, some-street, some-hause",
				},
			},
			wantErr: true,
		},
		{
			name: "Can't create account without required public-id-field",
			mockBehavior: func(args args) {
				mock.ExpectExec("INSERT INTO "+accountTable).WithArgs(args.account.PublicId,
					args.account.Name, args.account.Username, args.account.Password,
					args.account.Email, args.account.Address, domain.ROLE_CUSTOMER).
					WillReturnError(errors.New("some error"))
			},
			args: args{
				account: domain.Account{
					Username: "ivan",
					Password: "qwerty",
					Email:    "test@test.ru",
					Address:  "Some-city, some-street, some-hause",
				},
			},
			wantErr: true,
		},
		{
			name: "Can't create account without required username-field",
			mockBehavior: func(args args) {
				mock.ExpectExec("INSERT INTO "+accountTable).WithArgs(args.account.PublicId,
					args.account.Name, args.account.Username, args.account.Password,
					args.account.Email, args.account.Address, domain.ROLE_CUSTOMER).
					WillReturnError(errors.New("some error"))
			},
			args: args{
				account: domain.Account{
					Password: "qwerty",
					Email:    "test@test.ru",
					Address:  "Some-city, some-street, some-hause",
				},
			},
			wantErr: true,
		},
		{
			name: "Can't create account without required password-field",
			mockBehavior: func(args args) {
				mock.ExpectExec("INSERT INTO "+accountTable).WithArgs(args.account.PublicId,
					args.account.Name, args.account.Username, args.account.Password,
					args.account.Email, args.account.Address, domain.ROLE_CUSTOMER).
					WillReturnError(errors.New("some error"))
			},
			args: args{
				account: domain.Account{
					Email:   "test@test.ru",
					Address: "Some-city, some-street, some-hause",
				},
			},
			wantErr: true,
		},
		{
			name: "Can't create account without required email-field",
			mockBehavior: func(args args) {
				mock.ExpectExec("INSERT INTO "+accountTable).WithArgs(args.account.PublicId,
					args.account.Name, args.account.Username, args.account.Password,
					args.account.Email, args.account.Address, domain.ROLE_CUSTOMER).
					WillReturnError(errors.New("some error"))
			},
			args: args{
				account: domain.Account{
					Address: "Some-city, some-street, some-hause",
				},
			},
			wantErr: true,
		},
		{
			name: "Can't create account without required address-field",
			mockBehavior: func(args args) {
				mock.ExpectExec("INSERT INTO "+accountTable).WithArgs(args.account.PublicId,
					args.account.Name, args.account.Username, args.account.Password,
					args.account.Email, args.account.Address, domain.ROLE_CUSTOMER).
					WillReturnError(errors.New("some error"))
			},
			args: args{
				account: domain.Account{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockBehavior(tt.args)

			err := repo.CreateAccount(tt.args.account)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAccount_GetAccount(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	assert.Equal(t, nil, err)
	defer db.Close()

	repo := NewAccount(db)

	publicId := "265cee57-2ff9-4ed3-85e1-d3373fa2a1a5"
	uuidPublicId, err := uuid.Parse(publicId)
	assert.Equal(t, nil, err)

	type args struct {
		account  domain.Account
		publicId string
	}
	type mockBehavior func(args args)

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantErr      bool
	}{
		{
			name: "Can get account with right publicId",
			mockBehavior: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "public_id", "name", "username", "password_hash", "email", "address", "role"}).
					AddRow(1, args.account.PublicId, args.account.Name, args.account.Username,
						args.account.Password, args.account.Email, args.account.Address, domain.ROLE_CUSTOMER)

				mock.ExpectQuery("^SELECT (.+) FROM " + accountTable + " WHERE public_id=").WithArgs(args.publicId).WillReturnRows(rows)
			},
			args: args{
				publicId: publicId,
				account: domain.Account{
					Id:       1,
					PublicId: uuidPublicId,
					Name:     "Ivan",
					Username: "ivan",
					Password: "qwerty",
					Email:    "test@test.ru",
					Address:  "Some-city, some-street, some-hause",
				},
			},
		},
		{
			name: "Can't get account without empty publicId",
			mockBehavior: func(args args) {
				mock.ExpectQuery("^SELECT (.+) FROM " + accountTable + " WHERE public_id=").WithArgs(args.publicId).WillReturnError(fmt.Errorf("some error"))
			},
			args:    args{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockBehavior(tt.args)

			account, err := repo.GetAccount(tt.args.publicId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.args.account, account)
			}
		})
	}
}
