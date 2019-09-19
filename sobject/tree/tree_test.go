package tree

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/aheber/go-sfdc"
	"github.com/aheber/go-sfdc/session"
)

type mockInserter struct {
	sobject string
	records []*Record
}

func (mock *mockInserter) SObject() string {
	return mock.sobject
}
func (mock *mockInserter) Records() []*Record {
	return mock.records
}
func TestResource_Insert(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
	}
	type args struct {
		inserter Inserter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Value
		wantErr bool
	}{
		{
			name: "No Inserter",
			fields: fields{
				session: &mockSessionFormatter{
					url: "123://wrong",
				},
			},
			args:    args{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Request Error",
			fields: fields{
				session: &mockSessionFormatter{
					url: "123://wrong",
				},
			},
			args: args{
				inserter: &mockInserter{
					sobject: "Account",
					records: []*Record{
						{
							Attributes: Attributes{
								Type:        "Account",
								ReferenceID: "ref1",
							},
							Fields: map[string]interface{}{
								"name":              "SampleAccount",
								"phone":             "1234567890",
								"website":           "www.salesforce.com",
								"numberOfEmployees": 100,
								"industry":          "Banking",
							},
							Records: make(map[string][]*Record),
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Response HTTP Error No JSON",
			fields: fields{
				session: &mockSessionFormatter{
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
				inserter: &mockInserter{
					sobject: "Account",
					records: []*Record{
						{
							Attributes: Attributes{
								Type:        "Account",
								ReferenceID: "ref1",
							},
							Fields: map[string]interface{}{
								"name":              "SampleAccount",
								"phone":             "1234567890",
								"website":           "www.salesforce.com",
								"numberOfEmployees": 100,
								"industry":          "Banking",
							},
							Records: make(map[string][]*Record),
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				session: &mockSessionFormatter{
					url: "something.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {

						if strings.HasPrefix(req.URL.String(), "something.com/composite/tree/Account") == false {
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
						{
							"hasErrors" : false,
							"results" : [{
							 "referenceId" : "ref1",
							 "id" : "001D000000K0fXOIAZ"
							 }]
						}`

						return &http.Response{
							StatusCode: http.StatusCreated,
							Status:     "Some Status",
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				inserter: &mockInserter{
					sobject: "Account",
					records: []*Record{
						{
							Attributes: Attributes{
								Type:        "Account",
								ReferenceID: "ref1",
							},
							Fields: map[string]interface{}{
								"name":              "SampleAccount",
								"phone":             "1234567890",
								"website":           "www.salesforce.com",
								"numberOfEmployees": 100,
								"industry":          "Banking",
							},
							Records: make(map[string][]*Record),
						},
					},
				},
			},
			want: &Value{
				HasErrors: false,
				Results: []InsertValue{
					{
						ID:          "001D000000K0fXOIAZ",
						ReferenceID: "ref1",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "success",
			fields: fields{
				session: &mockSessionFormatter{
					url: "something.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {

						if strings.HasPrefix(req.URL.String(), "something.com/composite/tree/Account") == false {
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
						{
							"hasErrors" : true,
							"results" : [{
							  "referenceId" : "ref2",
							  "errors" : [{
								"statusCode" : "INVALID_EMAIL_ADDRESS",
								"message" : "Email: invalid email address: 123",
								"fields" : [ "Email" ]
								}]
							  }]
						 }`

						return &http.Response{
							StatusCode: http.StatusCreated,
							Status:     "Some Status",
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				inserter: &mockInserter{
					sobject: "Account",
					records: []*Record{
						{
							Attributes: Attributes{
								Type:        "Account",
								ReferenceID: "ref1",
							},
							Fields: map[string]interface{}{
								"name":              "SampleAccount",
								"phone":             "1234567890",
								"website":           "www.salesforce.com",
								"numberOfEmployees": 100,
								"industry":          "Banking",
							},
							Records: make(map[string][]*Record),
						},
					},
				},
			},
			want: &Value{
				HasErrors: true,
				Results: []InsertValue{
					{
						ReferenceID: "ref2",
						Errors: []sfdc.Error{
							{
								ErrorCode: "INVALID_EMAIL_ADDRESS",
								Message:   "Email: invalid email address: 123",
								Fields:    []string{"Email"},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resource{
				session: tt.fields.session,
			}
			got, err := r.Insert(tt.args.inserter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resource.Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resource.Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewResource(t *testing.T) {
	type args struct {
		session session.ServiceFormatter
	}
	tests := []struct {
		name    string
		args    args
		want    *Resource
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				session: &mockSessionFormatter{},
			},
			want: &Resource{
				session: &mockSessionFormatter{},
			},
			wantErr: false,
		},
		{
			name:    "no session",
			args:    args{},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewResource(tt.args.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
