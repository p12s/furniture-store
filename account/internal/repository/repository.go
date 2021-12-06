package repository

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
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
	if err != nil {
		log.Fatal("create account.account table error", err.Error())
	}
	statement.Exec()
	fmt.Println("account.account table created 🗂")
}
