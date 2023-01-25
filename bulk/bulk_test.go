package bulk

import (
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/g8rswimmer/go-sfdc/session"
)

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
			name: "Created",
			args: args{
				session: &mockSessionFormatter{},
			},
			want: &Resource{
				session: &mockSessionFormatter{},
			},
			wantErr: false,
		},
		{
			name:    "failed",
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

func TestResource_CreateJob(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
	}
	type args struct {
		options Options
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Passing",
			fields: fields{
				session: &mockSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {
						if req.URL.String() != "https://test.salesforce.com/jobs/ingest" {
							return &http.Response{
								StatusCode: 500,
								Status:     "Invalid URL",
								Body:       io.NopCloser(strings.NewReader(req.URL.String())),
								Header:     make(http.Header),
							}
						}

						resp := `{
							"apiVersion": 44.0,
							"columnDelimiter": "COMMA",
							"concurrencyMode": "Parallel",
							"contentType": "CSV",
							"contentUrl": "services/v44.0/jobs",
							"createdById": "1234",
							"createdDate": "1/1/1970",
							"externalIdFieldName": "namename",
							"id": "9876",
							"jobType": "V2Ingest",
							"lineEnding": "LF",
							"object": "Account",
							"operation": "Insert",
							"state": "Open",
							"systemModstamp": "1/1/1980"
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
				options: Options{
					ColumnDelimiter:     Comma,
					ContentType:         CSV,
					ExternalIDFieldName: "Some External Field",
					LineEnding:          Linefeed,
					Object:              "Account",
					Operation:           Insert,
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
			_, err := r.CreateJob(tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resource.CreateJob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestResource_AllJobs(t *testing.T) {
	mockSession := &mockSessionFormatter{
		url: "https://test.salesforce.com",
		client: mockHTTPClient(func(req *http.Request) *http.Response {
			if req.URL.String() != "https://test.salesforce.com/jobs/ingest?isPkChunkingEnabled=false&jobType=V2Ingest" {
				return &http.Response{
					StatusCode: 500,
					Status:     "Invalid URL",
					Body:       io.NopCloser(strings.NewReader(req.URL.String())),
					Header:     make(http.Header),
				}
			}

			resp := `{
				"done": true,
				"records": [
					{
						"apiVersion": 44.0,
						"columnDelimiter": "COMMA",
						"concurrencyMode": "Parallel",
						"contentType": "CSV",
						"contentUrl": "services/v44.0/jobs",
						"createdById": "1234",
						"createdDate": "1/1/1970",
						"externalIdFieldName": "namename",
						"id": "9876",
						"jobType": "V2Ingest",
						"lineEnding": "LF",
						"object": "Account",
						"operation": "Insert",
						"state": "Open",
						"systemModstamp": "1/1/1980"
					}								
				]
			}`
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "Good",
				Body:       io.NopCloser(strings.NewReader(resp)),
				Header:     make(http.Header),
			}
		}),
	}

	type fields struct {
		session session.ServiceFormatter
	}
	type args struct {
		parameters Parameters
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Jobs
		wantErr bool
	}{
		{
			name: "Passing",
			fields: fields{
				session: mockSession,
			},
			args: args{
				parameters: Parameters{
					JobType: V2Ingest,
				},
			},
			want: &Jobs{
				session: mockSession,
				response: jobResponse{
					Done: true,
					Records: []Response{
						{
							APIVersion:          44.0,
							ColumnDelimiter:     "COMMA",
							ConcurrencyMode:     "Parallel",
							ContentType:         "CSV",
							ContentURL:          "services/v44.0/jobs",
							CreatedByID:         "1234",
							CreatedDate:         "1/1/1970",
							ExternalIDFieldName: "namename",
							ID:                  "9876",
							JobType:             "V2Ingest",
							LineEnding:          "LF",
							Object:              "Account",
							Operation:           "Insert",
							State:               "Open",
							SystemModstamp:      "1/1/1980",
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
			got, err := r.AllJobs(tt.args.parameters)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resource.AllJobs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resource.AllJobs() = %v, want %v", got, tt.want)
			}
		})
	}
}
