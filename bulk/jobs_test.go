package bulk

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/TuSKan/go-sfdc/session"
)

func TestJobs_do(t *testing.T) {
	type fields struct {
		session  session.ServiceFormatter
		response jobResponse
	}
	type args struct {
		request *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    jobResponse
		wantErr bool
	}{
		{
			name: "Passing",
			fields: fields{
				session: &mockSessionFormatter{
					client: mockHTTPClient(func(req *http.Request) *http.Response {
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
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}

					}),
				},
			},
			args: args{
				request: testNewRequest(),
			},
			want: jobResponse{
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
			wantErr: false,
		},
		{
			name: "failing",
			fields: fields{
				session: &mockSessionFormatter{
					client: mockHTTPClient(func(req *http.Request) *http.Response {
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
				request: testNewRequest(),
			},
			want:    jobResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Jobs{
				session:  tt.fields.session,
				response: tt.fields.response,
			}
			got, err := j.do(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Jobs.do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Jobs.do() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newJobs(t *testing.T) {

	mockSession := &mockSessionFormatter{
		url: "https://test.salesforce.com",
		client: mockHTTPClient(func(req *http.Request) *http.Response {
			if req.URL.String() != "https://test.salesforce.com/jobs/ingest?isPkChunkingEnabled=false&jobType=V2Ingest" {
				return &http.Response{
					StatusCode: 500,
					Status:     "Invalid URL",
					Body:       ioutil.NopCloser(strings.NewReader(req.URL.String())),
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
				Body:       ioutil.NopCloser(strings.NewReader(resp)),
				Header:     make(http.Header),
			}
		}),
	}

	type args struct {
		session    session.ServiceFormatter
		parameters Parameters
	}
	tests := []struct {
		name    string
		args    args
		want    *Jobs
		wantErr bool
	}{
		{
			name: "Passing",
			args: args{
				session: mockSession,
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
			got, err := newJobs(tt.args.session, tt.args.parameters)
			if (err != nil) != tt.wantErr {
				t.Errorf("newJobs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newJobs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJobs_Done(t *testing.T) {
	type fields struct {
		session  session.ServiceFormatter
		response jobResponse
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Passing",
			fields: fields{
				response: jobResponse{
					Done: true,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Jobs{
				session:  tt.fields.session,
				response: tt.fields.response,
			}
			if got := j.Done(); got != tt.want {
				t.Errorf("Jobs.Done() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJobs_Records(t *testing.T) {
	type fields struct {
		session  session.ServiceFormatter
		response jobResponse
	}
	tests := []struct {
		name   string
		fields fields
		want   []Response
	}{
		{
			name: "Passing",
			fields: fields{
				response: jobResponse{
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
			want: []Response{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Jobs{
				session:  tt.fields.session,
				response: tt.fields.response,
			}
			if got := j.Records(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Jobs.Records() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJobs_Next(t *testing.T) {
	mockSession := &mockSessionFormatter{
		url: "https://test.salesforce.com",
		client: mockHTTPClient(func(req *http.Request) *http.Response {
			if req.URL.String() != "https://test.salesforce.com/jobs/ingest?isPkChunkingEnabled=false&jobType=V2Ingest" {
				return &http.Response{
					StatusCode: 500,
					Status:     "Invalid URL",
					Body:       ioutil.NopCloser(strings.NewReader(req.URL.String())),
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
				Body:       ioutil.NopCloser(strings.NewReader(resp)),
				Header:     make(http.Header),
			}
		}),
	}

	type fields struct {
		session  session.ServiceFormatter
		response jobResponse
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Jobs
		wantErr bool
	}{
		{
			name: "Passing",
			fields: fields{
				session: mockSession,
				response: jobResponse{
					NextRecordsURL: "https://test.salesforce.com/jobs/ingest?isPkChunkingEnabled=false&jobType=V2Ingest",
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
			j := &Jobs{
				session:  tt.fields.session,
				response: tt.fields.response,
			}
			got, err := j.Next()
			if (err != nil) != tt.wantErr {
				t.Errorf("Jobs.Next() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Jobs.Next() = %v, want %v", got, tt.want)
			}
		})
	}
}
