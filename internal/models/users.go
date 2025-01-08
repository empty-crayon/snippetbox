package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

// to add a new user to record to the users table
func (m *UserModel) Insert (name, email,password string) error {
	return nil
}

// Autheticate whether a user exists in the database, return relevant id if  yes or error
func (m *UserModel) Authenticate (email, password string) (int, error) {
	return 0, nil
} 

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}


