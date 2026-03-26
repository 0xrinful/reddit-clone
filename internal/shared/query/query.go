package query

import (
	"fmt"
	"strings"
)

type queryType int

const (
	updateQuery queryType = iota
)

type Raw string

type Query struct {
	kind queryType

	table string

	sets   []string
	wheres []string
	args   []any
}

func (q *Query) Update(table string) {
	q.kind = updateQuery
	q.table = table
}

func (q *Query) Set(col string, value any) {
	switch v := value.(type) {
	case Raw:
		q.sets = append(q.sets, fmt.Sprintf("%s = %s", col, v))
	default:
		q.args = append(q.args, value)
		q.sets = append(q.sets, fmt.Sprintf("%s = $%d", col, len(q.args)))
	}
}

func (q *Query) Where(sql string, values ...any) {
	var b strings.Builder
	argIdx := 0
	for i := 0; i < len(sql); i++ {
		if sql[i] != '?' {
			b.WriteByte(sql[i])
			continue
		}
		q.args = append(q.args, values[argIdx])
		fmt.Fprintf(&b, "$%d", len(q.args))
		argIdx++
	}
	q.wheres = append(q.wheres, b.String())
}

func (q *Query) buildUpdate() (string, []any) {
	sql := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s",
		q.table,
		strings.Join(q.sets, ", "),
		strings.Join(q.wheres, " AND "),
	)
	return sql, q.args
}

func (q *Query) ToSql() (string, []any) {
	switch q.kind {
	case updateQuery:
		return q.buildUpdate()
	}
	panic("unsupported query type")
}
