// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: sessions.sql

package sqlc

import (
	"context"

	"github.com/RobinThrift/belt/internal/storage/database/sqlite/types"
)

const createSession = `-- name: CreateSession :exec
INSERT INTO sessions(
    token,
    data,
    expires_at
) VALUES (?, ?, ?)
ON CONFLICT (token)
DO UPDATE SET
    data       = excluded.data,
    expires_at = excluded.expires_at
`

type CreateSessionParams struct {
	Token     string
	Data      []byte
	ExpiresAt types.SQLiteDatetime
}

func (q *Queries) CreateSession(ctx context.Context, db DBTX, arg CreateSessionParams) error {
	_, err := db.ExecContext(ctx, createSession, arg.Token, arg.Data, arg.ExpiresAt)
	return err
}

const deleteExpired = `-- name: DeleteExpired :exec
DELETE FROM sessions WHERE datetime(expires_at) > CURRENT_TIMESTAMP
`

func (q *Queries) DeleteExpired(ctx context.Context, db DBTX) error {
	_, err := db.ExecContext(ctx, deleteExpired)
	return err
}

const deleteSession = `-- name: DeleteSession :exec
DELETE FROM sessions WHERE token = ?
`

func (q *Queries) DeleteSession(ctx context.Context, db DBTX, token string) error {
	_, err := db.ExecContext(ctx, deleteSession, token)
	return err
}

const getSession = `-- name: GetSession :one
SELECT
    sessions.id, sessions.token, sessions.data, sessions.expires_at
FROM sessions
WHERE token = ?
LIMIT 1
`

func (q *Queries) GetSession(ctx context.Context, db DBTX, token string) (Session, error) {
	row := db.QueryRowContext(ctx, getSession, token)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.Token,
		&i.Data,
		&i.ExpiresAt,
	)
	return i, err
}
