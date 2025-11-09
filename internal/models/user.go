package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Don't include in JSON responses
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Phone        string    `json:"phone"`
	UserType     string    `json:"user_type"` // patient, doctor, admin
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateUser inserts a new user into the database
func CreateUser(db *sql.DB, user *User) error {
	query := `
		INSERT INTO users (email, password_hash, first_name, last_name, phone, user_type)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := db.QueryRow(query, user.Email, user.PasswordHash, user.FirstName,
		user.LastName, user.Phone, user.UserType).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, email, password_hash, first_name, last_name, phone, user_type, created_at, updated_at
		FROM users WHERE email = $1
	`

	err := db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName,
		&user.LastName, &user.Phone, &user.UserType, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(db *sql.DB, id int) (*User, error) {
	user := &User{}
	query := `
		SELECT id, email, password_hash, first_name, last_name, phone, user_type, created_at, updated_at
		FROM users WHERE id = $1
	`

	err := db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName,
		&user.LastName, &user.Phone, &user.UserType, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetAllPatients retrieves all users with user_type = 'patient'
func GetAllPatients(db *sql.DB) ([]User, error) {
	query := `
		SELECT id, email, first_name, last_name, phone, user_type, created_at, updated_at
		FROM users WHERE user_type = 'patient'
		ORDER BY created_at DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName,
			&user.Phone, &user.UserType, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		patients = append(patients, user)
	}

	return patients, nil
}

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}
