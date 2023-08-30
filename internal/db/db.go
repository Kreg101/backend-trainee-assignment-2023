package db

import (
	"context"
	"database/sql"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// User represents id and segments for adding, deleting, etc.
type User struct {
	Id       int64    `json:"id"`
	Segments []string `json:"segments,omitempty"`
}

// PostgresStore implements server.Storage interface
type PostgresStore struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

// NewStorage creates and checks connection to database
func NewStorage(config string, logger *zap.SugaredLogger) (*PostgresStore, error) {
	conn, err := sql.Open("pgx", config)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresStore{
		db:     conn,
		logger: logger,
	}, nil
}

// Init creates tables for database
func (s *PostgresStore) Init() error {
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		s.logger.Info("can't begin transaction", zap.Error(err))
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS users (
	 id BIGINT UNIQUE
	);`)

	if err != nil {
		s.logger.Info("can't create users table", zap.Error(err))
		return err
	}

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS segments  (
	 id SERIAL PRIMARY KEY,
	 name VARCHAR(50) UNIQUE NOT NULL
	);`)

	if err != nil {
		s.logger.Info("can't create segments table", zap.Error(err))
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
		s.logger.Info("can't create user_segments table", zap.Error(err))
		return err
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Info("can't commit transaction in init", zap.Error(err))
		return err
	}

	return nil
}

// CreateSegment creates new segment in database
func (s *PostgresStore) CreateSegment(name string) error {
	_, err := s.db.Exec(`INSERT INTO segments (name) VALUES ($1)`, name)
	return err
}

// DeleteSegment deletes segment from database
func (s *PostgresStore) DeleteSegment(name string) error {
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		s.logger.Info("can't begin transaction", zap.Error(err))
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`DELETE FROM user_segments
				WHERE segment_id IN (
				  SELECT id FROM segments
				  WHERE name = $1
				);`,
		name)

	if err != nil {
		s.logger.Info("can't delete segment from user_segments", zap.Error(err))
		return err
	}

	_, err = tx.Exec(
		`DELETE FROM segments
		  WHERE name = $1;`,
		name)

	if err != nil {
		s.logger.Info("can't delete from segments", zap.Error(err))
		return err
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Info("can't commit in deleteSegment", zap.Error(err))
	}

	return nil
}

// CreateUser creates new user in database and returns new id
func (s *PostgresStore) CreateUser(id int64) error {
	_, err := s.db.Exec(`INSERT INTO users (id) VALUES ($1);`, id)
	return err
}

// AddSegmentsToUser appends segments to existing user in database
func (s *PostgresStore) AddSegmentsToUser(user User) error {
	for _, name := range user.Segments {
		_, err := s.db.Exec(
			`INSERT INTO user_segments (user_id, segment_id)
					VALUES (
					   (SELECT id FROM users WHERE id = $1),
					   (SELECT id FROM segments WHERE name = $2)
					);`,
			user.Id, name)

		if err != nil {
			s.logger.Info("can't add segment to user", zap.Error(err))
		}
	}

	return nil
}

// DeleteSegmentsFromUser deletes segments from existing user in database
func (s *PostgresStore) DeleteSegmentsFromUser(user User) error {
	for _, name := range user.Segments {
		_, err := s.db.Exec(
			`DELETE FROM user_segments
					WHERE user_id = (SELECT id FROM users WHERE id = $1)
					AND segment_id = (SELECT id FROM segments WHERE name = $2);`,
			user.Id, name)

		if err != nil {
			s.logger.Info("can't delete segment from user", zap.Error(err))
		}
	}

	return nil
}

// GetUser gets user from database
func (s *PostgresStore) GetUser(id int64) (User, error) {
	user := User{
		Id:       id,
		Segments: make([]string, 0),
	}

	rows, err := s.db.Query(
		`SELECT us.user_id, s.name
				FROM user_segments us
				JOIN segments s ON us.segment_id = s.id
				WHERE us.user_id = $1;`,
		user.Id)

	if err != nil {
		s.logger.Info("can't get correct response from db", zap.Error(err))
		return User{}, err
	}

	for rows.Next() {
		var segment string
		err = rows.Scan(&user.Id, &segment)
		if err != nil {
			s.logger.Info("can't scan user-segment from storage", zap.Error(err))
			continue
		}
		user.Segments = append(user.Segments, segment)
	}

	if rows.Err() != nil {
		s.logger.Info("can't scan rows in getUser", zap.Error(err))
		return User{}, rows.Err()
	}

	return user, nil
}
