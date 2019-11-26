package collections

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"go-sfdc/session"
	"./sobject"
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
			q := &query{
				session: tt.fields.session,
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
	}
	type args struct {
		records []sobject.Querier
		sobject string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "payload",
			args: args{
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
			fields:  fields{},
			wantErr: false,
		},
		{
			name:   "payload",
			fields: fields{},
			args: args{
				sobject: "Account",
				records: []sobject.Querier{
					&mockQuery{
						sobject: "Contact",
						id:      "001xx000003DGb1AAG",
						fields: []string{
							"id",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &query{
				session: tt.fields.session,
			}
			_, err := q.payload(tt.args.sobject, tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query.payload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestQuery_Callout(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
	}
	type args struct {
		records []sobject.Querier
		sobject string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
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
			q := &query{
				session: tt.fields.session,
			}
			_, err := q.callout(tt.args.sobject, tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query.Callout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
