package data

import (
	"context"
	"database/sql"
	"time"
)

type Flow struct {
	ID        int64     `json:"id"`
	SessionID int64     `json:"sessionId"`
	StateID   int64     `json:"stateId"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
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
	err := m.DB.QueryRowContext(ctx, query, sessionID).Scan(&f.ID, &f.SessionID, &f.StateID, &f.Comment, &f.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &f, nil
}

func (m FlowModel) GetFullHistory(sessionID int64) ([]*Flow, error) {
	query := `SELECT id, session_id, state_id, comment, created_at FROM flows
				WHERE session_id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), THREE_SECONDS)

	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, sessionID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var flows []*Flow
	for rows.Next() {

		var f Flow

		err = rows.Scan(f.ID, f.SessionID, f.StateID, f.Comment, f.CreatedAt)
		if err != nil {
			return nil, err
		}

		flows = append(flows, &f)
	}

	return flows, nil
}
