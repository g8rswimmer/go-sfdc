package sobject

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/g8rswimmer/goforce"
	"github.com/g8rswimmer/goforce/session"
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

type mockRecord struct {
	sobject string
	url     string
	fields  map[string]interface{}
}

func testNewRecord(data []byte) *goforce.Record {
	var record goforce.Record
	err := json.Unmarshal(data, &record)
	if err != nil {
		return &goforce.Record{}
	}
	return &record
}
func Test_query_Query(t *testing.T) {
	type fields struct {
		session session.Formatter
	}
	type args struct {
		querier Querier
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *goforce.Record
		wantErr bool
	}{
		{
			name: "Request Error",
			fields: fields{
				session: &mockMetadataSessionFormatter{
					url: "123://wrong",
				},
			},
			args: args{
				querier: &mockQuery{
					sobject: "Account",
					fields: []string{
						"Name",
						"Email",
					},
					id: "SomeID",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Response HTTP Error No JSON",
			fields: fields{
				session: &mockMetadataSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {

						return &http.Response{
							StatusCode: 500,
							Status:     "Some Status",
							Body:       ioutil.NopCloser(strings.NewReader("resp")),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				querier: &mockQuery{
					sobject: "Account",
					fields: []string{
						"Name",
						"Email",
					},
					id: "SomeID",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Response HTTP Error JSON",
			fields: fields{
				session: &mockMetadataSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {

						resp := `
							[ 
								{
									"message" : "Email: invalid email address: Not a real email address",
									"errorCode" : "INVALID_EMAIL_ADDRESS"
							  	} 
							]`
						return &http.Response{
							StatusCode: 500,
							Status:     "Some Status",
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				querier: &mockQuery{
					sobject: "Account",
					fields: []string{
						"Name",
						"Email",
					},
					id: "SomeID",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Response JSON Error",
			fields: fields{
				session: &mockMetadataSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {
						resp := `
						{`

						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				querier: &mockQuery{
					sobject: "Account",
					fields: []string{
						"Name",
						"Email",
					},
					id: "SomeID",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Response Passing",
			fields: fields{
				session: &mockMetadataSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {
						resp := `
						{
							"AccountNumber" : "CD656092",
							"BillingPostalCode" : "27215"
						}`

						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				querier: &mockQuery{
					sobject: "Account",
					fields: []string{
						"AccountNumber",
						"BillingPostalCode",
					},
					id: "SomeID",
				},
			},
			want: testNewRecord([]byte(`
			{
				"AccountNumber" : "CD656092",
				"BillingPostalCode" : "27215"
			}`)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &query{
				session: tt.fields.session,
			}
			got, err := q.Query(tt.args.querier)
			if (err != nil) != tt.wantErr {
				t.Errorf("query.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("query.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockExternalQuery struct {
	sobject  string
	id       string
	fields   []string
	external string
}

func (mock *mockExternalQuery) SObject() string {
	return mock.sobject
}
func (mock *mockExternalQuery) ID() string {
	return mock.id
}
func (mock *mockExternalQuery) Fields() []string {
	return mock.fields
}
func (mock *mockExternalQuery) ExternalField() string {
	return mock.external
}
func Test_query_ExternalQuery(t *testing.T) {
	type fields struct {
		session session.Formatter
	}
	type args struct {
		querier ExternalQuerier
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *goforce.Record
		wantErr bool
	}{
		{
			name: "Request Error",
			fields: fields{
				session: &mockMetadataSessionFormatter{
					url: "123://wrong",
				},
			},
			args: args{
				querier: &mockExternalQuery{
					sobject: "Account",
					fields: []string{
						"Name",
						"Email",
					},
					id:       "SomeID",
					external: "ExternalField",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Response HTTP Error No JSON",
			fields: fields{
				session: &mockMetadataSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {

						return &http.Response{
							StatusCode: 500,
							Status:     "Some Status",
							Body:       ioutil.NopCloser(strings.NewReader("resp")),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				querier: &mockExternalQuery{
					sobject: "Account",
					fields: []string{
						"Name",
						"Email",
					},
					id:       "SomeID",
					external: "ExternalField",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Response HTTP Error JSON",
			fields: fields{
				session: &mockMetadataSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {

						resp := `
							[ 
								{
									"message" : "Email: invalid email address: Not a real email address",
									"errorCode" : "INVALID_EMAIL_ADDRESS"
							  	} 
							]`
						return &http.Response{
							StatusCode: 500,
							Status:     "Some Status",
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				querier: &mockExternalQuery{
					sobject: "Account",
					fields: []string{
						"Name",
						"Email",
					},
					id:       "SomeID",
					external: "ExternalField",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Response JSON Error",
			fields: fields{
				session: &mockMetadataSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {
						resp := `
						{`

						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				querier: &mockExternalQuery{
					sobject: "Account",
					fields: []string{
						"Name",
						"Email",
					},
					id:       "SomeID",
					external: "ExternalField",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Response Passing",
			fields: fields{
				session: &mockMetadataSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {
						resp := `
						{
							"AccountNumber" : "CD656092",
							"BillingPostalCode" : "27215"
						}`

						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				querier: &mockExternalQuery{
					sobject: "Account",
					fields: []string{
						"Name",
						"Email",
					},
					id:       "SomeID",
					external: "ExternalField",
				},
			},
			want: testNewRecord([]byte(`
			{
				"AccountNumber" : "CD656092",
				"BillingPostalCode" : "27215"
			}`)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &query{
				session: tt.fields.session,
			}
			got, err := q.ExternalQuery(tt.args.querier)
			if (err != nil) != tt.wantErr {
				t.Errorf("query.ExternalQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("query.ExternalQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
