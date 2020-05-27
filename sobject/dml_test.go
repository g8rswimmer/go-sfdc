package sobject

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/g8rswimmer/go-sfdc"

	"github.com/g8rswimmer/go-sfdc/session"
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
		session session.ServiceFormatter
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
				session: &mockSessionFormatter{
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
				session: &mockSessionFormatter{
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
				session: &mockSessionFormatter{
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
				session: &mockSessionFormatter{
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
				Errors:  make([]sfdc.Error, 0),
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
			got, err := d.insertCallout(tt.args.inserter)
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
		session session.ServiceFormatter
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
				session: &mockSessionFormatter{
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
				session: &mockSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {

						resp := `[
							{
								"message" : "The requested resource does not exist",
								"errorCode" : "NOT_FOUND"
							}							
						]
						`
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
				session: &mockSessionFormatter{
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
			if err := d.updateCallout(tt.args.updater); (err != nil) != tt.wantErr {
				t.Errorf("dml.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type mockUpsert struct {
	sobject  string
	id       string
	fields   map[string]interface{}
	external string
}

func (mock *mockUpsert) SObject() string {
	return mock.sobject
}
func (mock *mockUpsert) ID() string {
	return mock.id
}
func (mock *mockUpsert) Fields() map[string]interface{} {
	return mock.fields
}
func (mock *mockUpsert) ExternalField() string {
	return mock.external
}
func Test_dml_Upsert(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
	}
	type args struct {
		upserter Upserter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    UpsertValue
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
				upserter: &mockUpsert{
					sobject:  "Account",
					id:       "12345",
					external: "external__c",
					fields: map[string]interface{}{
						"Name":   "Some Test Name",
						"Active": false,
					},
				},
			},
			want:    UpsertValue{},
			wantErr: true,
		},
		{
			name: "Response HTTP Error No JSON",
			fields: fields{
				session: &mockSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {

						resp := `[
							{
								"message" : "The requested resource does not exist",
								"errorCode" : "NOT_FOUND"
							}							
						]
						`
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
				upserter: &mockUpsert{
					sobject:  "Account",
					id:       "12345",
					external: "external__c",
					fields: map[string]interface{}{
						"Name":   "Some Test Name",
						"Active": false,
					},
				},
			},
			want:    UpsertValue{},
			wantErr: true,
		},
		{
			name: "Response HTTP Error JSON",
			fields: fields{
				session: &mockSessionFormatter{
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
				upserter: &mockUpsert{
					sobject:  "Account",
					id:       "12345",
					external: "external__c",
					fields: map[string]interface{}{
						"Name":   "Some Test Name",
						"Active": false,
					},
				},
			},
			want:    UpsertValue{},
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
							StatusCode: http.StatusCreated,
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				upserter: &mockUpsert{
					sobject:  "Account",
					id:       "12345",
					external: "external__c",
					fields: map[string]interface{}{
						"Name":   "Some Test Name",
						"Active": false,
					},
				},
			},
			want:    UpsertValue{},
			wantErr: true,
		},
		{
			name: "Insert Response Passing",
			fields: fields{
				session: &mockSessionFormatter{
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
				upserter: &mockUpsert{
					sobject:  "Account",
					id:       "12345",
					external: "external__c",
					fields: map[string]interface{}{
						"Name":   "Some Test Name",
						"Active": false,
					},
				},
			},
			want: UpsertValue{
				Inserted: true,
				InsertValue: InsertValue{
					Success: true,
					Errors:  make([]sfdc.Error, 0),
					ID:      "001D000000IqhSLIAZ",
				},
			},
			wantErr: false,
		},
		{
			name: "Upsert Response Updated Passing",
			fields: fields{
				session: &mockSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {
						resp := `
						{
							"id" : "001D000000IqhSLIAZ",
							"errors" : [ ],
							"success" : true
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
				upserter: &mockUpsert{
					sobject:  "Account",
					id:       "12345",
					external: "external__c",
					fields: map[string]interface{}{
						"Name":   "Some Test Name",
						"Active": false,
					},
				},
			},
			want: UpsertValue{
				Inserted: false,
				InsertValue: InsertValue{
					Success: true,
					Errors:  make([]sfdc.Error, 0),
					ID:      "001D000000IqhSLIAZ",
				},
			},
			wantErr: false,
		},
		{
			name: "Upsert Response Passing",
			fields: fields{
				session: &mockSessionFormatter{
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
				upserter: &mockUpsert{
					sobject:  "Account",
					id:       "12345",
					external: "external__c",
					fields: map[string]interface{}{
						"Name":   "Some Test Name",
						"Active": false,
					},
				},
			},
			want: UpsertValue{
				Inserted: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dml{
				session: tt.fields.session,
			}
			got, err := d.upsertCallout(tt.args.upserter)
			if (err != nil) != tt.wantErr {
				t.Errorf("dml.Upsert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dml.Upsert() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockDelete struct {
	sobject string
	id      string
}

func (mock *mockDelete) SObject() string {
	return mock.sobject
}
func (mock *mockDelete) ID() string {
	return mock.id
}

func Test_dml_Delete(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
	}
	type args struct {
		deleter Deleter
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
				session: &mockSessionFormatter{
					url: "123://wrong",
				},
			},
			args: args{
				deleter: &mockDelete{
					sobject: "Account",
					id:      "12345",
				},
			},
			wantErr: true,
		},
		{
			name: "Response HTTP Error JSON",
			fields: fields{
				session: &mockSessionFormatter{
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
				deleter: &mockDelete{
					sobject: "Account",
					id:      "12345",
				},
			},
			wantErr: true,
		},
		{
			name: "Passing",
			fields: fields{
				session: &mockSessionFormatter{
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
				deleter: &mockDelete{
					sobject: "Account",
					id:      "12345",
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
			if err := d.deleteCallout(tt.args.deleter); (err != nil) != tt.wantErr {
				t.Errorf("dml.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
