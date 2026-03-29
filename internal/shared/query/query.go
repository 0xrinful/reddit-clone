package query

import (
	"strconv"
	"strings"
)

type queryType int

const (
	selectQuery queryType = iota
	updateQuery
)

type Query struct {
	kind queryType

	prefix string
	table  string
	cols   string
	order  string
	group  string
	limit  int

	joins  strings.Builder
	sets   strings.Builder
	wheres strings.Builder
	args   []any
}

func New() *Query {
	return &Query{
		args: make([]any, 0, 8),
	}
}

func (q *Query) bindArgs(b *strings.Builder, clause string, values ...any) {
	argIdx := 0
	start := 0
	for i := 0; i < len(clause); i++ {
		if clause[i] == '?' {
			q.args = append(q.args, values[argIdx])
			argIdx++
			b.WriteString(clause[start:i])
			b.WriteByte('$')
			b.WriteString(strconv.Itoa(len(q.args)))
			start = i + 1
		}
	}
	b.WriteString(clause[start:])
}

func (q *Query) Update(table string) {
	q.sets.Grow(64)
	q.wheres.Grow(64)
	q.kind = updateQuery
	q.table = table
}

func (q *Query) Select(cols string, table string) {
	q.wheres.Grow(64)
	q.joins.Grow(64)
	q.kind = selectQuery
	q.table = table
	q.cols = cols
}

func (q *Query) Prefix(sql string, values ...any) {
	var b strings.Builder
	b.Grow(len(sql))
	q.bindArgs(&b, sql, values...)
	q.prefix = b.String()
}

func (q *Query) Set(clause string, values ...any) {
	if q.sets.Len() > 0 {
		q.sets.WriteString(", ")
	}
	q.bindArgs(&q.sets, clause, values...)
}

func (q *Query) Where(clause string, values ...any) {
	if q.wheres.Len() == 0 {
		q.wheres.WriteString(" WHERE ")
	} else {
		q.wheres.WriteString(" AND ")
	}
	q.bindArgs(&q.wheres, clause, values...)
}

func (q *Query) Order(clause string) {
	q.order = "ORDER BY " + clause
}

func (q *Query) Group(clause string) {
	q.group = "GROUP BY " + clause
}

func (q *Query) Join(clause string) {
	q.joins.WriteString(" INNER JOIN ")
	q.joins.WriteString(clause)
}

func (q *Query) LeftJoin(clause string) {
	q.joins.WriteString(" LEFT JOIN ")
	q.joins.WriteString(clause)
}

func (q *Query) RightJoin(clause string) {
	q.joins.WriteString(" RIGHT JOIN ")
	q.joins.WriteString(clause)
}

func (q *Query) Limit(limit int) {
	q.args = append(q.args, limit)
	q.limit = len(q.args)
}

func (q *Query) buildUpdate() (string, []any) {
	var b strings.Builder
	b.Grow(q.sets.Len() + q.wheres.Len() + 32)

	b.WriteString("UPDATE ")
	b.WriteString(q.table)
	b.WriteString(" SET ")

	b.WriteString(q.sets.String())
	b.WriteString(q.wheres.String())

	return b.String(), q.args
}

func (q *Query) buildSelect() (string, []any) {
	var b strings.Builder
	b.Grow(len(q.prefix) + q.joins.Len() + q.wheres.Len() + len(q.group) + len(q.order) + 64)

	if q.prefix != "" {
		b.WriteString(q.prefix)
		b.WriteByte(' ')
	}

	b.WriteString("SELECT ")
	b.WriteString(q.cols)
	b.WriteString(" FROM ")
	b.WriteString(q.table)

	b.WriteString(q.joins.String())

	b.WriteString(q.wheres.String())

	if q.group != "" {
		b.WriteByte(' ')
		b.WriteString(q.group)
	}

	if q.order != "" {
		b.WriteByte(' ')
		b.WriteString(q.order)
	}

	if q.limit != 0 {
		b.WriteString(" LIMIT $")
		b.WriteString(strconv.Itoa(q.limit))
	}

	return b.String(), q.args
}

func (q *Query) ToSql() (string, []any) {
	switch q.kind {
	case updateQuery:
		return q.buildUpdate()
	case selectQuery:
		return q.buildSelect()
	}
	panic("unsupported query type")
}
