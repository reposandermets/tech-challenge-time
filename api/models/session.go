package models

import (
	"database/sql"
	"errors"
)

// TimeSession ...
type TimeSession struct {
	TimeSessionPartialID    string         `json:"time_session_partial_id,omitempty"`
	TimeSessionName         sql.NullString `json:"time_session_name,omitempty"`
	TimeSessionPartialStart sql.NullString `json:"time_session_partial_start,omitempty"`
	TimeSessionPartialEnd   sql.NullString `json:"time_session_partial_end,omitempty"`
	TimeSessionID           string         `json:"time_session_id,omitempty"`
	TimeSessionCompleted    sql.NullBool   `json:"time_session_completed,omitempty"`
}

func scanRows(rows *sql.Rows) ([]*TimeSession, error) {
	timeSessions := make([]*TimeSession, 0)

	for rows.Next() {
		timeSession := new(TimeSession)
		err := rows.Scan(
			&timeSession.TimeSessionPartialID,
			&timeSession.TimeSessionName,
			&timeSession.TimeSessionPartialStart,
			&timeSession.TimeSessionPartialEnd,
			&timeSession.TimeSessionID,
			&timeSession.TimeSessionCompleted,
		)
		if err != nil {
			return nil, err
		}
		timeSessions = append(timeSessions, timeSession)
	}

	err := rows.Err()
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return timeSessions, nil
}

// ReadByID ...
func ReadByID(db *sql.DB, sessionPartialID string, userID string) ([]*TimeSession, error) {
	rows, err := db.Query(`
		SELECT
			time_session_partial_id,
			time_session_name,
			time_session_partial_start,
			time_session_partial_end,
			time_session_id,
			time_session_completed
		FROM time_session_partial
		WHERE time_session_partial_id = $1 AND user_id = $2`,
		sessionPartialID, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanRows(rows)
}

// ReadBySessionID ...
func ReadBySessionID(db *sql.DB, sessionID string, userID string) ([]*TimeSession, error) {
	rows, err := db.Query(`
		SELECT
			time_session_partial_id,
			time_session_name,
			time_session_partial_start,
			time_session_partial_end,
			time_session_id,
			time_session_completed
		FROM time_session_partial
		WHERE time_session_id = $1 AND user_id = $2`,
		sessionID, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanRows(rows)
}

// ReadOngoingTracking ...
func ReadOngoingTracking(db *sql.DB, userID string) ([]*TimeSession, error) {
	rows, err := db.Query(`
		SELECT
			time_session_partial_id,
			time_session_name,
			time_session_partial_start,
			time_session_partial_end,
			time_session_id,
			time_session_completed
		FROM time_session_partial
		WHERE time_session_partial_end IS NULL AND user_id = $1`,
		userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanRows(rows)
}

// ReadOngoingSession ...
func ReadOngoingSession(db *sql.DB, userID string) ([]*TimeSession, error) {
	rows, err := db.Query(`
		SELECT
			time_session_partial_id,
			time_session_name,
			time_session_partial_start,
			time_session_partial_end,
			time_session_id,
			time_session_completed
		FROM time_session_partial
		WHERE time_session_completed IS NULL AND user_id = $1`,
		userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanRows(rows)
}

// WriteStartSession new session is being started
func WriteStartSession(db *sql.DB, sessionPartialID string, name string, sessionID string, userID string) ([]*TimeSession, error) {
	timeSessions, err := ReadOngoingTracking(db, userID)
	if err != nil {
		return nil, err
	}
	if len(timeSessions) == 1 {
		return nil, errors.New("can't start session with ongoing tracking")
	}
	sqlStatement := `
		INSERT INTO time_session_partial (
			time_session_partial_id,
			time_session_name,
			time_session_partial_start,
			time_session_partial_end,
			time_session_id,
			time_session_completed,
			user_id
		)
		VALUES (
			$1,
			$2,
			NOW(),
			NULL,
			$3,
			NULL,
			$4
		);`

	_, err = db.Exec(sqlStatement, sessionPartialID, name, sessionID, userID)
	if err != nil {
		return nil, err
	}

	return ReadByID(db, sessionPartialID, userID)
}

// WriteEndSession Session will be completed
func WriteEndSession(db *sql.DB, sessionID string, userID string) ([]*TimeSession, error) {
	timeSessions, err := ReadOngoingTracking(db, userID)
	if err != nil {
		return nil, err
	}
	if len(timeSessions) == 1 {
		return nil, errors.New("can't end any session with ongoing tracking")
	}
	sqlStatement := `
	UPDATE time_session_partial 
	SET time_session_completed = true
	WHERE time_session_id = $1 AND user_id = $2;`
	_, err = db.Exec(sqlStatement, sessionID, userID)

	if err != nil {
		return nil, err
	}

	return ReadBySessionID(db, sessionID, userID)
}

// WriteStopSession user can add another time slot
func WriteStopSession(db *sql.DB, sessionPartialID string, userID string) ([]*TimeSession, error) {
	// session can be only stopped if there's ongoing tracking
	timeSessions, err := ReadOngoingTracking(db, userID)
	if err != nil {
		return nil, err
	}
	if len(timeSessions) != 1 {
		return nil, errors.New("no active tracking to stop")
	}
	sqlStatement := `
	UPDATE time_session_partial 
	SET time_session_partial_end = NOW()
	WHERE time_session_partial_id = $1 AND user_id = $2;`
	_, err = db.Exec(sqlStatement, sessionPartialID, userID)
	if err != nil {
		return nil, err
	}

	return ReadByID(db, sessionPartialID, userID)
}
