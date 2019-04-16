package collections

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/g8rswimmer/goforce"
	"github.com/g8rswimmer/goforce/session"
	"github.com/g8rswimmer/goforce/sobject"
)

type mockInserter struct {
	sobject string
	fields  map[string]interface{}
}

func (mock *mockInserter) SObject() string {
	return mock.sobject
}
func (mock *mockInserter) Fields() map[string]interface{} {
	return mock.fields
}

func TestInsert_payload(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		records []sobject.Inserter
	}
	type args struct {
		allOrNone bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "payload",
			fields: fields{
				records: []sobject.Inserter{
					&mockInserter{
						sobject: "Account",
						fields: map[string]interface{}{
							"Name":        "example.com",
							"BillingCity": "San Francisco",
						},
					},
					&mockInserter{
						sobject: "Contact",
						fields: map[string]interface{}{
							"LastName":  "Johnson",
							"FirstName": "Erica",
						},
					},
				},
			},
			args: args{
				allOrNone: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Insert{
				session: tt.fields.session,
				records: tt.fields.records,
			}
			_, err := i.payload(tt.args.allOrNone)
			if (err != nil) != tt.wantErr {
				t.Errorf("Insert.payload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestInsert_Records(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		records []sobject.Inserter
	}
	type args struct {
		records []sobject.Inserter
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Insert
	}{
		{
			name:   "add records",
			fields: fields{},
			args: args{
				records: []sobject.Inserter{
					&mockInserter{
						sobject: "Account",
						fields: map[string]interface{}{
							"Name":        "example.com",
							"BillingCity": "San Francisco",
						},
					},
					&mockInserter{
						sobject: "Contact",
						fields: map[string]interface{}{
							"LastName":  "Johnson",
							"FirstName": "Erica",
						},
					},
				},
			},
			want: &Insert{
				records: []sobject.Inserter{
					&mockInserter{
						sobject: "Account",
						fields: map[string]interface{}{
							"Name":        "example.com",
							"BillingCity": "San Francisco",
						},
					},
					&mockInserter{
						sobject: "Contact",
						fields: map[string]interface{}{
							"LastName":  "Johnson",
							"FirstName": "Erica",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Insert{
				session: tt.fields.session,
				records: tt.fields.records,
			}
			i.Records(tt.args.records...)
			if !reflect.DeepEqual(i, tt.want) {
				t.Errorf("Insert.Records() = %v, want %v", i, tt.want)
			}
		})
	}
}

func TestInsert_Callout(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		records []sobject.Inserter
	}
	type args struct {
		allOrNone bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []sobject.InsertValue
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				records: []sobject.Inserter{
					&mockInserter{
						sobject: "Account",
						fields: map[string]interface{}{
							"Name":        "example.com",
							"BillingCity": "San Francisco",
						},
					},
					&mockInserter{
						sobject: "Contact",
						fields: map[string]interface{}{
							"LastName":  "Johnson",
							"FirstName": "Erica",
						},
					},
				},
				session: &mockSessionFormatter{
					url: "something.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {

						if strings.HasPrefix(req.URL.String(), "something.com/composite/sobjects") == false {
							return &http.Response{
								StatusCode: 500,
								Status:     "Bad URL: " + req.URL.String(),
								Body:       ioutil.NopCloser(strings.NewReader("resp")),
								Header:     make(http.Header),
							}
						}

						if req.Method != http.MethodPost {
							return &http.Response{
								StatusCode: 500,
								Status:     "Bad Method",
								Body:       ioutil.NopCloser(strings.NewReader("resp")),
								Header:     make(http.Header),
							}
						}

						resp := `
						[
							{
							   "success" : false,
							   "errors" : [
								  {
									 "statusCode" : "DUPLICATES_DETECTED",
									 "message" : "Use one of these records?",
									 "fields" : [ ]
								  }
							   ]
							},
							{
							   "id" : "003RM0000068xVCYAY",
							   "success" : true,
							   "errors" : [ ]
							}
						 ]`

						return &http.Response{
							StatusCode: http.StatusOK,
							Status:     "Some Status",
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				allOrNone: true,
			},
			want: []sobject.InsertValue{
				{
					Success: false,
					Errors: []goforce.Error{
						{
							StatusCode: "DUPLICATES_DETECTED",
							Message:    "Use one of these records?",
							Fields:     make([]string, 0),
						},
					},
				},
				{
					Success: true,
					ID:      "003RM0000068xVCYAY",
					Errors:  make([]goforce.Error, 0),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Insert{
				session: tt.fields.session,
				records: tt.fields.records,
			}
			got, err := i.Callout(tt.args.allOrNone)
			if (err != nil) != tt.wantErr {
				t.Errorf("Insert.Callout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Insert.Callout() = %v, want %v", got, tt.want)
			}
		})
	}
}
