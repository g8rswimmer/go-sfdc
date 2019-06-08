package bulk

import (
	"bufio"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/g8rswimmer/go-sfdc/session"
)

func TestJob_formatOptions(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		info    Response
	}
	type args struct {
		options *Options
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Options
		wantErr bool
	}{
		{
			name:   "passing",
			fields: fields{},
			args: args{
				options: &Options{
					ColumnDelimiter:     Caret,
					ContentType:         CSV,
					ExternalIDFieldName: "Some External Field",
					LineEnding:          Linefeed,
					Object:              "Account",
					Operation:           Insert,
				},
			},
			want: &Options{
				ColumnDelimiter:     Caret,
				ContentType:         CSV,
				ExternalIDFieldName: "Some External Field",
				LineEnding:          Linefeed,
				Object:              "Account",
				Operation:           Insert,
			},
			wantErr: false,
		},
		{
			name:   "defaults",
			fields: fields{},
			args: args{
				options: &Options{
					ExternalIDFieldName: "Some External Field",
					Object:              "Account",
					Operation:           Insert,
				},
			},
			want: &Options{
				ColumnDelimiter:     Comma,
				ContentType:         CSV,
				ExternalIDFieldName: "Some External Field",
				LineEnding:          Linefeed,
				Object:              "Account",
				Operation:           Insert,
			},
			wantErr: false,
		},
		{
			name:   "no object",
			fields: fields{},
			args: args{
				options: &Options{
					ExternalIDFieldName: "Some External Field",
					Operation:           Insert,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "no operation",
			fields: fields{},
			args: args{
				options: &Options{
					ExternalIDFieldName: "Some External Field",
					Object:              "Account",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "no external fields",
			fields: fields{},
			args: args{
				options: &Options{
					Object:    "Account",
					Operation: Upsert,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				session: tt.fields.session,
				info:    tt.fields.info,
			}
			err := j.formatOptions(tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("Job.formatOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(tt.args.options, tt.want) {
				t.Errorf("Job.formatOptions() = %v, want %v", tt.args.options, tt.want)
			}
		})
	}
}

func TestJob_newline(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		info    Response
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Carrage return",
			fields: fields{
				info: Response{
					LineEnding: "CRLF",
				},
			},
			want: "\r\n",
		},
		{
			name: "Line feed",
			fields: fields{
				info: Response{
					LineEnding: "LF",
				},
			},
			want: "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				session: tt.fields.session,
				info:    tt.fields.info,
			}
			if got := j.newline(); got != tt.want {
				t.Errorf("Job.newline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_delimiter(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		info    Response
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "tab",
			fields: fields{
				info: Response{
					ColumnDelimiter: "TAB",
				},
			},
			want: "\t",
		},
		{
			name: "back quote",
			fields: fields{
				info: Response{
					ColumnDelimiter: "BACKQUOTE",
				},
			},
			want: "`",
		},
		{
			name: "caret",
			fields: fields{
				info: Response{
					ColumnDelimiter: "CARET",
				},
			},
			want: "^",
		},
		{
			name: "comma",
			fields: fields{
				info: Response{
					ColumnDelimiter: "COMMA",
				},
			},
			want: ",",
		},
		{
			name: "pipe",
			fields: fields{
				info: Response{
					ColumnDelimiter: "PIPE",
				},
			},
			want: "|",
		},
		{
			name: "semi colon",
			fields: fields{
				info: Response{
					ColumnDelimiter: "SEMICOLON",
				},
			},
			want: ";",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				session: tt.fields.session,
				info:    tt.fields.info,
			}
			if got := j.delimiter(); got != tt.want {
				t.Errorf("Job.delimiter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_record(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		info    Response
	}
	type args struct {
		fields []string
		values []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
	}{
		{
			name:   "make record",
			fields: fields{},
			args: args{
				fields: []string{
					"first",
					"last",
					"DOB",
				},
				values: []string{
					"john",
					"doe",
					"1/1/1970",
				},
			},
			want: map[string]string{
				"first": "john",
				"last":  "doe",
				"DOB":   "1/1/1970",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				session: tt.fields.session,
				info:    tt.fields.info,
			}
			if got := j.record(tt.args.fields, tt.args.values); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.record() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_fields(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		info    Response
	}
	type args struct {
		scanner   *bufio.Scanner
		delimiter string
		offset    int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:   "passing",
			fields: fields{},
			args: args{
				scanner:   bufio.NewScanner(strings.NewReader("sf_id|first|last|DOB")),
				delimiter: "|",
				offset:    1,
			},
			want: []string{
				"first",
				"last",
				"DOB",
			},
			wantErr: false,
		},
		{
			name:   "error",
			fields: fields{},
			args: args{
				scanner:   bufio.NewScanner(strings.NewReader("")),
				delimiter: "|",
				offset:    1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				session: tt.fields.session,
				info:    tt.fields.info,
			}
			got, err := j.fields(tt.args.scanner, tt.args.delimiter, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("Job.fields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.fields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testNewRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "https://test.salesforce.com", nil)
	return req
}
func TestJob_response(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		info    Response
	}
	type args struct {
		request *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Response
		wantErr bool
	}{
		{
			name: "Passing",
			fields: fields{
				session: &mockSessionFormatter{
					client: mockHTTPClient(func(req *http.Request) *http.Response {
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
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}

					}),
				},
			},
			args: args{
				request: testNewRequest(),
			},
			want: Response{
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
			want:    Response{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				session: tt.fields.session,
				info:    tt.fields.info,
			}
			got, err := j.response(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Job.response() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.response() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_createCallout(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		info    Response
	}
	type args struct {
		options Options
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Response
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
								Body:       ioutil.NopCloser(strings.NewReader(req.URL.String())),
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
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
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
			want: Response{
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				session: tt.fields.session,
				info:    tt.fields.info,
			}
			got, err := j.createCallout(tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("Job.createCallout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.createCallout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_create(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		info    Response
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
								Body:       ioutil.NopCloser(strings.NewReader(req.URL.String())),
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
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
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
			j := &Job{
				session: tt.fields.session,
				info:    tt.fields.info,
			}
			if err := j.create(tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("Job.create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJob_setState(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		info    Response
	}
	type args struct {
		state State
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Response
		wantErr bool
	}{
		{
			name: "Passing",
			fields: fields{
				info: Response{
					ID: "1234",
				},
				session: &mockSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {
						if req.URL.String() != "https://test.salesforce.com/jobs/ingest/1234" {
							return &http.Response{
								StatusCode: 500,
								Status:     "Invalid URL",
								Body:       ioutil.NopCloser(strings.NewReader(req.URL.String())),
								Header:     make(http.Header),
							}
						}

						if req.Method != http.MethodPatch {
							return &http.Response{
								StatusCode: 500,
								Status:     "Invalid Method",
								Body:       ioutil.NopCloser(strings.NewReader(req.Method)),
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
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}

					}),
				},
			},
			args: args{
				state: UpdateComplete,
			},
			want: Response{
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
			wantErr: false,
		},
		{
			name: "failing",
			fields: fields{
				info: Response{
					ID: "1234",
				},
				session: &mockSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {
						if req.URL.String() != "https://test.salesforce.com/jobs/ingest/1234" {
							return &http.Response{
								StatusCode: 500,
								Status:     "Invalid URL",
								Body:       ioutil.NopCloser(strings.NewReader(req.URL.String())),
								Header:     make(http.Header),
							}
						}

						if req.Method != http.MethodPatch {
							return &http.Response{
								StatusCode: 500,
								Status:     "Invalid Method",
								Body:       ioutil.NopCloser(strings.NewReader(req.Method)),
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
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				state: UpdateComplete,
			},
			want:    Response{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				session: tt.fields.session,
				info:    tt.fields.info,
			}
			got, err := j.setState(tt.args.state)
			if (err != nil) != tt.wantErr {
				t.Errorf("Job.setState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.setState() = %v, want %v", got, tt.want)
			}
		})
	}
}
