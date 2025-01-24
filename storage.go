package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type storage interface {
	CreateAccount(*Account) (int64, error)
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountById(int) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

// connecting to postgres
func NewPostgresStore() (*PostgresStore, error) {
	conn := "host=localhost port=5433 user=postgres password=gobank dbname=bank sslmode=disable"

	// connect to database
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Println("Error while connecting db!")
		return nil, err
	}

	// ping database
	if err = db.Ping(); err != nil {
		log.Println("Error while pinging db")
		return nil, err
	}

	log.Println("Database Successfully connected!")

	// return db connection
	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()

}

// create table accounts
func (s *PostgresStore) createAccountTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS	accounts(
			id	SERIAL PRIMARY KEY,
			firstName	VARCHAR(255) NOT NULL, 
			lastName	VARCHAR(255) NOT NULL, 
			number		BIGINT UNIQUE NOT NULL,
			balance		BIGINT NOT NULL DEFAULT 0,
			createdAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := s.db.Exec(query)
	if err != nil {
		log.Println("Error while creating table.")
		return err
	}

	log.Println("Table created successfully.")
	return nil

}

// create an Account
func (p *PostgresStore) CreateAccount(account *Account) (int64, error) {

	var lastInsertedId int64

	// query
	query := "INSERT INTO accounts(firstname, lastname, number, balance, createdat) VALUES($1, $2, $3, $4, $5) RETURNING ID"
	insertQueryErr := p.db.QueryRow(query, account.FirstName, account.LastName, account.Number, account.Balance, account.CreatedAt).Scan(&lastInsertedId)

	if insertQueryErr != nil {
		log.Println("Error while implementing insert query.")
		return 0, insertQueryErr
	}

	log.Println("Row inserted! ", lastInsertedId)
	return lastInsertedId, nil
}

func (p *PostgresStore) GetAccounts() ([]*Account, error) {

	var accounts []*Account
	query := "SELECT * FROM accounts"

	result, queryErr := p.db.Query(query)
	if queryErr != nil {
		log.Println("Error while querying the database")
		return []*Account{}, nil
	}

	for result.Next() {
		var account = &Account{}
		scanErr := result.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Number, &account.Balance, &account.CreatedAt)
		if scanErr != nil {
			return []*Account{}, nil
		}

		accounts = append(accounts, account)
	}

	// Print the accounts (for demonstration)
	for _, account := range accounts {
		fmt.Printf("ID: %d, Name: %s %s, Number: %d, Balance: %d, Created At: %v\n", account.ID, account.FirstName, account.LastName, account.Number, account.Balance, account.CreatedAt)
	}

	return accounts, nil
}

func (p *PostgresStore) DeleteAccount(id int) error {
	return nil
}
func (p *PostgresStore) UpdateAccount(account *Account) error {
	return nil
}
func (p *PostgresStore) GetAccountById(id int) (*Account, error) {
	return nil, nil
}
