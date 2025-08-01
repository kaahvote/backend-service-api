package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"unicode/utf8"

	"github.com/kaahvote/backend-service-api/internal/validator"
)

type Session struct {
	ID                 int64     `json:"id"`
	Name               string    `json:"name"`
	PublicID           string    `json:"publicId"`
	ExpiresAt          time.Time `json:"expiresAt"`
	VotingPolicyID     int64     `json:"votingPolicyId"`
	VotersPolicyID     int64     `json:"votersPolicyId"`
	CandidatesPolicyID int64     `json:"candidatesPolicyId"`
	CreatedBy          int64     `json:"createdBy"`
	CreatedAt          time.Time `json:"createdAt"`
}

func ValidateSession(v *validator.Validator, s *Session) {
	v.Check(s.Name != "", "name", "must be provided")
	v.Check(utf8.RuneCountInString(s.Name) > 10, "name", "must at least 10 characters")

	v.Check(s.VotingPolicyID > 0, "votingPolicyID", "must be a valid positive integer")
	v.Check(s.VotersPolicyID > 0, "votersPolicyID", "must be a valid positive integer")
	v.Check(s.CandidatesPolicyID > 0, "candidatesPolicyID", "must be a valid positive integer")
	v.Check(s.CreatedBy > 0, "createdBy", "must be a valid positive integer")

	v.Check(s.ExpiresAt.After(time.Now()), "expiresAt", "cannot be in the past")
}

type SessionModel struct {
	DB *sql.DB
}

func (m SessionModel) Get(publicId string) (*Session, error) {
	query := `
		SELECT id, name, public_id, expires_at, voting_policy_id, 
		voters_policy_id, candidate_policy_id, created_by, created_at
		FROM sessions WHERE public_id=$1`

	var s Session

	ctx, cancel := context.WithTimeout(context.Background(), THREE_SECONDS)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, publicId).Scan(
		&s.ID,
		&s.Name,
		&s.PublicID,
		&s.ExpiresAt,
		&s.VotingPolicyID,
		&s.VotersPolicyID,
		&s.CandidatesPolicyID,
		&s.CreatedBy,
		&s.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &s, nil
}

func (m SessionModel) Insert(s *Session) error {

	ctx, cancel := context.WithTimeout(context.Background(), THREE_SECONDS)
	defer cancel()

	txn, err := m.DB.Begin()
	if err != nil {
		return err
	}

	err = m.InsertSessionWithTransaction(ctx, txn, s)
	if err != nil {
		return err
	}

	err = m.InsertFirstSessionFlowWithTransaction(ctx, txn, s.ID)
	if err != nil {
		return err
	}

	return txn.Commit()
}

func (m SessionModel) InsertSessionWithTransaction(ctx context.Context, txn *sql.Tx, s *Session) error {
	query := `INSERT INTO sessions (name, public_id, expires_at, voting_policy_id, voters_policy_id, candidate_policy_id, created_by) 
				VALUES ($1, $2, $3, $4, $5, $6, $7) 
				RETURNING id, created_at`

	stmt, err := txn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	args := []any{s.Name, s.PublicID, s.ExpiresAt, s.VotingPolicyID, s.VotersPolicyID, s.CandidatesPolicyID, s.CreatedBy}
	err = stmt.QueryRowContext(ctx, args...).Scan(&s.ID, &s.CreatedAt)
	if err != nil {
		return err
	}

	return stmt.Close()
}

func (m SessionModel) InsertFirstSessionFlowWithTransaction(ctx context.Context, txn *sql.Tx, sessionId int64) error {
	query := `INSERT INTO flows (session_id, state_id) VALUES ($1, $2)`

	args := []any{sessionId, SESSION_IN_DRAFT}
	stmt, err := txn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		return err
	}

	return stmt.Close()
}

func (m SessionModel) Update(s *Session) error {
	query := `UPDATE sessions SET 
			name=$1, expires_at=$2, voting_policy_id=$3, voters_policy_id=$4, candidate_policy_id=$5
			WHERE id = $6`

	ctx, cancel := context.WithTimeout(context.Background(), THREE_SECONDS)

	defer cancel()

	args := []any{s.Name, s.ExpiresAt, s.VotingPolicyID, s.VotersPolicyID, s.CandidatesPolicyID, s.ID}
	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m SessionModel) Delete(id int64) error {
	query := `DELETE FROM sessions WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), THREE_SECONDS)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

func (m SessionModel) ListSessionsByUserID(userID int64) ([]*Session, error) {
	query := `SELECT id, name, public_id, expires_at, 
				voting_policy_id, voters_policy_id, candidate_policy_id, 
				created_by, created_at 
			FROM sessions
			WHERE created_by = $1`

	ctx, cancel := context.WithTimeout(context.Background(), THREE_SECONDS)
	defer cancel()

	args := []any{userID}
	rows, err := m.DB.QueryContext(ctx, query, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		} else {
			return nil, err
		}
	}

	var sessions []*Session
	defer rows.Close()

	for rows.Next() {

		var session Session

		err := rows.Scan(
			&session.ID,
			&session.Name,
			&session.PublicID,
			&session.ExpiresAt,
			&session.VotingPolicyID,
			&session.VotersPolicyID,
			&session.CandidatesPolicyID,
			&session.CreatedBy,
			&session.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		sessions = append(sessions, &session)
	}

	return sessions, nil

}
