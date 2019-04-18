package collections

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/g8rswimmer/goforce"
	"github.com/g8rswimmer/goforce/session"
	"github.com/g8rswimmer/goforce/sobject"
)

func Test_collection_send(t *testing.T) {
	type fields struct {
		method   string
		endpoint string
		values   *url.Values
		body     io.Reader
	}
	type args struct {
		session session.ServiceFormatter
		value   interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Collection get with values",
			fields: fields{
				method:   http.MethodGet,
				endpoint: "/some/cool/endpoint",
				values: &url.Values{
					"one": []string{"this is fun"},
					"two": []string{"whatever,"},
				},
				body: nil,
			},
			args: args{
				session: &mockSessionFormatter{
					url: "something.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {

						if strings.HasPrefix(req.URL.String(), "something.com/some/cool/endpoint") == false {
							return &http.Response{
								StatusCode: 500,
								Status:     "Bad URL: " + req.URL.String(),
								Body:       ioutil.NopCloser(strings.NewReader("resp")),
								Header:     make(http.Header),
							}
						}

						if req.Method != http.MethodGet {
							return &http.Response{
								StatusCode: 500,
								Status:     "Bad Method",
								Body:       ioutil.NopCloser(strings.NewReader("resp")),
								Header:     make(http.Header),
							}
						}

						values := req.URL.Query()
						if _, ok := values["one"]; ok == false {

							return &http.Response{
								StatusCode: 500,
								Status:     "No one value",
								Body:       ioutil.NopCloser(strings.NewReader("resp")),
								Header:     make(http.Header),
							}
						}
						if _, ok := values["two"]; ok == false {

							return &http.Response{
								StatusCode: 500,
								Status:     "No two value",
								Body:       ioutil.NopCloser(strings.NewReader("resp")),
								Header:     make(http.Header),
							}
						}
						resp := `
							{
								"message" : "Email: invalid email address: Not a real email address",
								"errorCode" : "INVALID_EMAIL_ADDRESS",
								"fields" : [ "Email" ]
							}`

						return &http.Response{
							StatusCode: http.StatusOK,
							Status:     "Some Status",
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
				value: &goforce.Error{},
			},
			wantErr: false,
		},
		{
			name: "Collection with errors",
			fields: fields{
				method:   http.MethodGet,
				endpoint: "/some/cool/endpoint",
				values: &url.Values{
					"one": []string{"this is fun"},
					"two": []string{"whatever,"},
				},
				body: nil,
			},
			args: args{
				session: &mockSessionFormatter{
					url: "something.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {

						if strings.HasPrefix(req.URL.String(), "something.com/some/cool/endpoint") == false {
							return &http.Response{
								StatusCode: 500,
								Status:     "Bad URL: " + req.URL.String(),
								Body:       ioutil.NopCloser(strings.NewReader("resp")),
								Header:     make(http.Header),
							}
						}

						if req.Method != http.MethodGet {
							return &http.Response{
								StatusCode: 500,
								Status:     "Bad Method",
								Body:       ioutil.NopCloser(strings.NewReader("resp")),
								Header:     make(http.Header),
							}
						}

						values := req.URL.Query()
						if _, ok := values["one"]; ok == false {

							return &http.Response{
								StatusCode: 500,
								Status:     "No one value",
								Body:       ioutil.NopCloser(strings.NewReader("resp")),
								Header:     make(http.Header),
							}
						}
						if _, ok := values["two"]; ok == false {

							return &http.Response{
								StatusCode: 500,
								Status:     "No two value",
								Body:       ioutil.NopCloser(strings.NewReader("resp")),
								Header:     make(http.Header),
							}
						}
						resp := `
						[
							{
								"message" : "Email: invalid email address: Not a real email address",
								"errorCode" : "INVALID_EMAIL_ADDRESS",
								"fields" : [ "Email" ]
							}
						]`

						return &http.Response{
							StatusCode: http.StatusConflict,
							Status:     "Some Status",
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
				value: &goforce.Error{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &collection{
				method:   tt.fields.method,
				endpoint: tt.fields.endpoint,
				values:   tt.fields.values,
				body:     tt.fields.body,
			}
			if err := c.send(tt.args.session, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("collection.send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewResource(t *testing.T) {
	type args struct {
		session session.ServiceFormatter
	}
	tests := []struct {
		name string
		args args
		want *Resource
	}{
		{
			name: "get resource",
			args: args{
				session: &mockSessionFormatter{
					url: "some.url.com",
				},
			},
			want: &Resource{
				session: &mockSessionFormatter{
					url: "some.url.com",
				},
				update: &update{
					session: &mockSessionFormatter{
						url: "some.url.com",
					},
				},
				query: &query{
					session: &mockSessionFormatter{
						url: "some.url.com",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResource(tt.args.session); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResource_NewInsert(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
	}
	tests := []struct {
		name   string
		fields fields
		want   *Insert
	}{
		{
			name: "get resource",
			fields: fields{
				session: &mockSessionFormatter{
					url: "some.url.com",
				},
			},
			want: &Insert{
				session: &mockSessionFormatter{
					url: "some.url.com",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resource{
				session: tt.fields.session,
			}
			if got := r.NewInsert(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resource.NewInsert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResource_NewDelete(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
	}
	tests := []struct {
		name   string
		fields fields
		want   *Delete
	}{
		{
			name: "get resource",
			fields: fields{
				session: &mockSessionFormatter{
					url: "some.url.com",
				},
			},
			want: &Delete{
				session: &mockSessionFormatter{
					url: "some.url.com",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resource{
				session: tt.fields.session,
			}
			if got := r.NewDelete(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resource.NewDelete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResource_Update(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		update  *update
	}
	type args struct {
		allOrNone bool
		records   []sobject.Updater
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []UpdateValue
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				update: &update{
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

							if req.Method != http.MethodPatch {
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
										 "statusCode" : "MALFORMED_ID",
										 "message" : "Use one of these records?",
										 "fields" : [ "Id" ]
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
			},
			args: args{
				allOrNone: false,
				records: []sobject.Updater{
					&mockUpdater{
						sobject: "Account",
						fields: map[string]interface{}{
							"NumberOfEmployees": 27000,
						},
						id: "001xx000003DGb2AAG",
					},
					&mockUpdater{
						sobject: "Contact",
						fields: map[string]interface{}{
							"Title": "Lead Engineer",
						},
						id: "003xx000004TmiQAAS",
					},
				},
			},
			want: []UpdateValue{
				UpdateValue{
					InsertValue: sobject.InsertValue{
						Success: false,
						Errors: []goforce.Error{
							{
								ErrorCode: "MALFORMED_ID",
								Message:   "Use one of these records?",
								Fields:    []string{"Id"},
							},
						},
					},
				},
				UpdateValue{
					InsertValue: sobject.InsertValue{
						Success: true,
						ID:      "003RM0000068xVCYAY",
						Errors:  make([]goforce.Error, 0),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "No records",
			fields: fields{
				update: &update{},
			},
			args:    args{},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "No update",
			fields:  fields{},
			args:    args{},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resource{
				session: tt.fields.session,
				update:  tt.fields.update,
			}
			got, err := r.Update(tt.args.allOrNone, tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resource.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resource.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResource_Query(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		update  *update
		query   *query
	}
	type args struct {
		sobject string
		records []sobject.Querier
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				query: &query{
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
			},
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
			wantErr: false,
		},
		{
			name:    "not initialized",
			fields:  fields{},
			args:    args{},
			wantErr: true,
		},
		{
			name: "no records",
			fields: fields{
				query: &query{},
			},
			args:    args{},
			wantErr: true,
		},
		{
			name: "no records",
			fields: fields{
				query: &query{},
			},
			args: args{
				records: make([]sobject.Querier, 2),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resource{
				session: tt.fields.session,
				update:  tt.fields.update,
				query:   tt.fields.query,
			}
			_, err := r.Query(tt.args.sobject, tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resource.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
