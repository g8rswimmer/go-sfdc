package sobject

import (
	"net/http"
	"testing"

	"github.com/g8rswimmer/goforce/session"
)

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
func Test_update_Update(t *testing.T) {
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
			u := &update{
				session: tt.fields.session,
			}
			if err := u.Update(tt.args.updater); (err != nil) != tt.wantErr {
				t.Errorf("update.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
