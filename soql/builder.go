package soql

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Builder struct {
	fieldList  []string
	objectType string
	subQuery   []Querier
	where      WhereClauser
	order      Orderer
	limit      int
	offset     int
}

type Querier interface {
	Query() (string, error)
}

func NewBuilder(object string) (*Builder, error) {
	if object == "" {
		return nil, errors.New("builder: object type can not be an empty string")
	}
	return &Builder{
		objectType: object,
	}, nil
}
func (b *Builder) FieldList(fields ...string) {
	if b == nil {
		panic("builder can not be nil")
	}
	b.fieldList = append(b.fieldList, fields...)
}
func (b *Builder) SubQuery(query Querier) {
	if b == nil {
		panic("builder can not be nil")
	}
	b.subQuery = append(b.subQuery, query)
}
func (b *Builder) Where(where WhereClauser) {
	if b == nil {
		panic("builder can not be nil")
	}
	b.where = where
}
func (b *Builder) OrderBy(order Orderer) {
	if b == nil {
		panic("builder can not be nil")
	}
	b.order = order
}
func (b *Builder) Limit(limit int) {
	if b == nil {
		panic("builder can not be nil")
	}
	b.limit = limit
}
func (b *Builder) Offset(offset int) {
	if b == nil {
		panic("builder can not be nil")
	}
	b.offset = offset
}
func (b *Builder) Query() (string, error) {
	if b.objectType == "" {
		return "", errors.New("builder: object type can not be an empty string")
	}
	if len(b.fieldList) == 0 {
		return "", errors.New("builder: field list must be have fields present")
	}

	soql := "SELECT " + strings.Join(b.fieldList, ",")
	if b.subQuery != nil {
		for _, query := range b.subQuery {
			var sub string
			var err error
			if sub, err = query.Query(); err == nil {
				soql += fmt.Sprintf(",(%s)", sub)
			} else {
				return "", err
			}
		}
	}
	soql += " FROM " + b.objectType
	if b.where != nil {
		soql += " " + b.where.Clause()
	}
	if b.order != nil {
		soql += " " + b.order.Order()
	}
	if b.limit > 0 {
		soql += fmt.Sprintf(" LIMIT %d", b.limit)
	}
	if b.offset > 0 {
		soql += fmt.Sprintf(" OFFSET %d", b.offset)
	}
	return soql, nil
}

type WhereClause struct {
	expression string
}

type WhereExpression interface {
	Expression() string
}
type WhereClauser interface {
	Clause() string
}

func WhereLike(field string, value string) (*WhereClause, error) {
	return &WhereClause{
		expression: fmt.Sprintf("%s LIKE '%s'", field, value),
	}, nil
}
func WhereGreaterThan(field string, value interface{}, equals bool) (*WhereClause, error) {
	var v string
	if value != nil {
		switch value.(type) {
		case string, bool:
			return nil, errors.New("where greater than: value can not be a string or bool")
		case time.Time:
			date := value.(time.Time)
			v = date.Format(time.RFC3339)
		default:
			v = fmt.Sprintf("%v", value)
		}
	} else {
		return nil, errors.New("where greater than: value can not be nil")
	}

	operator := ">"
	if equals {
		operator += "="
	}

	return &WhereClause{
		expression: fmt.Sprintf("%s %s %s", field, operator, v),
	}, nil
}
func WhereLessThan(field string, value interface{}, equals bool) (*WhereClause, error) {
	var v string
	if value != nil {
		switch value.(type) {
		case string, bool:
			return nil, errors.New("where less than: value can not be a string")
		case time.Time:
			date := value.(time.Time)
			v = date.Format(time.RFC3339)
		default:
			v = fmt.Sprintf("%v", value)
		}
	} else {
		return nil, errors.New("where less than: value can not be nil")
	}

	operator := "<"
	if equals {
		operator += "="
	}

	return &WhereClause{
		expression: fmt.Sprintf("%s %s %s", field, operator, v),
	}, nil
}
func WhereEquals(field string, value interface{}) (*WhereClause, error) {
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
		expression: fmt.Sprintf("%s = %s", field, v),
	}, nil
}
func WhereNotEquals(field string, value interface{}) (*WhereClause, error) {
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
		expression: fmt.Sprintf("%s != %s", field, v),
	}, nil
}
func WhereIn(field string, values []interface{}) (*WhereClause, error) {
	set := make([]string, len(values))
	for idx, value := range values {
		switch value.(type) {
		case string:
			set[idx] = fmt.Sprintf("'%s'", value.(string))
		case bool:
			return nil, errors.New("where in: boolean is not a value set value")
		case time.Time:
			date := value.(time.Time)
			set[idx] = date.Format(time.RFC3339)
		default:
			set[idx] = fmt.Sprintf("%v", value)
		}
	}

	return &WhereClause{
		expression: fmt.Sprintf("%s IN (%s)", field, strings.Join(set, ",")),
	}, nil
}

func WhereNotIn(field string, values []interface{}) (*WhereClause, error) {
	set := make([]string, len(values))
	for idx, value := range values {
		switch value.(type) {
		case string:
			set[idx] = fmt.Sprintf("'%s'", value.(string))
		case bool:
			return nil, errors.New("where not in: boolean is not a value set value")
		case time.Time:
			date := value.(time.Time)
			set[idx] = date.Format(time.RFC3339)
		default:
			set[idx] = fmt.Sprintf("%v", value)
		}
	}

	return &WhereClause{
		expression: fmt.Sprintf("%s NOT IN (%s)", field, strings.Join(set, ",")),
	}, nil
}

func (wc *WhereClause) Clause() string {
	if wc == nil {
		panic("WhereClause can not be nil")
	}

	return fmt.Sprintf("WHERE %s", wc.expression)
}
func (wc *WhereClause) Group() {
	if wc == nil {
		panic("WhereClause can not be nil")
	}
	wc.expression = fmt.Sprintf("(%s)", wc.expression)
}
func (wc *WhereClause) And(where WhereExpression) {
	if wc == nil {
		panic("WhereClause can not be nil")
	}
	wc.expression = fmt.Sprintf("%s AND %s", wc.expression, where.Expression())
}
func (wc *WhereClause) Or(where WhereExpression) {
	if wc == nil {
		panic("WhereClause can not be nil")
	}
	wc.expression = fmt.Sprintf("%s OR %s", wc.expression, where.Expression())
}
func (wc *WhereClause) Expression() string {
	if wc == nil {
		panic("WhereClause can not be nil")
	}
	return wc.expression
}

type OrderResult string

const (
	OrderAsc  OrderResult = "ASC"
	OrderDesc OrderResult = "DESC"
)

type OrderNulls string

const (
	OrderNullsLast  OrderNulls = "NULLS LAST"
	OrderNullsFirst OrderNulls = "NULLS FIRST"
)

type OrderBy struct {
	fieldOrder []string
	result     OrderResult
	nulls      OrderNulls
}

type Orderer interface {
	Order() string
}

func NewOrderBy(result OrderResult) (*OrderBy, error) {
	switch result {
	case OrderAsc, OrderDesc:
	default:
		return nil, fmt.Errorf("order by: %s is not a valid result ordering type", string(result))
	}

	return &OrderBy{
		result: result,
	}, nil

}
func (o *OrderBy) FieldOrder(fields ...string) {
	if o == nil {
		panic("OrderBy can not be nil")
	}

	o.fieldOrder = append(o.fieldOrder, fields...)
}

func (o *OrderBy) NullOrdering(nulls OrderNulls) error {
	if o == nil {
		panic("OrderBy can not be nil")
	}

	switch nulls {
	case OrderNullsLast, OrderNullsFirst:
	default:
		return fmt.Errorf("order by: %s is not a valid null ordering type", string(nulls))
	}
	o.nulls = nulls
	return nil
}

func (o *OrderBy) Order() string {
	if o == nil {
		panic("OrderBy can not be nil")
	}

	orderBy := "ORDER BY " + strings.Join(o.fieldOrder, ",") + " " + string(o.result)
	if o.nulls != "" {
		orderBy += " " + string(o.nulls)
	}
	return orderBy
}
