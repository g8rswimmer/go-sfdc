package batch

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/g8rswimmer/go-sfdc/session"
)

type mockSubrequester struct {
	url             string
	method          string
	richInput       map[string]interface{}
	binaryPartName  string
	binaryPartAlais string
}

func (mock *mockSubrequester) URL() string {
	return mock.url
}
func (mock *mockSubrequester) Method() string {
	return mock.method
}
func (mock *mockSubrequester) BinaryPartName() string {
	return mock.binaryPartName
}
func (mock *mockSubrequester) BinaryPartNameAlias() string {
	return mock.binaryPartAlais
}
func (mock *mockSubrequester) RichInput() map[string]interface{} {
	return mock.richInput
}

func TestResource_validateSubrequests(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
	}
	type args struct {
		requesters []Subrequester
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Valid",
			fields: fields{},
			args: args{
				requesters: []Subrequester{
					&mockSubrequester{
						url:    "www.something.com",
						method: http.MethodGet,
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "No Url",
			fields: fields{},
			args: args{
				requesters: []Subrequester{
					&mockSubrequester{
						method: http.MethodGet,
					},
				},
			},
			wantErr: true,
		},
		{
			name:   "Valid",
			fields: fields{},
			args: args{
				requesters: []Subrequester{
					&mockSubrequester{
						url:    "www.something.com",
						method: http.MethodTrace,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resource{
				session: tt.fields.session,
			}
			if err := r.validateSubrequests(tt.args.requesters); (err != nil) != tt.wantErr {
				t.Errorf("Resource.validateSubrequests() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResource_payload(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
	}
	type args struct {
		haltOnError bool
		requesters  []Subrequester
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *bytes.Reader
		wantErr bool
	}{
		{
			name:   "success",
			fields: fields{},
			args: args{
				haltOnError: true,
				requesters: []Subrequester{
					&mockSubrequester{
						url:             "www.something.com",
						method:          http.MethodGet,
						binaryPartAlais: "Alais",
						binaryPartName:  "Name",
						richInput: map[string]interface{}{
							"Name": "NewName",
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
			_, err := r.payload(tt.args.haltOnError, tt.args.requesters)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resource.payload() error = %v, wantErr %v", err, tt.wantErr)
				return
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
			name: "no session",
			args: args{
				session: nil,
			},
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

func TestResource_Retrieve(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
	}
	type args struct {
		haltOnError bool
		requesters  []Subrequester
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Value
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				session: &mockSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {
						if req.URL.String() != "https://test.salesforce.com/composite/batch" {
							return &http.Response{
								StatusCode: 500,
								Status:     "Invalid URL",
								Body:       ioutil.NopCloser(strings.NewReader(req.URL.String())),
								Header:     make(http.Header),
							}
						}

						resp := `{
							"hasErrors" : false,
							"results" : [{
							   "statusCode" : 204,
							   "result" : null
							   },{
							   "statusCode" : 200,
							   "result": {
								  "attributes" : {
									 "type" : "Account",
									 "url" : "/services/data/v34.0/sobjects/Account/001D000000K0fXOIAZ"
								  },
								  "Name" : "NewName",
								  "BillingPostalCode" : "94105",
								  "Id" : "001D000000K0fXOIAZ"
							   }
							}]
						 }`
						return &http.Response{
							StatusCode: http.StatusOK,
							Status:     "Good",
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}

					}),
				},
			},
			args: args{
				haltOnError: false,
				requesters: []Subrequester{
					&mockSubrequester{
						url:    "www.something.com",
						method: http.MethodGet,
					},
				},
			},
			want: Value{
				HasErrors: false,
				Results: []Subvalue{
					{
						StatusCode: 204,
					},
					{
						Result: map[string]interface{}{
							"attributes": map[string]interface{}{
								"type": "Account",
								"url":  "/services/data/v34.0/sobjects/Account/001D000000K0fXOIAZ",
							},
							"Name":              "NewName",
							"BillingPostalCode": "94105",
							"Id":                "001D000000K0fXOIAZ",
						},
						StatusCode: 200,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Errors",
			fields: fields{
				session: &mockSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {
						if req.URL.String() != "https://test.salesforce.com/composite/batch" {
							return &http.Response{
								StatusCode: 500,
								Status:     "Invalid URL",
								Body:       ioutil.NopCloser(strings.NewReader(req.URL.String())),
								Header:     make(http.Header),
							}
						}

						resp := `[
							{
								"fields" : [ "Id" ],
								"message" : "Account ID: id value of incorrect type: 001900K0001pPuOAAU",
								"errorCode" : "MALFORMED_ID"
							}							
						]`
						return &http.Response{
							StatusCode: http.StatusBadRequest,
							Status:     "Bad",
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}

					}),
				},
			},
			args: args{
				haltOnError: false,
				requesters: []Subrequester{
					&mockSubrequester{
						url:    "www.something.com",
						method: http.MethodGet,
					},
				},
			},
			want:    Value{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resource{
				session: tt.fields.session,
			}
			got, err := r.Retrieve(tt.args.haltOnError, tt.args.requesters)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resource.Retrieve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resource.Retrieve() = %v, want %v", got, tt.want)
			}
		})
	}
}
