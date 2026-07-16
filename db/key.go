package db

type Key interface {
	Columns() []string
	Values() []any
}

type ID int64

func (i ID) Columns() []string { return []string{"id"} }
func (i ID) Values() []any     { return []any{int64(i)} }
