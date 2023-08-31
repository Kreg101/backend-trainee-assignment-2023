package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Kreg101/backend-trainee-assignment-2023/internal/server"
	"go.uber.org/zap"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// PostgresStore implements server.Storage interface
type PostgresStore struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

// NewStorage creates and checks connection to database
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
	       segment_id INT NOT NULL,
	       time_added BIGINT NOT NULL,
	       time_removed BIGINT,
	       FOREIGN KEY (user_id) REFERENCES users(id),
	       FOREIGN KEY (segment_id) REFERENCES segments(id)
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

// CreateSegment creates new segment in database
func (s *PostgresStore) CreateSegment(segment server.Segment) error {
	// check autoPercent validity
	if segment.AutoPercent < 0 || segment.AutoPercent > 100 {
		return errors.New("invalid percent")
	}

	// should we use auto addition?
	if segment.AutoPercent == 0 {
		_, err := s.db.Exec(`INSERT INTO segments (name) VALUES ($1)`, segment.Name)
		if err != nil {
			s.logger.Info(zap.Error(err))
			return err
		}
		return nil
	}

	// begin transaction with auto addition
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return err
	}
	tx.Rollback()

	// insert segment to database
	_, err = tx.Exec(`INSERT INTO segments (name) VALUES ($1)`, segment.Name)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return err
	}

	// count users
	row := tx.QueryRow(`SELECT COUNT(*) FROM users`)
	var all int64
	if err = row.Scan(&all); err != nil {
		s.logger.Info(zap.Error(err))
		return err
	}

	// calculate amount of users to add segment to
	count := (all * int64(segment.AutoPercent)) / 100

	// get users id
	rows, err := tx.Query(
		`SELECT (id) FROM users
			   ORDER BY RAND() LIMIT $1;`,
		count)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return err
	}

	// scan all users
	ids := make([]int64, 0)
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			s.logger.Info(zap.Error(err))
			return err
		}
		ids = append(ids, id)
	}

	if row.Err() != nil {
		s.logger.Info(zap.Error(err))
		return row.Err()
	}

	// add user-segment to users_segments and user_segment_history database
	for _, id := range ids {
		err := dbAddUserSegment(tx, id, segment.Name, nil)
		if err != nil {
			s.logger.Info(zap.Error(err))
			return err
		}
	}

	return nil
}

// DeleteSegment deletes segment by it's name 	from database
func (s *PostgresStore) DeleteSegment(name string) error {
	// use transactions because we need to delete segment from
	// two different tables (segment, user_segment)
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return err
	}
	defer tx.Rollback()

	// delete segments from user_segments
	_, err = tx.Exec(
		`DELETE FROM user_segments WHERE segment_id IN (
			   SELECT id FROM segments WHERE name = $1);`,
		name)

	if err != nil {
		s.logger.Info(zap.Error(err))
		return err
	}

	// TODO: solve bug

	// update history connected to this segment
	_, err = tx.Exec(
		`UPDATE user_segment_history 
               SET time_removed = $1
               WHERE time_removed > $2 AND
               segment_id = (SELECT id FROM segments WHERE name = $3
			  );`,
		time.Now().Unix(), time.Now().Unix(), name)

	if err != nil {
		s.logger.Info(zap.Error(err))
		return err
	}

	// delete segment from user
	_, err = tx.Exec(
		`DELETE FROM segments WHERE name = $1;`,
		name)

	if err != nil {
		s.logger.Info(zap.Error(err))
		return err
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Info(zap.Error(err))
		return err
	}

	return nil
}

// CreateUser creates new user in database and returns new id
func (s *PostgresStore) CreateUser(id int64) error {
	_, err := s.db.Exec(`INSERT INTO users (id) VALUES ($1);`, id)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return err
	}
	return nil
}

// AddSegmentsToUser appends segments to existing user in database
func (s *PostgresStore) AddSegmentsToUser(user server.User) error {
	deleteTime := new(int64)
	now := time.Now().Unix()

	// check that active time is valid
	if user.ActiveTime > 0 {
		*deleteTime = now + user.ActiveTime
	} else {
		deleteTime = nil
	}

	// begin transaction
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return err
	}
	defer tx.Rollback()

	// add all user-segment pairs
	for _, name := range user.Segments {
		err := dbAddUserSegment(tx, user.Id, name, deleteTime)
		if err != nil {
			s.logger.Info(zap.Error(err))
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Info(zap.Error(err))
		return err
	}

	return nil
}

// DeleteSegmentsFromUser deletes segments from existing user in database
func (s *PostgresStore) DeleteSegmentsFromUser(user server.User) error {
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return err
	}
	defer tx.Rollback()

	// delete all user-segment pairs
	for _, name := range user.Segments {

		// delete from user_segment table
		_, err := tx.Exec(
			`DELETE FROM user_segments
				   WHERE user_id = (SELECT id FROM users WHERE id = $1)
				   AND segment_id = (SELECT id FROM segments WHERE name = $2);`,
			user.Id, name)

		if err != nil {
			s.logger.Info(zap.Error(err))
			return err
		}

		// update history
		_, err = tx.Exec(
			`UPDATE user_segment_history 
				   SET time_removed = $1
                   WHERE user_id = $2
                   AND segment_id = (SELECT id FROM segments
                   WHERE name = $3
				   );`,
			time.Now().Unix(), user.Id, name)

		if err != nil {
			s.logger.Info(zap.Error(err))
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Info(zap.Error(err))
		return err
	}

	return nil
}

// GetUser gets user from database. If user is in database, it returns (*User, nil).
// If error happened, it returns (nil, error). If user isn't in db, it returns (nil, nil).
func (s *PostgresStore) GetUser(id int64) (*server.User, error) {
	user := &server.User{Id: id, Segments: make([]string, 0)}

	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return nil, err
	}
	defer tx.Rollback()

	// check that user exists
	row := tx.QueryRow(
		`SELECT EXISTS ( SELECT 1 FROM users
     	       WHERE id = $1) AS user_exists;`, id)

	var exists bool
	err = row.Scan(&exists)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return nil, err
	}

	// if it doesn't exist it's NotFount
	if !exists {
		return nil, nil
	}

	// get all pairs
	rows, err := tx.Query(
		`SELECT us.user_id, s.name
			   FROM user_segments us
			   JOIN segments s ON us.segment_id = s.id
			   WHERE us.user_id = $1;`,
		user.Id)

	if err != nil {
		s.logger.Info(zap.Error(err))
		return nil, err
	}

	// scan all pairs
	for rows.Next() {
		var segment string
		err = rows.Scan(&user.Id, &segment)
		if err != nil {
			s.logger.Info(zap.Error(err))
			continue
		}
		user.Segments = append(user.Segments, segment)
	}

	// check error
	if rows.Err() != nil {
		s.logger.Info(zap.Error(err))
		return nil, rows.Err()
	}

	return user, nil
}

func (s *PostgresStore) GetUserHistory(user server.User) ([]server.TimeUser, error) {
	firstOfMonth := time.Date(user.Year, time.Month(user.Month), 0, 0, 0, 0, 0, time.UTC)

	// should be careful with +- 1 days, so longer not shorter
	lastOfMonth := firstOfMonth.AddDate(0, 1, 1)
	start := firstOfMonth.Unix()
	end := lastOfMonth.Unix()

	// get (id, segment_name, time_in, time_out) user from start to end
	rows, err := s.db.Query(
		`SELECT ush.user_id, s.name AS segment_name, ush.time_added, ush.time_removed
			   FROM user_segment_history ush
			   JOIN segments s ON s.id = ush.segment_id
			   WHERE ush.user_id = $1 AND ush.time_added > $2 AND ush.time_removed < $3`,
		user.Id, start, end)

	if err != nil {
		s.logger.Info(zap.Error(err))
		return nil, err
	}

	history := make([]server.TimeUser, 0)

	// scan all rows
	for rows.Next() {
		var tu server.TimeUser

		err := rows.Scan(&tu.Id, &tu.SegmentName, &start, &end)
		if err != nil {
			s.logger.Info(zap.Error(err))
			return nil, err
		}

		tu.TimeIn = time.Unix(start, 0).String()
		tu.TimeOut = time.Unix(end, 0).String()

		history = append(history, tu)
	}

	if rows.Err() != nil {
		s.logger.Info(zap.Error(rows.Err()))
		return nil, err
	}

	return history, nil
}

// ttl every 30 seconds finds expired records in user_segment database
// and delete them
func (s *PostgresStore) ttl() {
	for range time.Tick(30 * time.Second) {
		_, err := s.db.Exec(
			`DELETE FROM user_segments 
       			   WHERE time_out IS NOT NULL AND time_out <= $1`,
			time.Now().Unix())

		if err != nil {
			s.logger.Info(zap.Error(err))
		}
	}
}

// dbAddUserSegment for adding user-segment pair to both tables
func dbAddUserSegment(tx *sql.Tx, id int64, name string, deleteTime *int64) error {
	_, err := tx.Exec(
		`INSERT INTO user_segments (user_id, segment_id, time_in, time_out)
				   VALUES (
				   (SELECT id FROM users WHERE id = $1),
				   (SELECT id FROM segments WHERE name = $2), $3, $4);`,
		id, name, time.Now().Unix(), deleteTime)

	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`INSERT INTO user_segment_history (user_id, segment_id, time_added, time_removed)
				   VALUES ($1, (SELECT id FROM segments WHERE name = $2), $3, $4)`,
		id, name, time.Now().Unix(), deleteTime)

	if err != nil {
		return err
	}

	return nil
}
