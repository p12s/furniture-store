package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// Repository - repo
type Repository struct {
	Accounter
}

// NewRepository - constructor
func NewRepository(db *sqlx.DB) *Repository {
	createAccountTable(db)

	return &Repository{
		Accounter: NewAccount(db),
	}
}

// Deliberately removed the obligation of important fields (name, username, password_hash, ...),
// because the architecture is asynchronous, a Business-event with only role (role)
// can come before a CUD-event with all other data.

// createAccountTable
func createAccountTable(db *sqlx.DB) {
	query := `CREATE TABLE IF NOT EXISTS account (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"public_id" TEXT,
		"name" TEXT,
		"username" TEXT,		
		"password_hash" TEXT,
		"email" TEXT,
		"address" TEXT,
		"role" INTEGER DEFAULT 0,
		"created_at" DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
	  );`
	statement, err := db.Prepare(query)
	defer statement.Close() // nolint
	if err != nil {
		statement.Close()
		logrus.Fatal("create account.account table fail: ", err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		logrus.Fatal("exec creating account.account table fail: ", err.Error())
	}

	fmt.Println("account.account table created 🗂")
}
