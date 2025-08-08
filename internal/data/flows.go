package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Flow struct {
	ID              int64     `json:"id"`
	SessionID       int64     `json:"-"`
	SessionPublicID string    `json:"sessionId"`
	StateID         int64     `json:"stateId"`
	StateName       string    `json:"stateName,omitzero"`
	Comment         string    `json:"comment,omitzero"`
	CreatedAt       time.Time `json:"createdAt"`
}

type FlowDetail struct {
	ID          int64     `json:"id"`
	Comment     string    `json:"comment,omitzero"`
	CreatedAt   time.Time `json:"createdAt"`
	StateDetail State     `json:"stateDetail"`
}

func (f Flow) Equals(flow *Flow) bool {
	return f.SessionID == flow.SessionID &&
		f.StateID == flow.StateID
}

type FlowModel struct {
	DB *sql.DB
}

func (m FlowModel) Insert(f *Flow) error {
	query := `INSERT INTO flows (session_id, state_id, comment)
			  VALUES ($1, $2, $3)
			  RETURNING id, created_at`
	args := []any{f.SessionID, f.StateID, f.Comment}

	ctx, cancel := context.WithTimeout(context.Background(), THREE_SECONDS)

	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&f.ID, &f.CreatedAt)
}

func (m FlowModel) InsertFirstFlow(sessionId int64) error {
	flow := &Flow{
		SessionID: sessionId,
		StateID:   1,
		Comment:   "",
	}
	return m.Insert(flow)
}

func (m FlowModel) UpdateState(f *Flow) error {
	query := `UPDATE flows SET state_id = $1 WHERE id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), THREE_SECONDS)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, f.StateID, f.ID)
	return err
}

func (m FlowModel) GetCurrentState(sessionID int64) (*Flow, error) {
	query := `SELECT id, session_id, state_id, comment, created_at 
			  FROM flows WHERE session_id = $1
			  ORDER BY created_at DESC LIMIT 1`

	ctx, cancel := context.WithTimeout(context.Background(), THREE_SECONDS)

	defer cancel()

	var f Flow
	var comment sql.NullString
	err := m.DB.QueryRowContext(ctx, query, sessionID).Scan(&f.ID, &f.SessionID, &f.StateID, &comment, &f.CreatedAt)

	if err != nil {
		return nil, err
	}

	if comment.Valid {
		f.Comment = comment.String
	}

	return &f, nil
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

func (m FlowModel) GetFullHistory(filters FlowFilters) ([]*Flow, Metadata, error) {
	query := fmt.Sprintf(`SELECT COUNT(*) OVER(), 
							s.public_id, f.id, f.session_id, f.state_id, ss.name, f.comment, f.created_at 
						FROM flows f
						INNER JOIN sessions s ON s.id = f.session_id
						INNER JOIN session_states ss ON f.state_id = ss.id
			  			WHERE 1=1 
			  			AND f.session_id=$1
			  			ORDER BY f.%s %s, f.id ASC
			  			LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	args := []any{filters.SessionID, filters.limit(), filters.offset()}

	ctx, cancel := context.WithTimeout(context.Background(), THREE_SECONDS)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	var flows []*Flow
	for rows.Next() {

		var f Flow
		var comment sql.NullString

		err = rows.Scan(&totalRecords, &f.SessionPublicID, &f.ID, &f.SessionID, &f.StateID, &f.StateName, &comment, &f.CreatedAt)
		if err != nil {
			return nil, Metadata{}, err
		}

		if comment.Valid {
			f.Comment = comment.String
		}

		flows = append(flows, &f)
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return flows, metadata, nil
}
