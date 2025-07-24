package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

const THREE_SECONDS = 3 * time.Second

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
	query := `INSERT INTO sessions (name, public_id, expires_at, voting_policy_id, voters_policy_id, candidate_policy_id, created_by) VALUES
				($1, $2, $3, $4, $5, $6, $7) 
				RETURNING id, created_at`

	args := []any{s.Name, s.PublicID, s.ExpiresAt, s.VotingPolicyID, s.VotersPolicyID, s.CandidatesPolicyID, s.CreatedBy}

	ctx, cancel := context.WithTimeout(context.Background(), THREE_SECONDS)

	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&s.ID, &s.CreatedAt)
}
