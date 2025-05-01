package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail    = errors.New("a user with that email already exists")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
)

type UserStore struct {
	db *sql.DB
}

type User struct {
	ID          int64    `json:"id"`
	Email       string   `json:"email"`
	Username    string   `json:"username"`
	DisplayName *string  `json:"display_name"`
	Password    password `json:"-"`
	IsActive    bool     `json:"is_active"`
	CreatedAt   string   `json:"created_at"`
}

type password struct {
	text *string
	hash []byte
}

func (s UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		INSERT INTO users (username, password, email, display_name)
   		VALUES ($1, $2, $3, $4)
    	RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password.hash,
		user.Email,
		user.DisplayName,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}
	return nil
}

func (s *UserStore) GetByID(ctx context.Context, userID int64) (*User, error) {
	query := `
	SELECT id, email, username, display_name, password, created_at
	FROM users
	WHERE id = $1 AND is_active = true
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user User
	err := s.db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.DisplayName,
		&user.Password,
		&user.CreatedAt,
	)

	switch err {
	case nil:
		return &user, nil
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash
	return nil
}

func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}

		if err := s.createUserInvitation(ctx, tx, token, invitationExp, user.ID); err != nil {
			return err
		}
		return nil
	})

}

func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userID int64) error {
	query := `INSERT INTO user_invitations (token, user_id, expiry) VALUES ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(exp))

	if err != nil {
		return err
	}

	return nil
}
