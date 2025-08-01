package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	PublicID  string    `json:"publicId"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Get(publicID string) (*User, error) {
	query := `SELECT id, public_id, name, email, password, created_at 
			  FROM users 
			  WHERE public_id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), THREE_SECONDS)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, publicID)

	var user User

	err := row.Scan(&user.ID, &user.PublicID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		} else {
			return nil, err
		}
	}

	return &user, nil
}
