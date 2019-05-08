package soql

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// QueryInput is used to provide SOQL inputs.
//
// ObjectType is the Salesforce Object, like Account
//
// FieldList is the Salesforce Object's fields to query
//
// SubQuery is the inner query
//
// Where is the SOQL where cause
//
// Order is the SOQL ordering
//
// Limit is the SOQL record limit
//
// Offset is the SOQL record offset
type QueryInput struct {
	FieldList  []string
	ObjectType string
	SubQuery   []Querier
	Where      WhereClauser
	Order      Orderer
	Limit      int
	Offset     int
}

// Builder is the struture used to build a SOQL query.
type Builder struct {
	fieldList  []string
	objectType string
	subQuery   []Querier
	where      WhereClauser
	order      Orderer
	limit      int
	offset     int
}

// Querier is the interface to return the SOQL query.
//
// Query returns the SOQL query.
type Querier interface {
	Query() (string, error)
}

// NewBuilder creates a new builder.  If the object is an
// empty string, then an error is returned.
func NewBuilder(input QueryInput) (*Builder, error) {
	if input.ObjectType == "" {
		return nil, errors.New("builder: object type can not be an empty string")
	}
	if len(input.FieldList) == 0 {
		return nil, errors.New("builder: field list can not be empty")
	}

	return &Builder{
		objectType: input.ObjectType,
		fieldList:  input.FieldList,
		subQuery:   input.SubQuery,
		where:      input.Where,
		order:      input.Order,
		limit:      input.Limit,
		offset:     input.Offset,
	}, nil
}

// Query will return the SOQL query.  If the builder has an empty string or
// the field list is zero, an error is returned.
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
		order, err := b.order.Order()
		if err == nil {
			soql += " " + order
		} else {
			return "", err
		}
	}
	if b.limit > 0 {
		soql += fmt.Sprintf(" LIMIT %d", b.limit)
	}
	if b.offset > 0 {
		soql += fmt.Sprintf(" OFFSET %d", b.offset)
	}
	return soql, nil
}

// WhereClause is the structure that will contain a SOQL where clause.
type WhereClause struct {
	expression string
}

// WhereExpression is an interface to return the where cause's expression.
type WhereExpression interface {
	Expression() string
}

// WhereClauser is an interface to return the where cause.
type WhereClauser interface {
	Clause() string
}

// WhereLike will form the LIKE expression.
func WhereLike(field string, value string) (*WhereClause, error) {
	if field == "" {
		return nil, errors.New("soql where: field can not be empty")
	}
	if value == "" {
		return nil, errors.New("soql where: value can not be empty")
	}
	return &WhereClause{
		expression: fmt.Sprintf("%s LIKE '%s'", field, value),
	}, nil
}

// WhereGreaterThan will form the greater or equal than expression.  If the value is a
// string or boolean, an error is returned.
func WhereGreaterThan(field string, value interface{}, equals bool) (*WhereClause, error) {
	if field == "" {
		return nil, errors.New("soql where: field can not be empty")
	}
	if value == nil {
		return nil, errors.New("soql where: value can not be nil")
	}
	var v string
	switch value.(type) {
	case string, bool:
		return nil, errors.New("where greater than: value can not be a string or bool")
	case time.Time:
		date := value.(time.Time)
		v = date.Format(time.RFC3339)
	default:
		v = fmt.Sprintf("%v", value)
	}

	operator := ">"
	if equals {
		operator += "="
	}

	return &WhereClause{
		expression: fmt.Sprintf("%s %s %s", field, operator, v),
	}, nil
}

// WhereLessThan will form the less or equal than expression.  If the value is a
// string or boolean, an error is returned.
func WhereLessThan(field string, value interface{}, equals bool) (*WhereClause, error) {
	if field == "" {
		return nil, errors.New("soql where: field can not be empty")
	}
	if value == nil {
		return nil, errors.New("soql where: value can not be nil")
	}
	var v string
	switch value.(type) {
	case string, bool:
		return nil, errors.New("where less than: value can not be a string")
	case time.Time:
		date := value.(time.Time)
		v = date.Format(time.RFC3339)
	default:
		v = fmt.Sprintf("%v", value)
	}

	operator := "<"
	if equals {
		operator += "="
	}

	return &WhereClause{
		expression: fmt.Sprintf("%s %s %s", field, operator, v),
	}, nil
}

// WhereEquals forms the equals where expression.
func WhereEquals(field string, value interface{}) (*WhereClause, error) {
	if field == "" {
		return nil, errors.New("soql where: field can not be empty")
	}
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

// WhereNotEquals forms the not equals where expression.
func WhereNotEquals(field string, value interface{}) (*WhereClause, error) {
	if field == "" {
		return nil, errors.New("soql where: field can not be empty")
	}
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

// WhereIn forms the field in a set expression.
func WhereIn(field string, values []interface{}) (*WhereClause, error) {
	if field == "" {
		return nil, errors.New("soql where: field can not be empty")
	}
	if values == nil {
		return nil, errors.New("soql where: value array can not be nil")
	}
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

// WhereNotIn forms the field is not in a set expression.
func WhereNotIn(field string, values []interface{}) (*WhereClause, error) {
	if field == "" {
		return nil, errors.New("soql where: field can not be empty")
	}
	if values == nil {
		return nil, errors.New("soql where: value array can not be nil")
	}
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

// Clause returns the where cluase.
func (wc *WhereClause) Clause() string {
	return fmt.Sprintf("WHERE %s", wc.expression)
}

// Group will form a grouping around the expression.
func (wc *WhereClause) Group() {
	wc.expression = fmt.Sprintf("(%s)", wc.expression)
}

// And will logical AND the expressions.
func (wc *WhereClause) And(where WhereExpression) {
	wc.expression = fmt.Sprintf("%s AND %s", wc.expression, where.Expression())
}

// Or will logical OR the expressions.
func (wc *WhereClause) Or(where WhereExpression) {
	wc.expression = fmt.Sprintf("%s OR %s", wc.expression, where.Expression())
}

// Expression will return the where expression.
func (wc *WhereClause) Expression() string {
	return wc.expression
}

// OrderResult is the type of ordering of the query result.
type OrderResult string

const (
	// OrderAsc will place the results in ascending order.
	OrderAsc OrderResult = "ASC"
	// OrderDesc will place the results in descending order.
	OrderDesc OrderResult = "DESC"
)

// OrderNulls is where the null values are placed in the ordering.
type OrderNulls string

const (
	// OrderNullsLast places the null values at the end of the ordering.
	OrderNullsLast OrderNulls = "NULLS LAST"
	// OrderNullsFirst places the null values at the start of the ordering.
	OrderNullsFirst OrderNulls = "NULLS FIRST"
)

// OrderBy is the ordering structure of the SOQL query.
type OrderBy struct {
	fieldOrder []string
	result     OrderResult
	nulls      OrderNulls
}

// Orderer is the interface for returning the SOQL ordering.
type Orderer interface {
	Order() (string, error)
}

// NewOrderBy creates an OrderBy structure.  If the order results is not ASC or DESC, an
// error will be returned.
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

// FieldOrder is a list of fields in the ordering.
func (o *OrderBy) FieldOrder(fields ...string) {
	o.fieldOrder = append(o.fieldOrder, fields...)
}

// NullOrdering sets the ordering, first or last, of the null values.
func (o *OrderBy) NullOrdering(nulls OrderNulls) error {
	switch nulls {
	case OrderNullsLast, OrderNullsFirst:
	default:
		return fmt.Errorf("order by: %s is not a valid null ordering type", string(nulls))
	}
	o.nulls = nulls
	return nil
}

// Order returns the order by SOQL string.
func (o *OrderBy) Order() (string, error) {
	switch o.result {
	case OrderAsc, OrderDesc:
	default:
		return "", fmt.Errorf("order by: %s is not a valid result ordering type", string(o.result))
	}

	orderBy := "ORDER BY " + strings.Join(o.fieldOrder, ",") + " " + string(o.result)
	if o.nulls != "" {
		orderBy += " " + string(o.nulls)
	}
	return orderBy, nil
}
