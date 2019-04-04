package sobject

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/g8rswimmer/goforce/session"
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

func Test_dml_Insert(t *testing.T) {
	type fields struct {
		session session.Formatter
	}
	type args struct {
		inserter Inserter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    InsertValue
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
				inserter: &mockInserter{
					sobject: "Account",
					fields: map[string]interface{}{
						"Name":  "Test Account",
						"Email": "something@test.com",
					},
				},
			},
			want:    InsertValue{},
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
				inserter: &mockInserter{
					sobject: "Account",
					fields: map[string]interface{}{
						"Name":  "Test Account",
						"Email": "something@test.com",
					},
				},
			},
			want:    InsertValue{},
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
									"errorCode" : "INVALID_EMAIL_ADDRESS",
									"fields" : [ "Email" ]
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
				inserter: &mockInserter{
					sobject: "Account",
					fields: map[string]interface{}{
						"Name":  "Test Account",
						"Email": "something@test.com",
					},
				},
			},
			want:    InsertValue{},
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
							StatusCode: http.StatusCreated,
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				inserter: &mockInserter{
					sobject: "Account",
					fields: map[string]interface{}{
						"Name":  "Test Account",
						"Email": "something@test.com",
					},
				},
			},
			want:    InsertValue{},
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
							"id" : "001D000000IqhSLIAZ",
							"errors" : [ ],
							"success" : true
						}`

						return &http.Response{
							StatusCode: http.StatusCreated,
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				inserter: &mockInserter{
					sobject: "Account",
					fields: map[string]interface{}{
						"Name":  "Test Account",
						"Email": "something@test.com",
					},
				},
			},
			want: InsertValue{
				Success: true,
				Errors:  make([]string, 0),
				ID:      "001D000000IqhSLIAZ",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dml{
				session: tt.fields.session,
			}
			got, err := d.Insert(tt.args.inserter)
			if (err != nil) != tt.wantErr {
				t.Errorf("dml.Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dml.Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockUpdate struct {
	sobject string
	id      string
	fields  map[string]interface{}
}

func (mock *mockUpdate) SObject() string {
	return mock.sobject
}
func (mock *mockUpdate) ID() string {
	return mock.id
}
func (mock *mockUpdate) Fields() map[string]interface{} {
	return mock.fields
}

func Test_dml_Update(t *testing.T) {
	type fields struct {
		session session.Formatter
	}
	type args struct {
		updater Updater
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
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
				updater: &mockUpdate{
					sobject: "Account",
					id:      "someid",
					fields: map[string]interface{}{
						"Name":   "Some Test Name",
						"Active": false,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Response HTTP Error",
			fields: fields{
				session: &mockMetadataSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {

						return &http.Response{
							StatusCode: 500,
							Status:     "Some Status",
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				updater: &mockUpdate{
					sobject: "Account",
					id:      "someid",
					fields: map[string]interface{}{
						"Name":   "Some Test Name",
						"Active": false,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Response Passing",
			fields: fields{
				session: &mockMetadataSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {
						return &http.Response{
							StatusCode: http.StatusNoContent,
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				updater: &mockUpdate{
					sobject: "Account",
					id:      "someid",
					fields: map[string]interface{}{
						"Name":   "Some Test Name",
						"Active": false,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dml{
				session: tt.fields.session,
			}
			if err := d.Update(tt.args.updater); (err != nil) != tt.wantErr {
				t.Errorf("dml.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
