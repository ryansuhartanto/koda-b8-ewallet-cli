package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Repo[T any] struct {
	db    Querier
	table string
}

func NewRepo[T any](db Querier, table string) Repo[T] {
	return Repo[T]{db: db, table: table}
}

func (r Repo[T]) List(ctx context.Context) ([]T, error) {
	rows, err := r.db.Query(ctx, `SELECT * FROM `+r.table)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[T])
}

type RepoKeyed[T any, K Key] struct {
	Repo[T]
}

func NewRepoKeyed[T any, K Key](db Querier, table string) RepoKeyed[T, K] {
	return RepoKeyed[T, K]{NewRepo[T](db, table)}
}

func wherePredicates(cols []string) string {
	preds := make([]string, len(cols))
	for i, c := range cols {
		preds[i] = fmt.Sprintf(`%q = $%d`, c, i+1)
	}
	return strings.Join(preds, " AND ")
}

func (r RepoKeyed[T, K]) Get(ctx context.Context, key K) (T, error) {
	var zero T

	rows, err := r.db.Query(
		ctx,
		`SELECT * FROM `+r.table+` WHERE `+wherePredicates(key.Columns()),
		key.Values()...,
	)
	if err != nil {
		return zero, err
	}
	return pgx.CollectOneRow(rows, pgx.RowToStructByName[T])
}

func (r RepoKeyed[T, K]) Delete(ctx context.Context, key K) error {
	_, err := r.db.Exec(
		ctx,
		`DELETE FROM `+r.table+` WHERE `+wherePredicates(key.Columns()),
		key.Values()...,
	)
	return err
}

type RepoSoftDelete[T any, K Key] struct {
	RepoKeyed[T, K]
}

func NewRepoSoftDelete[T any, K Key](db Querier, table string) RepoSoftDelete[T, K] {
	return RepoSoftDelete[T, K]{NewRepoKeyed[T, K](db, table)}
}

func (r RepoSoftDelete[T, K]) List(ctx context.Context) ([]T, error) {
	rows, err := r.db.Query(ctx, `SELECT * FROM `+r.table+` WHERE "deleted_at" IS NULL`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[T])
}

func (r RepoSoftDelete[T, K]) Get(ctx context.Context, key K) (T, error) {
	var zero T

	rows, err := r.db.Query(
		ctx,
		`SELECT * FROM `+r.table+
			` WHERE `+wherePredicates(key.Columns())+` AND "deleted_at" IS NULL`,
		key.Values()...,
	)
	if err != nil {
		return zero, err
	}
	return pgx.CollectOneRow(rows, pgx.RowToStructByName[T])
}

func (r RepoSoftDelete[T, K]) Delete(ctx context.Context, key K) error {
	_, err := r.db.Exec(
		ctx,
		`UPDATE `+r.table+
			` SET "deleted_at" = CURRENT_TIMESTAMP`+
			` WHERE `+wherePredicates(key.Columns())+` AND "deleted_at" IS NULL`,
		key.Values()...,
	)
	return err
}
