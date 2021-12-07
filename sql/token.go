package sql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/leg100/otf"
)

var (
	_ otf.TokenStore = (*TokenDB)(nil)
)

type TokenDB struct {
	*sqlx.DB
}

func NewTokenDB(db *sqlx.DB) *TokenDB {
	return &TokenDB{
		DB: db,
	}
}

func (db TokenDB) Create(ctx context.Context, token otf.Token) error {
	sql, args, err := psql.
		Insert("tokens").
		Columns("token_id", "created_at", "updated_at", "description", "hash").
		Values(token.ID, token.CreatedAt, token.UpdatedAt, token.Description, token.Hash).
		ToSql()
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (db TokenDB) List(ctx context.Context) ([]*otf.Token, error) {
	selectBuilder := psql.Select("*").From("tokens").OrderBy("created_at DESC")

	sql, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var tokens []*otf.Token
	if err := db.SelectContext(ctx, &tokens, sql, args...); err != nil {
		return nil, fmt.Errorf("unable to scan tokens from DB: %w", err)
	}

	return tokens, nil
}

func (db TokenDB) Get(ctx context.Context, id string) (*otf.Token, error) {
	selectBuilder := psql.Select("*").From("tokens").Where("id = $1", id)

	sql, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("generating sql: %w", err)
	}

	var token otf.Token
	if err := db.GetContext(ctx, &token, sql, args...); err != nil {
		return nil, databaseError(err)
	}

	return &token, nil
}

// Delete deletes a token from the DB
func (db TokenDB) Delete(ctx context.Context, id string) error {
	var deleted string
	if err := db.GetContext(ctx, &deleted, "DELETE FROM tokens WHERE token_id = $1 RETURNING token_id", id); err != nil {
		return fmt.Errorf("unable to delete token: %w", err)
	}

	return nil
}
