package composite

import (
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/g8rswimmer/go-sfdc/session"
)

type mockSubrequester struct {
	url         string
	referenceID string
	method      string
	httpHeaders http.Header
	body        map[string]interface{}
}

func (mock *mockSubrequester) URL() string {
	return mock.url
}
func (mock *mockSubrequester) ReferenceID() string {
	return mock.referenceID
}
func (mock *mockSubrequester) Method() string {
	return mock.method
}
func (mock *mockSubrequester) HTTPHeaders() http.Header {
	return mock.httpHeaders
}
func (mock *mockSubrequester) Body() map[string]interface{} {
	return mock.body
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
						url:         "www.something.com",
						referenceID: "someID",
						method:      http.MethodGet,
						httpHeaders: http.Header(map[string][]string{
							"Accept-Language": []string{"en-us,"},
						}),
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "No URL",
			fields: fields{},
			args: args{
				requesters: []Subrequester{
					&mockSubrequester{
						url:         "",
						referenceID: "someID",
						method:      http.MethodGet,
						httpHeaders: http.Header(map[string][]string{
							"Accept-Language": []string{"en-us,"},
						}),
					},
				},
			},
			wantErr: true,
		},
		{
			name:   "No reference ID",
			fields: fields{},
			args: args{
				requesters: []Subrequester{
					&mockSubrequester{
						url:         "www.something.com",
						referenceID: "",
						method:      http.MethodGet,
						httpHeaders: http.Header(map[string][]string{
							"Accept-Language": []string{"en-us,"},
						}),
					},
				},
			},
			wantErr: true,
		},
		{
			name:   "Incorrect method",
			fields: fields{},
			args: args{
				requesters: []Subrequester{
					&mockSubrequester{
						url:         "www.something.com",
						referenceID: "someID",
						method:      http.MethodTrace,
						httpHeaders: http.Header(map[string][]string{
							"Accept-Language": []string{"en-us,"},
						}),
					},
				},
			},
			wantErr: true,
		},
		{
			name:   "No method",
			fields: fields{},
			args: args{
				requesters: []Subrequester{
					&mockSubrequester{
						url:         "www.something.com",
						referenceID: "someID",
						method:      "",
						httpHeaders: http.Header(map[string][]string{
							"Accept-Language": []string{"en-us,"},
						}),
					},
				},
			},
			wantErr: true,
		},
		{
			name:   "Incorrect http header",
			fields: fields{},
			args: args{
				requesters: []Subrequester{
					&mockSubrequester{
						url:         "www.something.com",
						referenceID: "someID",
						method:      http.MethodGet,
						httpHeaders: http.Header(map[string][]string{
							"Accept": []string{"application/json"},
						}),
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
		allOrNone  bool
		requesters []Subrequester
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "success",
			fields: fields{},
			args: args{
				allOrNone: true,
				requesters: []Subrequester{
					&mockSubrequester{
						url:         "www.something.com",
						referenceID: "someID",
						method:      http.MethodGet,
						httpHeaders: http.Header(map[string][]string{
							"Accept-Language": []string{"en-us,"},
						}),
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
			_, err := r.payload(tt.args.allOrNone, tt.args.requesters)
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
		allOrNone  bool
		requesters []Subrequester
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
						if req.URL.String() != "https://test.salesforce.com/composite" {
							return &http.Response{
								StatusCode: 500,
								Status:     "Invalid URL",
								Body:       io.NopCloser(strings.NewReader(req.URL.String())),
								Header:     make(http.Header),
							}
						}

						resp := `{
							"compositeResponse" : [{
								"body" : {
									"id" : "001R00000033JNuIAM",
									"success" : true,
									"errors" : [ ]
								},
								"httpHeaders" : {
								  "Location" : "/services/data/v38.0/sobjects/Account/001R00000033JNuIAM"
								},
								"httpStatusCode" : 201,
								"referenceId" : "NewAccount"
							},{
								"body" : {
									"what": "all the account data"
								},
								"httpHeaders" : {
									"ETag" : "\"Jbjuzw7dbhaEG3fd90kJbx6A0ow=\"",
									"Last-Modified" : "Fri, 22 Jul 2016 20:19:37 GMT"
								},
								"httpStatusCode" : 200,
								"referenceId" : "NewAccountInfo"
							},{
								"body" : {
									"id" : "003R00000025REHIA2",
									"success" : true,
									"errors" : [ ]
								},
								"httpHeaders" : {
									"Location" : "/services/data/v38.0/sobjects/Contact/003R00000025REHIA2"
								},
								"httpStatusCode" : 201,
								"referenceId" : "NewContact"
							},{
								"body" : {
									"attributes" : {
									"type" : "User",
									"url" : "/services/data/v38.0/sobjects/User/005R0000000I90CIAS"
									},
									"Name" : "Jane Doe",
									"CompanyName" : "Salesforce",
									"Title" : "Director",
									"City" : "San Francisco",
									"State" : "CA",
									"Id" : "005R0000000I90CIAS"
								},
								"httpHeaders" : { },
								"httpStatusCode" : 200,
								"referenceId" : "NewAccountOwner"
							},{
								"body" : null,
								"httpHeaders" : {
									"ETag" : "\"f2293620\"",
									"Last-Modified" : "Fri, 22 Jul 2016 18:45:56 GMT"
								 },
								"httpStatusCode" : 304,
								"referenceId" : "AccountMetadata"
							}]
						}`
						return &http.Response{
							StatusCode: http.StatusOK,
							Status:     "Good",
							Body:       io.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}

					}),
				},
			},
			args: args{
				allOrNone: false,
				requesters: []Subrequester{
					&mockSubrequester{
						url:         "www.something.com",
						referenceID: "someID",
						method:      http.MethodGet,
						httpHeaders: nil,
					},
				},
			},
			want: Value{
				Response: []Subvalue{
					{
						Body: map[string]interface{}{
							"id":      "001R00000033JNuIAM",
							"success": true,
							"errors":  make([]interface{}, 0),
						},
						HTTPHeaders: map[string]string{
							"Location": "/services/data/v38.0/sobjects/Account/001R00000033JNuIAM",
						},
						HTTPStatusCode: 201,
						ReferenceID:    "NewAccount",
					},
					{
						Body: map[string]interface{}{
							"what": "all the account data",
						},
						HTTPHeaders: map[string]string{
							"ETag":          "\"Jbjuzw7dbhaEG3fd90kJbx6A0ow=\"",
							"Last-Modified": "Fri, 22 Jul 2016 20:19:37 GMT",
						},
						HTTPStatusCode: 200,
						ReferenceID:    "NewAccountInfo",
					},
					{
						Body: map[string]interface{}{
							"id":      "003R00000025REHIA2",
							"success": true,
							"errors":  make([]interface{}, 0),
						},
						HTTPHeaders: map[string]string{
							"Location": "/services/data/v38.0/sobjects/Contact/003R00000025REHIA2",
						},
						HTTPStatusCode: 201,
						ReferenceID:    "NewContact",
					},
					{
						Body: map[string]interface{}{
							"attributes": map[string]interface{}{
								"type": "User",
								"url":  "/services/data/v38.0/sobjects/User/005R0000000I90CIAS",
							},
							"Name":        "Jane Doe",
							"CompanyName": "Salesforce",
							"Title":       "Director",
							"City":        "San Francisco",
							"State":       "CA",
							"Id":          "005R0000000I90CIAS",
						},
						HTTPHeaders:    map[string]string{},
						HTTPStatusCode: 200,
						ReferenceID:    "NewAccountOwner",
					},
					{
						Body: nil,
						HTTPHeaders: map[string]string{
							"ETag":          "\"f2293620\"",
							"Last-Modified": "Fri, 22 Jul 2016 18:45:56 GMT",
						},
						HTTPStatusCode: 304,
						ReferenceID:    "AccountMetadata",
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
						if req.URL.String() != "https://test.salesforce.com/composite" {
							return &http.Response{
								StatusCode: 500,
								Status:     "Invalid URL",
								Body:       io.NopCloser(strings.NewReader(req.URL.String())),
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
							Body:       io.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}

					}),
				},
			},
			args: args{
				allOrNone: false,
				requesters: []Subrequester{
					&mockSubrequester{
						url:         "www.something.com",
						referenceID: "someID",
						method:      http.MethodGet,
						httpHeaders: nil,
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
			got, err := r.Retrieve(tt.args.allOrNone, tt.args.requesters)
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
