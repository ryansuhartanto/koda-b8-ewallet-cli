package db

import (
	"context"
	"iter"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Key interface {
	Bindings() iter.Seq2[string, any]
}

type ID int64

func (id ID) Bindings() iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		if !yield("id", id) {
			return
		}
	}
}

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

type RepoKeyed[T any, K Key] Repo[T]

func NewRepoKeyeded[T any, K Key](db Querier, table string) RepoKeyed[T, K] {
	return RepoKeyed[T, K](NewRepo[T](db, table))
}

func conjunction(preds iter.Seq2[string, any]) (string, []any) {
	var sb strings.Builder
	var args []any

	for col, val := range preds {
		if len(args) > 0 {
			sb.WriteString(" AND ")
		}
		sb.WriteString(col)
		sb.WriteString(" = $")
		sb.WriteString(strconv.Itoa(len(args) + 1))

		args = append(args, val)
	}

	return sb.String(), args
}

func (r RepoKeyed[T, K]) Get(ctx context.Context, key K) (T, error) {
	var zero T
	predicates, values := conjunction(key.Bindings())

	rows, err := r.db.Query(
		ctx,
		`SELECT * FROM `+r.table+` WHERE `+predicates,
		values...,
	)
	if err != nil {
		return zero, err
	}
	return pgx.CollectOneRow(rows, pgx.RowToStructByName[T])
}

func (r RepoKeyed[T, K]) Delete(ctx context.Context, key K) error {
	predicates, values := conjunction(key.Bindings())

	_, err := r.db.Exec(
		ctx,
		`DELETE FROM `+r.table+` WHERE `+predicates,
		values...,
	)
	return err
}

type RepoSoftDelete[T any, K Key] RepoKeyed[T, K]

func NewRepoSoftDelete[T any, K Key](db Querier, table string) RepoSoftDelete[T, K] {
	return RepoSoftDelete[T, K](NewRepoKeyeded[T, K](db, table))
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
	predicates, values := conjunction(key.Bindings())

	rows, err := r.db.Query(
		ctx,
		`SELECT * FROM `+r.table+
			` WHERE `+predicates+` AND "deleted_at" IS NULL`,
		values...,
	)
	if err != nil {
		return zero, err
	}
	return pgx.CollectOneRow(rows, pgx.RowToStructByName[T])
}

func (r RepoSoftDelete[T, K]) Delete(ctx context.Context, key K) error {
	predicates, values := conjunction(key.Bindings())

	_, err := r.db.Exec(
		ctx,
		`UPDATE `+r.table+
			` SET "deleted_at" = CURRENT_TIMESTAMP`+
			` WHERE `+predicates+` AND "deleted_at" IS NULL`,
		values...,
	)
	return err
}
