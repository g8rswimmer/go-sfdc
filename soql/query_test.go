package soql

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/nettolicious/go-sfdc/session"
)

type mockQuerier struct {
	stmt string
	err  error
}

func (mock *mockQuerier) Format() (string, error) {
	return mock.stmt, mock.err
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
			name: "New Resource",
			args: args{
				session: &mockSessionFormatter{
					url: "Something",
				},
			},
			want: &Resource{
				session: &mockSessionFormatter{
					url: "Something",
				},
			},
			wantErr: false,
		},
		{
			name:    "New Resource",
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

func TestResource_Query(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
	}
	type args struct {
		querier QueryFormatter
		all     bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *QueryResult
		wantErr bool
	}{
		{
			name: "Request Error",
			fields: fields{
				session: &mockSessionFormatter{
					url: "123://wrong",
				},
			},
			args: args{
				querier: &mockQuerier{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Response HTTP Error",
			fields: fields{
				session: &mockSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {

						return &http.Response{
							StatusCode: 500,
							Status:     "Some Status",
							Body:       ioutil.NopCloser(strings.NewReader("Error")),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				querier: &mockQuerier{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Response JSON Error",
			fields: fields{
				session: &mockSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {
						resp := `
						{`

						return &http.Response{
							StatusCode: 200,
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				querier: &mockQuerier{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Response Passing",
			fields: fields{
				session: &mockSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {
						if req.URL.String() != "https://test.salesforce.com/query/?q=SELECT+Name+FROM+Account" {
							return &http.Response{
								StatusCode: 500,
								Status:     "Some Status",
								Body:       ioutil.NopCloser(strings.NewReader("Error")),
								Header:     make(http.Header),
							}
						}
						resp := `
						{
							"done" : true,
							"totalSize" : 2,
							"records" : 
							[ 
								{  
									"attributes" : 
									{    
										"type" : "Account",    
										"url" : "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH"  
									},  
									"Name" : "Test 1"
								}, 
								{  
									"attributes" : 
									{    
										"type" : "Account",    
										"url" : "/services/data/v20.0/sobjects/Account/001D000000IomazIAB"  
									},  
									"Name" : "Test 2"
								}
							]
						}`

						return &http.Response{
							StatusCode: 200,
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				querier: &mockQuerier{
					stmt: "SELECT Name FROM Account",
				},
			},
			want: &QueryResult{
				response: queryResponse{
					Done:      true,
					TotalSize: 2,
					Records: []map[string]interface{}{
						{
							"attributes": map[string]interface{}{
								"type": "Account",
								"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
							},
							"Name": "Test 1",
						},
						{
							"attributes": map[string]interface{}{
								"type": "Account",
								"url":  "/services/data/v20.0/sobjects/Account/001D000000IomazIAB",
							},
							"Name": "Test 2",
						},
					},
				},
				records: testNewQueryRecords([]map[string]interface{}{
					{
						"attributes": map[string]interface{}{
							"type": "Account",
							"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
						},
						"Name": "Test 1",
					},
					{
						"attributes": map[string]interface{}{
							"type": "Account",
							"url":  "/services/data/v20.0/sobjects/Account/001D000000IomazIAB",
						},
						"Name": "Test 2",
					},
				}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resource{
				session: tt.fields.session,
			}
			got, err := r.Query(tt.args.querier, tt.args.all)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resource.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				tt.want.resource = r
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resource.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResource_next(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
	}
	type args struct {
		recordURL string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *QueryResult
		wantErr bool
	}{
		{
			name: "Response Passing",
			fields: fields{
				session: &mockSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {
						if req.URL.String() != "https://test.salesforce.com/services/data/v20.0/query/01gD0000002HU6KIAW-2000" {
							return &http.Response{
								StatusCode: 500,
								Status:     "Some Status",
								Body:       ioutil.NopCloser(strings.NewReader("Error")),
								Header:     make(http.Header),
							}
						}
						resp := `
						{
							"done" : true,
							"totalSize" : 2,
							"records" : 
							[ 
								{  
									"attributes" : 
									{    
										"type" : "Account",    
										"url" : "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH"  
									},  
									"Name" : "Test 1"
								}, 
								{  
									"attributes" : 
									{    
										"type" : "Account",    
										"url" : "/services/data/v20.0/sobjects/Account/001D000000IomazIAB"  
									},  
									"Name" : "Test 2"
								}
							]
						}`

						return &http.Response{
							StatusCode: 200,
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				recordURL: "/services/data/v20.0/query/01gD0000002HU6KIAW-2000",
			},
			want: &QueryResult{
				response: queryResponse{
					Done:      true,
					TotalSize: 2,
					Records: []map[string]interface{}{
						{
							"attributes": map[string]interface{}{
								"type": "Account",
								"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
							},
							"Name": "Test 1",
						},
						{
							"attributes": map[string]interface{}{
								"type": "Account",
								"url":  "/services/data/v20.0/sobjects/Account/001D000000IomazIAB",
							},
							"Name": "Test 2",
						},
					},
				},
				records: testNewQueryRecords([]map[string]interface{}{
					{
						"attributes": map[string]interface{}{
							"type": "Account",
							"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
						},
						"Name": "Test 1",
					},
					{
						"attributes": map[string]interface{}{
							"type": "Account",
							"url":  "/services/data/v20.0/sobjects/Account/001D000000IomazIAB",
						},
						"Name": "Test 2",
					},
				}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resource{
				session: tt.fields.session,
			}
			got, err := r.next(tt.args.recordURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resource.next() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				tt.want.resource = r
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resource.next() = %v, want %v", got, tt.want)
			}
		})
	}
}
