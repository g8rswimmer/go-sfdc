package sobject

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/aheber/go-sfdc/session"
)

func Test_metadata_Metadata(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
	}
	type args struct {
		sobject string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    MetadataValue
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
				sobject: "someobject",
			},
			want:    MetadataValue{},
			wantErr: true,
		},
		{
			name: "Response HTTP Error",
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
				sobject: "someobject",
			},
			want:    MetadataValue{},
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
				sobject: "someobject",
			},
			want:    MetadataValue{},
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
							"objectDescribe": {
								"activateable": false,
								"createable": true,
								"custom": false,
								"customSetting": false,
								"deletable": true,
								"deprecatedAndHidden": false,
								"feedEnabled": true,
								"hasSubtypes": true,
								"isSubtype": false,
								"keyPrefix": "001",
								"label": "Account",
								"labelPlural": "Accounts",
								"layoutable": true,
								"mergeable": true,
								"mruEnabled": true,
								"name": "Account",
								"queryable": true,
								"replicateable": true,
								"retrieveable": true,
								"searchable": true,
								"triggerable": true,
								"undeletable": true,
								"updateable": true,
								"urls": {
									"compactLayouts": "/services/data/v44.0/sobjects/Account/describe/compactLayouts",
									"rowTemplate": "/services/data/v44.0/sobjects/Account/{ID}",
									"approvalLayouts": "/services/data/v44.0/sobjects/Account/describe/approvalLayouts",
									"defaultValues": "/services/data/v44.0/sobjects/Account/defaultValues?recordTypeId&fields",
									"listviews": "/services/data/v44.0/sobjects/Account/listviews",
									"describe": "/services/data/v44.0/sobjects/Account/describe",
									"quickActions": "/services/data/v44.0/sobjects/Account/quickActions",
									"layouts": "/services/data/v44.0/sobjects/Account/describe/layouts",
									"sobject": "/services/data/v44.0/sobjects/Account"
								}
							},
							"recentItems": []
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
				sobject: "someobject",
			},
			want: MetadataValue{
				ObjectDescribe: ObjectDescribe{
					Creatable:     true,
					Deletable:     true,
					FeedEnabled:   true,
					HasSubtype:    true,
					KeyPrefix:     "001",
					Label:         "Account",
					LabelPlural:   "Accounts",
					Layoutable:    true,
					Mergeable:     true,
					MruEnabled:    true,
					Name:          "Account",
					Queryable:     true,
					Replicateable: true,
					Retrieveable:  true,
					Searchable:    true,
					Triggerable:   true,
					Undeletable:   true,
					Updateable:    true,
					URLs: ObjectURLs{
						CompactLayouts:  "/services/data/v44.0/sobjects/Account/describe/compactLayouts",
						RowTemplate:     "/services/data/v44.0/sobjects/Account/{ID}",
						ApprovalLayouts: "/services/data/v44.0/sobjects/Account/describe/approvalLayouts",
						DefaultValues:   "/services/data/v44.0/sobjects/Account/defaultValues?recordTypeId&fields",
						ListViews:       "/services/data/v44.0/sobjects/Account/listviews",
						Describe:        "/services/data/v44.0/sobjects/Account/describe",
						QuickActions:    "/services/data/v44.0/sobjects/Account/quickActions",
						Layouts:         "/services/data/v44.0/sobjects/Account/describe/layouts",
						SObject:         "/services/data/v44.0/sobjects/Account",
					},
				},
				RecentItems: make([]map[string]interface{}, 0),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := &metadata{
				session: tt.fields.session,
			}
			got, err := md.callout(tt.args.sobject)
			if (err != nil) != tt.wantErr {
				t.Errorf("metadata.Metadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("metadata.Metadata() = %v, want %v", got, tt.want)
			}
		})
	}
}
