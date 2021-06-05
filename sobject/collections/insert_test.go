package collections

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/nettolicious/go-sfdc"
	"github.com/nettolicious/go-sfdc/session"
	"github.com/nettolicious/go-sfdc/sobject"
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
	}
	type args struct {
		allOrNone bool
		records   []sobject.Inserter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "payload",
			fields: fields{},
			args: args{
				allOrNone: false,
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &insert{
				session: tt.fields.session,
			}
			_, err := i.payload(tt.args.allOrNone, tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("Insert.payload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestInsert_Callout(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
	}
	type args struct {
		allOrNone bool
		records   []sobject.Inserter
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
			want: []sobject.InsertValue{
				{
					Success: false,
					Errors: []sfdc.Error{
						{
							ErrorCode: "DUPLICATES_DETECTED",
							Message:   "Use one of these records?",
							Fields:    make([]string, 0),
						},
					},
				},
				{
					Success: true,
					ID:      "003RM0000068xVCYAY",
					Errors:  make([]sfdc.Error, 0),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &insert{
				session: tt.fields.session,
			}
			got, err := i.callout(tt.args.allOrNone, tt.args.records)
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
