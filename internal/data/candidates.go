package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Candidate struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ImageURL  string    `json:"imageUrl"`
	UserID    int64     `json:"userId"`
	SessionID int64     `json:"sessionId"`
	CreatedAt time.Time `json:"createdAt"`
}

type CandidateModel struct {
	DB *sql.DB
}

func (m CandidateModel) Insert(c *Candidate) error {
	query := `INSERT INTO candidates (name, image_url, user_id, session_id, created_at)
				 VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`

	args := []any{c.Name, c.ImageURL, c.UserID, c.SessionID, c.CreatedAt}

	ctx, cancel := context.WithTimeout(context.Background(), THREE_SECONDS)

	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&c.ID, &c.CreatedAt)

}

func (m CandidateModel) ListFiltering(filters *CandidateFilters) ([]*Candidate, Metadata, error) {
	query := fmt.Sprintf(`SELECT COUNT(*) OVER(), id, name, image_url, user_id, created_at 
							FROM candidates
							WHERE 1=1
							AND session_id = $1
							AND to_tsvector('simple', name) @@ plainto_tsquery('simple', $2) OR $2 = '' OR $2 IS NULL
							ORDER BY %s %s, id ASC
							LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	args := []any{filters.SessionID, filters.Name, filters.limit(), filters.offset()}

	ctx, cancel := context.WithTimeout(context.Background(), THREE_SECONDS)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	candidates := make([]*Candidate, 0)
	totalRecords := 0

	defer rows.Close()
	for rows.Next() {
		var c Candidate
		var imgURL sql.NullString

		err = rows.Scan(&totalRecords, &c.ID, &c.Name, &imgURL, &c.UserID, &c.CreatedAt)
		if err != nil {
			return nil, Metadata{}, err
		}

		if imgURL.Valid {
			c.ImageURL = imgURL.String
		}

		candidates = append(candidates, &c)
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return candidates, metadata, nil
}

func (m CandidateModel) Get(id int64) (*Candidate, error) {
	query := `SELECT id, name, user_id, created_at FROM candidates WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), THREE_SECONDS)
	defer cancel()

	var c Candidate
	var imgURL sql.NullString

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.Name, &c.UserID, &c.CreatedAt)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	if imgURL.Valid {
		c.ImageURL = imgURL.String
	}

	return &c, nil
}
