package db

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// PostgresStore implements server.Storage interface
type PostgresStore struct {
	db *sql.DB
}

// NewStorage creates and checks connection to database
func NewStorage(config string) (*PostgresStore, error) {
	conn, err := sql.Open("pgx", config)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: conn,
	}, nil
}

// Init creates tables for database
func (s *PostgresStore) Init() error {
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS users (
	 id BIGINT UNIQUE
	);`)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS segments  (
	 id SERIAL PRIMARY KEY,
	 name VARCHAR(50) UNIQUE NOT NULL
	);`)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS user_segments (
	 user_id INT,
	 segment_id INT,
	 PRIMARY KEY (user_id, segment_id),
	 FOREIGN KEY (user_id) REFERENCES users(id),
	 FOREIGN KEY (segment_id) REFERENCES segments(id)
	);`)

	if err != nil {
		return err
	}

	return tx.Commit()
}

// CreateSegment creates new segment in database
func (s *PostgresStore) CreateSegment(name string) error {
	return nil
}

// DeleteSegment deletes segment from database
func (s *PostgresStore) DeleteSegment(name string) error {
	return nil
}

// CreateUser creates new user in database and returns new id
func (s *PostgresStore) CreateUser(id int64) (int64, error) {
	_, err := s.db.Exec(`INSERT INTO users (id) VALUES ($1);`, id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// UpdateUser updates existing user in database
func (s *PostgresStore) UpdateUser(user User) error {
	return nil
}

// GetUser gets user from database
func (s *PostgresStore) GetUser(id int64) (*User, error) {
	return nil, nil
}
