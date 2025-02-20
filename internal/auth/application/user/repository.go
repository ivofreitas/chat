package user

import (
	"database/sql"
	"errors"
	"github.com/ivofreitas/chat/internal/auth/domain"
	"log"
)

type Repository interface {
	CreateUser(email, hashedPassword string) error
	GetUserByEmail(email string) (*domain.User, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// CreateUser inserts a new user into the database.
func (r *repository) CreateUser(email, hashedPassword string) error {
	query := "INSERT INTO users_schema.users (email, password) VALUES ($1, $2)"
	_, err := r.db.Exec(query, email, hashedPassword)
	if err != nil {
		log.Println("Error inserting user:", err)
		return err
	}
	return nil
}

// GetUserByEmail retrieves a user by their email.
func (r *repository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	query := "SELECT id, email, password FROM users_schema.users WHERE email = $1"
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.HashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Println("Error retrieving user:", err)
		return nil, err
	}
	return &user, nil
}
