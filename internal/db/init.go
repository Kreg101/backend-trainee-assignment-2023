package db

import (
	"context"
	"database/sql"
	"go.uber.org/zap"
)

// PostgresStore implements server.Storage interface
type PostgresStore struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

// NewStorage creates and checks connection to database Postgresql
func NewStorage(config string, logger *zap.SugaredLogger) (*PostgresStore, error) {
	// try to connect to database
	conn, err := sql.Open("pgx", config)
	if err != nil {
		return nil, err
	}

	// check connection
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

	//use transactions for creating all tables
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return err
	}
	defer tx.Rollback()

	// basic user table
	_, err = tx.Exec(
		`CREATE TABLE IF NOT EXISTS users (
		   id BIGINT UNIQUE CHECK (id > 0)
	);`)

	if err != nil {
		s.logger.Info(zap.Error(err))
		return err
	}

	// segments table, every segment name has his own id for easier connection
	// with user
	_, err = tx.Exec(
		`CREATE TABLE IF NOT EXISTS segments  (
	 	   id SERIAL PRIMARY KEY,
	       name VARCHAR(50) UNIQUE NOT NULL
	);`)

	if err != nil {
		s.logger.Info(zap.Error(err))
		return err
	}

	// user_segments table for connection users and segments
	// time_in - the time the record was added to the database
	// time_out - the time the record should be deleted from database
	// time_out = NULL means the record can only be removed manually
	_, err = tx.Exec(
		`CREATE TABLE IF NOT EXISTS user_segments (
		   user_id BIGINT NOT NULL,
		   segment_id INT NOT NULL,
		   time_in BIGINT NOT NULL,
		   time_out BIGINT,
		   PRIMARY KEY (user_id, segment_id),
		   FOREIGN KEY (user_id) REFERENCES users(id),
		   FOREIGN KEY (segment_id) REFERENCES segments(id)
	);`)

	if err != nil {
		s.logger.Info(zap.Error(err))
		return err
	}

	// user_segment_history table contains all information about user-segment
	// relationships with start and finish time
	_, err = tx.Exec(
		`CREATE TABLE IF NOT EXISTS user_segment_history (
	   	   user_id BIGINT NOT NULL,
	       segment_name VARCHAR(50) NOT NULL,
	       time_added BIGINT NOT NULL,
	       time_removed BIGINT,
	       FOREIGN KEY (user_id) REFERENCES users(id)
	  );`)

	if err != nil {
		s.logger.Info(zap.Error(err))
		return err
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Info(zap.Error(err))
		return err
	}

	// start ttl function for automatic removal
	go s.ttl()

	return nil
}
