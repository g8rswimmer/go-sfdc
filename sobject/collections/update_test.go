package collections

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/g8rswimmer/go-sfdc"
	"github.com/g8rswimmer/go-sfdc/session"
	"github.com/g8rswimmer/go-sfdc/sobject"
)

type mockUpdater struct {
	sobject string
	fields  map[string]interface{}
	id      string
}

func (mock *mockUpdater) SObject() string {
	return mock.sobject
}
func (mock *mockUpdater) Fields() map[string]interface{} {
	return mock.fields
}
func (mock *mockUpdater) ID() string {
	return mock.id
}

func TestUpdate_payload(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
	}
	type args struct {
		allOrNone bool
		records   []sobject.Updater
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "payloading",
			fields: fields{},
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &update{
				session: tt.fields.session,
			}
			_, err := u.payload(tt.args.allOrNone, tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update.payload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUpdate_Callout(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
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
			args: args{
				allOrNone: true,
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
						Errors: []sfdc.Error{
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
						Errors:  make([]sfdc.Error, 0),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &update{
				session: tt.fields.session,
			}
			got, err := u.callout(tt.args.allOrNone, tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update.Callout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update.Callout() = %v, want %v", got, tt.want)
			}
		})
	}
}
