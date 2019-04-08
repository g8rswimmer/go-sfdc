package soql

import (
	"fmt"
	"strings"
	"time"
)

type Builder struct {
	fieldList  []string
	objectType string
	subQuery   []fmt.Stringer
	where      fmt.Stringer
	order      fmt.Stringer
	limit      int
	offset     int
}

func (b *Builder) FieldList(fields ...string) {
	b.fieldList = append(b.fieldList, fields...)
}
func (b *Builder) ObjectType(object string) {
	b.objectType = object
}
func (b *Builder) SubQuery(query fmt.Stringer) {
	b.subQuery = append(b.subQuery, query)
}
func (b *Builder) Where(where fmt.Stringer) {
	b.where = where
}
func (b *Builder) OrderBy(order fmt.Stringer) {
	b.order = order
}
func (b *Builder) Limit(limit int) {
	b.limit = limit
}
func (b *Builder) Offset(offset int) {
	b.offset = offset
}
func (b *Builder) String() string {
	soql := "SELECT " + strings.Join(b.fieldList, ",")
	if b.subQuery != nil {
		for _, query := range b.subQuery {
			soql += fmt.Sprintf(",(%s)", query.String())
		}
	}
	soql += " FROM " + b.objectType
	if b.where != nil {
		soql += " WHERE " + b.where.String()
	}
	if b.order != nil {
		soql += " ORDER BY " + b.order.String()
	}
	if b.limit > 0 {
		soql += fmt.Sprintf(" LIMIT %d", b.limit)
	}
	if b.offset > 0 {
		soql += fmt.Sprintf(" OFFSET %d", b.offset)
	}
	return soql
}

type WhereClause struct {
	clause string
}

func WhereLike(field string, value string) *WhereClause {
	return &WhereClause{
		clause: fmt.Sprintf("%s LIKE '%s'", field, value),
	}
}
func WhereGreaterThan(field string, value interface{}, equals bool) *WhereClause {
	var v string
	if value != nil {
		switch value.(type) {
		case string:
			v = fmt.Sprintf("'%s'", value.(string))
		case time.Time:
			date := value.(time.Time)
			v = date.Format(time.RFC3339)
		default:
			v = fmt.Sprintf("%v", value)
		}
	}

	operator := ">"
	if equals {
		operator += "="
	}

	return &WhereClause{
		clause: fmt.Sprintf("%s %s %s", field, operator, v),
	}
}
func WhereLessThan(field string, value interface{}, equals bool) *WhereClause {
	var v string
	if value != nil {
		switch value.(type) {
		case string:
			v = fmt.Sprintf("'%s'", value.(string))
		case time.Time:
			date := value.(time.Time)
			v = date.Format(time.RFC3339)
		default:
			v = fmt.Sprintf("%v", value)
		}
	}

	operator := "<"
	if equals {
		operator += "="
	}

	return &WhereClause{
		clause: fmt.Sprintf("%s %s %s", field, operator, v),
	}
}
func WhereEquals(field string, value interface{}) *WhereClause {
	var v string
	if value != nil {
		switch value.(type) {
		case string:
			v = fmt.Sprintf("'%s'", value.(string))
		case time.Time:
			date := value.(time.Time)
			v = date.Format(time.RFC3339)
		default:
			v = fmt.Sprintf("%v", value)
		}
	} else {
		v = "null"
	}

	return &WhereClause{
		clause: fmt.Sprintf("%s = %s", field, v),
	}
}
func WhereNotEquals(field string, value interface{}) *WhereClause {
	var v string
	if value != nil {
		switch value.(type) {
		case string:
			v = fmt.Sprintf("'%s'", value.(string))
		case time.Time:
			date := value.(time.Time)
			v = date.Format(time.RFC3339)
		default:
			v = fmt.Sprintf("%v", value)
		}
	} else {
		v = "null"
	}

	return &WhereClause{
		clause: fmt.Sprintf("%s != %s", field, v),
	}
}
func WhereIn(field string, values []interface{}) *WhereClause {
	set := make([]string, len(values))
	for idx, value := range values {
		switch value.(type) {
		case string:
			set[idx] = fmt.Sprintf("'%s'", value.(string))
		case time.Time:
			date := value.(time.Time)
			set[idx] = date.Format(time.RFC3339)
		default:
			set[idx] = fmt.Sprintf("%v", value)
		}
	}

	return &WhereClause{
		clause: fmt.Sprintf("%s IN (%s)", field, strings.Join(set, ",")),
	}
}

func WhereNotIn(field string, values []interface{}) *WhereClause {
	set := make([]string, len(values))
	for idx, value := range values {
		switch value.(type) {
		case string:
			set[idx] = fmt.Sprintf("'%s'", value.(string))
		default:
			set[idx] = fmt.Sprintf("%v", value)
		}
	}

	return &WhereClause{
		clause: fmt.Sprintf("%s NOT IN (%s)", field, strings.Join(set, ",")),
	}
}
