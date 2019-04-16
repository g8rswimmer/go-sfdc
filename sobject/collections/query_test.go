package collections

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/g8rswimmer/goforce/session"
	"github.com/g8rswimmer/goforce/sobject"
)

type mockQuery struct {
	sobject string
	id      string
	fields  []string
}

func (mock *mockQuery) SObject() string {
	return mock.sobject
}
func (mock *mockQuery) ID() string {
	return mock.id
}
func (mock *mockQuery) Fields() []string {
	return mock.fields
}

func TestQuery_keyArray(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		records []sobject.Querier
		sobject string
	}
	type args struct {
		m map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name:   "set to array",
			fields: fields{},
			args: args{
				m: map[string]interface{}{
					"001xx000003DGb1AAG": nil,
					"001xx000003DGb0AAG": nil,
					"001xx000003DGb9AAG": nil,
				},
			},
			want: []string{
				"001xx000003DGb1AAG",
				"001xx000003DGb0AAG",
				"001xx000003DGb9AAG",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Query{
				session: tt.fields.session,
				records: tt.fields.records,
				sobject: tt.fields.sobject,
			}
			if got := q.keyArray(tt.args.m); len(got) != len(tt.want) {
				t.Errorf("Query.keyArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_payload(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		records []sobject.Querier
		sobject string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "payload",
			fields: fields{
				sobject: "Account",
				records: []sobject.Querier{
					&mockQuery{
						sobject: "Account",
						id:      "001xx000003DGb1AAG",
						fields: []string{
							"id",
						},
					},
					&mockQuery{
						sobject: "Account",
						id:      "001xx000003DGb0AAG",
						fields: []string{
							"id",
						},
					},
					&mockQuery{
						sobject: "Account",
						id:      "001xx000003DGb9AAG",
						fields: []string{
							"name",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "payload",
			fields: fields{
				sobject: "Account",
				records: []sobject.Querier{
					&mockQuery{
						sobject: "Contacts",
						id:      "001xx000003DGb1AAG",
						fields: []string{
							"id",
						},
					},
					&mockQuery{
						sobject: "Account",
						id:      "001xx000003DGb0AAG",
						fields: []string{
							"id",
						},
					},
					&mockQuery{
						sobject: "Account",
						id:      "001xx000003DGb9AAG",
						fields: []string{
							"name",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Query{
				session: tt.fields.session,
				records: tt.fields.records,
				sobject: tt.fields.sobject,
			}
			_, err := q.payload()
			if (err != nil) != tt.wantErr {
				t.Errorf("Query.payload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestQuery_Records(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		records []sobject.Querier
		sobject string
	}
	type args struct {
		records []sobject.Querier
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Query
	}{
		{
			name:   "records",
			fields: fields{},
			args: args{
				records: []sobject.Querier{
					&mockQuery{
						sobject: "Account",
						id:      "001xx000003DGb1AAG",
						fields: []string{
							"id",
						},
					},
					&mockQuery{
						sobject: "Account",
						id:      "001xx000003DGb0AAG",
						fields: []string{
							"id",
						},
					},
					&mockQuery{
						sobject: "Account",
						id:      "001xx000003DGb9AAG",
						fields: []string{
							"name",
						},
					},
				},
			},
			want: &Query{
				records: []sobject.Querier{
					&mockQuery{
						sobject: "Account",
						id:      "001xx000003DGb1AAG",
						fields: []string{
							"id",
						},
					},
					&mockQuery{
						sobject: "Account",
						id:      "001xx000003DGb0AAG",
						fields: []string{
							"id",
						},
					},
					&mockQuery{
						sobject: "Account",
						id:      "001xx000003DGb9AAG",
						fields: []string{
							"name",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Query{
				session: tt.fields.session,
				records: tt.fields.records,
				sobject: tt.fields.sobject,
			}
			q.Records(tt.args.records...)
			if !reflect.DeepEqual(q, tt.want) {
				t.Errorf("Query.Records() = %v, want %v", q, tt.want)
			}
		})
	}
}

func TestQuery_Callout(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		records []sobject.Querier
		sobject string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				sobject: "Account",
				records: []sobject.Querier{
					&mockQuery{
						sobject: "Account",
						id:      "001xx000003DGb1AAG",
						fields: []string{
							"id",
						},
					},
					&mockQuery{
						sobject: "Account",
						id:      "001xx000003DGb0AAG",
						fields: []string{
							"id",
						},
					},
					&mockQuery{
						sobject: "Account",
						id:      "001xx000003DGb9AAG",
						fields: []string{
							"name",
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
								"attributes" : {
									"type" : "Account",
									"url" : "/services/data/v42.0/sobjects/Account/001xx000003DGb1AAG"
								},
								"Id" : "001xx000003DGb1AAG",
								"Name" : "Acme"
							},
							{
								"attributes" : {
									"type" : "Account",
									"url" : "/services/data/v42.0/sobjects/Account/001xx000003DGb0AAG"
								},
								"Id" : "001xx000003DGb0AAG",
								"Name" : "Global Media"
							},
							null
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Query{
				session: tt.fields.session,
				records: tt.fields.records,
				sobject: tt.fields.sobject,
			}
			_, err := q.Callout()
			if (err != nil) != tt.wantErr {
				t.Errorf("Query.Callout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
