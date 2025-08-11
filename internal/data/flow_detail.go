package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrSessionUpdateNotAllowed = errors.New("session update is not allowed due to advanced flow status")
)

type FlowDetail struct {
	ID          int64     `json:"id"`
	Comment     string    `json:"comment,omitzero"`
	CreatedAt   time.Time `json:"createdAt"`
	StateDetail State     `json:"stateDetail"`
}

func (m FlowModel) GetCurrentFlow(sessionID int64) (*FlowDetail, error) {

	query := `SELECT
	  			f.id, s.id, s.name, f.comment, f.created_at
			  FROM flows f
			  INNER JOIN session_states s ON s.id = f.state_id
			  WHERE f.session_id = $1
			  ORDER BY created_at DESC
			  LIMIT 1`

	ctx, cancel := context.WithTimeout(context.Background(), THREE_SECONDS)
	defer cancel()

	var f FlowDetail
	var comment sql.NullString
	err := m.DB.QueryRowContext(ctx, query, sessionID).Scan(&f.ID, &f.StateDetail.ID, &f.StateDetail.Name, &comment, &f.CreatedAt)

	if err != nil {
		return nil, err
	}

	if comment.Valid {
		f.Comment = comment.String
	}

	return &f, nil
}

func (f *FlowDetail) IsInAdvancedVotingState() bool {
	return f.StateDetail.ID > 2
}

func (f *FlowDetail) IsFlowBackwarding(newStateId int64) bool {
	return f.StateDetail.ID > newStateId
}

func (f *FlowDetail) ShouldVotesBeDeleted(newStateId int64) bool {
	deleteVoteStates := []int64{SESSION_IN_DRAFT, SESSION_WAITING_CANDIDATES_OR_OPTIONS}

	if f.IsInAdvancedVotingState() {
		for _, v := range deleteVoteStates {
			if v == newStateId {
				return true
			}
		}
	}

	return false
}

func (f *FlowDetail) AllowSessionUpdate() (bool, error) {

	if f.IsInAdvancedVotingState() {
		return false, ErrSessionUpdateNotAllowed
	}

	return true, nil
}
