package soql

import (
	"reflect"
	"testing"
	"time"
)

func TestNewOrderBy(t *testing.T) {
	type args struct {
		result OrderResult
	}
	tests := []struct {
		name    string
		args    args
		want    *OrderBy
		wantErr bool
	}{
		{
			name: "No Ordering Result",
			args: args{
				result: OrderResult(""),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Asc Ordering Result",
			args: args{
				result: OrderAsc,
			},
			want: &OrderBy{
				result: OrderAsc,
			},
			wantErr: false,
		},
		{
			name: "Desc Ordering Result",
			args: args{
				result: OrderDesc,
			},
			want: &OrderBy{
				result: OrderDesc,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOrderBy(tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOrderBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOrderBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderBy_FieldOrder(t *testing.T) {
	type fields struct {
		fieldOrder []string
		result     OrderResult
		nulls      OrderNulls
	}
	type args struct {
		fields []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *OrderBy
	}{
		{
			name:   "Add Order Fields",
			fields: fields{},
			args: args{
				fields: []string{
					"Name",
					"Date",
					"Time",
				},
			},
			want: &OrderBy{
				fieldOrder: []string{
					"Name",
					"Date",
					"Time",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OrderBy{
				fieldOrder: tt.fields.fieldOrder,
				result:     tt.fields.result,
				nulls:      tt.fields.nulls,
			}
			o.FieldOrder(tt.args.fields...)
			if !reflect.DeepEqual(o, tt.want) {
				t.Errorf("OrderBy.FieldOrder() = %v, want %v", o, tt.want)
			}
		})
	}
}

func TestOrderBy_NullOrdering(t *testing.T) {
	type fields struct {
		fieldOrder []string
		result     OrderResult
		nulls      OrderNulls
	}
	type args struct {
		nulls OrderNulls
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *OrderBy
		wantErr bool
	}{
		{
			name:   "No Null Ordering",
			fields: fields{},
			args: args{
				nulls: OrderNulls(""),
			},
			want:    &OrderBy{},
			wantErr: true,
		},
		{
			name:   "Last Null Ordering",
			fields: fields{},
			args: args{
				nulls: OrderNullsLast,
			},
			want: &OrderBy{
				nulls: OrderNullsLast,
			},
			wantErr: false,
		},
		{
			name:   "First Null Ordering",
			fields: fields{},
			args: args{
				nulls: OrderNullsFirst,
			},
			want: &OrderBy{
				nulls: OrderNullsFirst,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OrderBy{
				fieldOrder: tt.fields.fieldOrder,
				result:     tt.fields.result,
				nulls:      tt.fields.nulls,
			}
			if err := o.NullOrdering(tt.args.nulls); (err != nil) != tt.wantErr {
				t.Errorf("OrderBy.NullOrdering() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(o, tt.want) {
				t.Errorf("OrderBy.NullOrdering() = %v, want %v", o, tt.want)
			}
		})
	}
}

func TestOrderBy_Order(t *testing.T) {
	type fields struct {
		fieldOrder []string
		result     OrderResult
		nulls      OrderNulls
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Error",
			fields: fields{
				fieldOrder: []string{
					"Name",
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Basic Ascending",
			fields: fields{
				fieldOrder: []string{
					"Name",
				},
				result: OrderAsc,
			},
			want:    "ORDER BY Name ASC",
			wantErr: false,
		},
		{
			name: "Basic Descending",
			fields: fields{
				fieldOrder: []string{
					"Name",
				},
				result: OrderDesc,
			},
			want:    "ORDER BY Name DESC",
			wantErr: false,
		},
		{
			name: "Multiple Ascending",
			fields: fields{
				fieldOrder: []string{
					"Name",
					"Date",
				},
				result: OrderAsc,
			},
			want:    "ORDER BY Name,Date ASC",
			wantErr: false,
		},
		{
			name: "Nulls Last",
			fields: fields{
				fieldOrder: []string{
					"Name",
				},
				result: OrderDesc,
				nulls:  OrderNullsLast,
			},
			want:    "ORDER BY Name DESC NULLS LAST",
			wantErr: false,
		},
		{
			name: "Nulls First",
			fields: fields{
				fieldOrder: []string{
					"Name",
				},
				result: OrderDesc,
				nulls:  OrderNullsFirst,
			},
			want:    "ORDER BY Name DESC NULLS FIRST",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OrderBy{
				fieldOrder: tt.fields.fieldOrder,
				result:     tt.fields.result,
				nulls:      tt.fields.nulls,
			}
			got, err := o.Order()
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderBy.NullOrdering() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("OrderBy.Order() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWhereLike(t *testing.T) {
	type args struct {
		field string
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    *WhereClause
		wantErr bool
	}{
		{
			name: "Like Where",
			args: args{
				field: "Name",
				value: "A%",
			},
			want: &WhereClause{
				expression: "Name LIKE 'A%'",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WhereLike(tt.args.field, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("WhereLike() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WhereLike() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWhereGreaterThan(t *testing.T) {
	type args struct {
		field  string
		value  interface{}
		equals bool
	}
	tests := []struct {
		name    string
		args    args
		want    *WhereClause
		wantErr bool
	}{
		{
			name: "value is string",
			args: args{
				field: "Name",
				value: "some string",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "value is boolean",
			args: args{
				field: "Name",
				value: false,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "value is date",
			args: args{
				field: "Name",
				value: time.Date(2019, time.April, 15, 0, 0, 0, 0, time.UTC),
			},
			want: &WhereClause{
				expression: "Name > 2019-04-15T00:00:00Z",
			},
			wantErr: false,
		},
		{
			name: "value is number",
			args: args{
				field: "Name",
				value: 10,
			},
			want: &WhereClause{
				expression: "Name > 10",
			},
			wantErr: false,
		},
		{
			name: "value is date and equal to or greater",
			args: args{
				field:  "Name",
				value:  time.Date(2019, time.April, 15, 0, 0, 0, 0, time.UTC),
				equals: true,
			},
			want: &WhereClause{
				expression: "Name >= 2019-04-15T00:00:00Z",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WhereGreaterThan(tt.args.field, tt.args.value, tt.args.equals)
			if (err != nil) != tt.wantErr {
				t.Errorf("WhereGreaterThan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WhereGreaterThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWhereLessThan(t *testing.T) {
	type args struct {
		field  string
		value  interface{}
		equals bool
	}
	tests := []struct {
		name    string
		args    args
		want    *WhereClause
		wantErr bool
	}{
		{
			name: "value is string",
			args: args{
				field: "Name",
				value: "some string",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "value is boolean",
			args: args{
				field: "Name",
				value: false,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "value is date",
			args: args{
				field: "Name",
				value: time.Date(2019, time.April, 15, 0, 0, 0, 0, time.UTC),
			},
			want: &WhereClause{
				expression: "Name < 2019-04-15T00:00:00Z",
			},
			wantErr: false,
		},
		{
			name: "value is number",
			args: args{
				field: "Name",
				value: 10,
			},
			want: &WhereClause{
				expression: "Name < 10",
			},
			wantErr: false,
		},
		{
			name: "value is date and equal to or less",
			args: args{
				field:  "Name",
				value:  time.Date(2019, time.April, 15, 0, 0, 0, 0, time.UTC),
				equals: true,
			},
			want: &WhereClause{
				expression: "Name <= 2019-04-15T00:00:00Z",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WhereLessThan(tt.args.field, tt.args.value, tt.args.equals)
			if (err != nil) != tt.wantErr {
				t.Errorf("WhereLessThan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WhereLessThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWhereEquals(t *testing.T) {
	type args struct {
		field string
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *WhereClause
		wantErr bool
	}{
		{
			name: "string",
			args: args{
				field: "Name",
				value: "Yeah",
			},
			want: &WhereClause{
				expression: "Name = 'Yeah'",
			},
			wantErr: false,
		},
		{
			name: "boolean",
			args: args{
				field: "Name",
				value: true,
			},
			want: &WhereClause{
				expression: "Name = true",
			},
			wantErr: false,
		},
		{
			name: "number",
			args: args{
				field: "Name",
				value: 10,
			},
			want: &WhereClause{
				expression: "Name = 10",
			},
			wantErr: false,
		},
		{
			name: "date",
			args: args{
				field: "Name",
				value: time.Date(2019, time.April, 15, 0, 0, 0, 0, time.UTC),
			},
			want: &WhereClause{
				expression: "Name = 2019-04-15T00:00:00Z",
			},
			wantErr: false,
		},
		{
			name: "null",
			args: args{
				field: "Name",
				value: nil,
			},
			want: &WhereClause{
				expression: "Name = null",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WhereEquals(tt.args.field, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("WhereEquals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WhereEquals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWhereNotEquals(t *testing.T) {
	type args struct {
		field string
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *WhereClause
		wantErr bool
	}{
		{
			name: "string",
			args: args{
				field: "Name",
				value: "Yeah",
			},
			want: &WhereClause{
				expression: "Name != 'Yeah'",
			},
			wantErr: false,
		},
		{
			name: "boolean",
			args: args{
				field: "Name",
				value: true,
			},
			want: &WhereClause{
				expression: "Name != true",
			},
			wantErr: false,
		},
		{
			name: "number",
			args: args{
				field: "Name",
				value: 10,
			},
			want: &WhereClause{
				expression: "Name != 10",
			},
			wantErr: false,
		},
		{
			name: "date",
			args: args{
				field: "Name",
				value: time.Date(2019, time.April, 15, 0, 0, 0, 0, time.UTC),
			},
			want: &WhereClause{
				expression: "Name != 2019-04-15T00:00:00Z",
			},
			wantErr: false,
		},
		{
			name: "null",
			args: args{
				field: "Name",
				value: nil,
			},
			want: &WhereClause{
				expression: "Name != null",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WhereNotEquals(tt.args.field, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("WhereNotEquals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WhereNotEquals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWhereIn(t *testing.T) {
	type args struct {
		field  string
		values []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *WhereClause
		wantErr bool
	}{
		{
			name: "Strings",
			args: args{
				field: "Name",
				values: []interface{}{
					"Yeah",
					"Yep",
					"Yes",
				},
			},
			want: &WhereClause{
				expression: "Name IN ('Yeah','Yep','Yes')",
			},
			wantErr: false,
		},
		{
			name: "Numbers",
			args: args{
				field: "Name",
				values: []interface{}{
					7,
					8,
					9,
				},
			},
			want: &WhereClause{
				expression: "Name IN (7,8,9)",
			},
			wantErr: false,
		},
		{
			name: "Dates",
			args: args{
				field: "Name",
				values: []interface{}{
					time.Date(2019, time.April, 15, 0, 0, 0, 0, time.UTC),
					time.Date(2018, time.April, 15, 0, 0, 0, 0, time.UTC),
					time.Date(2017, time.April, 15, 0, 0, 0, 0, time.UTC),
				},
			},
			want: &WhereClause{
				expression: "Name IN (2019-04-15T00:00:00Z,2018-04-15T00:00:00Z,2017-04-15T00:00:00Z)",
			},
			wantErr: false,
		},
		{
			name: "boolean",
			args: args{
				field: "Name",
				values: []interface{}{
					true,
					false,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WhereIn(tt.args.field, tt.args.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("WhereIn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WhereIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWhereNotIn(t *testing.T) {
	type args struct {
		field  string
		values []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *WhereClause
		wantErr bool
	}{
		{
			name: "Strings",
			args: args{
				field: "Name",
				values: []interface{}{
					"Yeah",
					"Yep",
					"Yes",
				},
			},
			want: &WhereClause{
				expression: "Name NOT IN ('Yeah','Yep','Yes')",
			},
			wantErr: false,
		},
		{
			name: "Numbers",
			args: args{
				field: "Name",
				values: []interface{}{
					7,
					8,
					9,
				},
			},
			want: &WhereClause{
				expression: "Name NOT IN (7,8,9)",
			},
			wantErr: false,
		},
		{
			name: "Dates",
			args: args{
				field: "Name",
				values: []interface{}{
					time.Date(2019, time.April, 15, 0, 0, 0, 0, time.UTC),
					time.Date(2018, time.April, 15, 0, 0, 0, 0, time.UTC),
					time.Date(2017, time.April, 15, 0, 0, 0, 0, time.UTC),
				},
			},
			want: &WhereClause{
				expression: "Name NOT IN (2019-04-15T00:00:00Z,2018-04-15T00:00:00Z,2017-04-15T00:00:00Z)",
			},
			wantErr: false,
		},
		{
			name: "boolean",
			args: args{
				field: "Name",
				values: []interface{}{
					true,
					false,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WhereNotIn(tt.args.field, tt.args.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("WhereNotIn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WhereNotIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWhereClause_Clause(t *testing.T) {
	type fields struct {
		expression string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Where Cause",
			fields: fields{
				expression: "Name = 'Yeah'",
			},
			want: "WHERE Name = 'Yeah'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wc := &WhereClause{
				expression: tt.fields.expression,
			}
			if got := wc.Clause(); got != tt.want {
				t.Errorf("WhereClause.Clause() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWhereClause_Group(t *testing.T) {
	type fields struct {
		expression string
	}
	tests := []struct {
		name   string
		fields fields
		want   *WhereClause
	}{
		{
			name: "Grouping",
			fields: fields{
				expression: "Name = 'Yeah'",
			},
			want: &WhereClause{
				expression: "(Name = 'Yeah')",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wc := &WhereClause{
				expression: tt.fields.expression,
			}
			wc.Group()
			if !reflect.DeepEqual(wc, tt.want) {
				t.Errorf("WhereClause.Group() = %v, want %v", wc, tt.want)
			}
		})
	}
}

func TestWhereClause_And(t *testing.T) {
	type fields struct {
		expression string
	}
	type args struct {
		where WhereExpression
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *WhereClause
	}{
		{
			name: "Anding",
			fields: fields{
				expression: "FirstName = 'Super'",
			},
			args: args{
				where: &WhereClause{
					expression: "LastName = 'Gary'",
				},
			},
			want: &WhereClause{
				expression: "FirstName = 'Super' AND LastName = 'Gary'",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wc := &WhereClause{
				expression: tt.fields.expression,
			}
			wc.And(tt.args.where)
			if !reflect.DeepEqual(wc, tt.want) {
				t.Errorf("WhereClause.And() = %v, want %v", wc, tt.want)
			}
		})
	}
}

func TestWhereClause_Or(t *testing.T) {
	type fields struct {
		expression string
	}
	type args struct {
		where WhereExpression
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *WhereClause
	}{
		{
			name: "Oring",
			fields: fields{
				expression: "FirstName = 'Super'",
			},
			args: args{
				where: &WhereClause{
					expression: "LastName = 'Gary'",
				},
			},
			want: &WhereClause{
				expression: "FirstName = 'Super' OR LastName = 'Gary'",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wc := &WhereClause{
				expression: tt.fields.expression,
			}
			wc.Or(tt.args.where)
			if !reflect.DeepEqual(wc, tt.want) {
				t.Errorf("WhereClause.Or() = %v, want %v", wc, tt.want)
			}
		})
	}
}

func TestWhereClause_Expression(t *testing.T) {
	type fields struct {
		expression string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Expression",
			fields: fields{
				expression: "Name = 'Yeah'",
			},
			want: "Name = 'Yeah'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wc := &WhereClause{
				expression: tt.fields.expression,
			}
			if got := wc.Expression(); got != tt.want {
				t.Errorf("WhereClause.Expression() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Query(t *testing.T) {
	type fields struct {
		fieldList  []string
		objectType string
		subQuery   []Querier
		where      WhereClauser
		order      Orderer
		limit      int
		offset     int
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name:    "No Object",
			fields:  fields{},
			want:    "",
			wantErr: true,
		},
		{
			name: "No Fields",
			fields: fields{
				objectType: "Account",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Basic",
			fields: fields{
				objectType: "Account",
				fieldList: []string{
					"Name",
					"CreatedBy",
				},
			},
			want:    "SELECT Name,CreatedBy FROM Account",
			wantErr: false,
		},
		{
			name: "Sub Query",
			fields: fields{
				objectType: "Account",
				fieldList: []string{
					"Name",
					"CreatedBy",
				},
				subQuery: []Querier{
					&Builder{
						objectType: "Contacts",
						fieldList: []string{
							"LastName",
						},
					},
				},
			},
			want:    "SELECT Name,CreatedBy,(SELECT LastName FROM Contacts) FROM Account",
			wantErr: false,
		},
		{
			name: "Where",
			fields: fields{
				objectType: "Account",
				fieldList: []string{
					"Name",
					"CreatedBy",
				},
				where: &WhereClause{
					expression: "Name = 'Super Gary'",
				},
			},
			want:    "SELECT Name,CreatedBy FROM Account WHERE Name = 'Super Gary'",
			wantErr: false,
		},
		{
			name: "Order By",
			fields: fields{
				objectType: "Account",
				fieldList: []string{
					"Name",
					"CreatedBy",
				},
				order: &OrderBy{
					fieldOrder: []string{
						"Name",
					},
					result: OrderAsc,
				},
			},
			want:    "SELECT Name,CreatedBy FROM Account ORDER BY Name ASC",
			wantErr: false,
		},
		{
			name: "Limit",
			fields: fields{
				objectType: "Account",
				fieldList: []string{
					"Name",
					"CreatedBy",
				},
				limit: 100,
			},
			want:    "SELECT Name,CreatedBy FROM Account LIMIT 100",
			wantErr: false,
		},
		{
			name: "Offset",
			fields: fields{
				objectType: "Account",
				fieldList: []string{
					"Name",
					"CreatedBy",
				},
				offset: 150,
			},
			want:    "SELECT Name,CreatedBy FROM Account OFFSET 150",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				fieldList:  tt.fields.fieldList,
				objectType: tt.fields.objectType,
				subQuery:   tt.fields.subQuery,
				where:      tt.fields.where,
				order:      tt.fields.order,
				limit:      tt.fields.limit,
				offset:     tt.fields.offset,
			}
			got, err := b.Query()
			if (err != nil) != tt.wantErr {
				t.Errorf("Builder.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Builder.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBuilder(t *testing.T) {
	type args struct {
		input QueryInput
	}
	tests := []struct {
		name    string
		args    args
		want    *Builder
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				input: QueryInput{
					ObjectType: "Account",
					FieldList: []string{
						"Name",
						"Id",
					},
				},
			},
			want: &Builder{
				objectType: "Account",
				fieldList: []string{
					"Name",
					"Id",
				},
			},
			wantErr: false,
		},
		{
			name: "no object type",
			args: args{
				input: QueryInput{
					FieldList: []string{
						"Name",
						"Id",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no field list",
			args: args{
				input: QueryInput{
					ObjectType: "Account",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBuilder(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBuilder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}
