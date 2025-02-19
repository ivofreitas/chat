package postgres

import (
	"database/sql"
	"fmt"
	"github.com/ivofreitas/chat/pkg/config"
	"log"

	_ "github.com/lib/pq"
)

func NewConnection() *sql.DB {
	env := config.GetEnv()
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		env.Database.Host, env.Database.Port, env.Database.User,
		env.Database.Password, env.Database.DBName, env.Database.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	return db
}
