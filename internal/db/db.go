package db

import (
	"context"
	"database/sql"
	"fmt"

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
	 user_id BIGINT NOT NULL,
	 segment_id INT NOT NULL,
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
	_, err := s.db.Exec(`INSERT INTO segments (name) VALUES ($1)`, name)
	return err
}

// DeleteSegment deletes segment from database
func (s *PostgresStore) DeleteSegment(name string) error {
	return nil
}

// CreateUser creates new user in database and returns new id
func (s *PostgresStore) CreateUser(id int64) error {
	_, err := s.db.Exec(`INSERT INTO users (id) VALUES ($1);`, id)
	return err
}

// UpdateUser updates existing user in database
// TODO handle errors
func (s *PostgresStore) UpdateUser(user User) error {
	for _, name := range user.AppendSegments {
		_, err := s.db.Exec(
			`INSERT INTO user_segments (user_id, segment_id)
					VALUES (
					   (SELECT id FROM users WHERE id = $1),
					   (SELECT id FROM segments WHERE name = $2)
					);`, user.Id, name)
		if err != nil {
			fmt.Printf("can't append segment to user %v\n", err)
		}
	}

	for _, name := range user.DeleteSegments {
		_, err := s.db.Exec(
			`DELETE FROM user_segments
					WHERE user_id = (SELECT id FROM users WHERE id = $1)
					AND segment_id = (SELECT id FROM segments WHERE name = $2);`,
			user.Id, name)
		if err != nil {
			fmt.Printf("can't delete segment from user %v\n", err)
		}
	}

	return nil
}

// GetUser gets user from database
func (s *PostgresStore) GetUser(id int64) (*User, error) {
	return nil, nil
}
