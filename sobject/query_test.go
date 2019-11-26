package sobject

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"go-sfdc"
	"go-sfdc/session"
)

type mockQuery struct {
	sobject string
	id      string
	fields  []string
}

func (mock *mockQuery) SObject() string {
	return mock.sobject
}
func (mock *mockQuery) ID() string {
	return mock.id
}
func (mock *mockQuery) Fields() []string {
	return mock.fields
}

type mockRecord struct {
	sobject string
	url     string
	fields  map[string]interface{}
}

func testNewRecord(data []byte) *sfdc.Record {
	var record sfdc.Record
	err := json.Unmarshal(data, &record)
	if err != nil {
		return &sfdc.Record{}
	}
	return &record
}

func testSalesforceParseTime(salesforceTime string) time.Time {
	date, _ := sfdc.ParseTime(salesforceTime)
	return date
}

func Test_query_Query(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
	}
	type args struct {
		querier Querier
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *sfdc.Record
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
				querier: &mockQuery{
					sobject: "Account",
					fields: []string{
						"Name",
						"Email",
					},
					id: "SomeID",
				},
			},
			want:    nil,
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
				querier: &mockQuery{
					sobject: "Account",
					fields: []string{
						"Name",
						"Email",
					},
					id: "SomeID",
				},
			},
			want:    nil,
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
									"errorCode" : "INVALID_EMAIL_ADDRESS"
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
				querier: &mockQuery{
					sobject: "Account",
					fields: []string{
						"Name",
						"Email",
					},
					id: "SomeID",
				},
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
							StatusCode: http.StatusOK,
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				querier: &mockQuery{
					sobject: "Account",
					fields: []string{
						"Name",
						"Email",
					},
					id: "SomeID",
				},
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
						resp := `
						{
							"AccountNumber" : "CD656092",
							"BillingPostalCode" : "27215"
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
				querier: &mockQuery{
					sobject: "Account",
					fields: []string{
						"AccountNumber",
						"BillingPostalCode",
					},
					id: "SomeID",
				},
			},
			want: testNewRecord([]byte(`
			{
				"AccountNumber" : "CD656092",
				"BillingPostalCode" : "27215"
			}`)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &query{
				session: tt.fields.session,
			}
			got, err := q.callout(tt.args.querier)
			if (err != nil) != tt.wantErr {
				t.Errorf("query.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("query.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockExternalQuery struct {
	sobject  string
	id       string
	fields   []string
	external string
}

func (mock *mockExternalQuery) SObject() string {
	return mock.sobject
}
func (mock *mockExternalQuery) ID() string {
	return mock.id
}
func (mock *mockExternalQuery) Fields() []string {
	return mock.fields
}
func (mock *mockExternalQuery) ExternalField() string {
	return mock.external
}
func Test_query_ExternalQuery(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
	}
	type args struct {
		querier ExternalQuerier
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *sfdc.Record
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
				querier: &mockExternalQuery{
					sobject: "Account",
					fields: []string{
						"Name",
						"Email",
					},
					id:       "SomeID",
					external: "ExternalField",
				},
			},
			want:    nil,
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
				querier: &mockExternalQuery{
					sobject: "Account",
					fields: []string{
						"Name",
						"Email",
					},
					id:       "SomeID",
					external: "ExternalField",
				},
			},
			want:    nil,
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
									"errorCode" : "INVALID_EMAIL_ADDRESS"
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
				querier: &mockExternalQuery{
					sobject: "Account",
					fields: []string{
						"Name",
						"Email",
					},
					id:       "SomeID",
					external: "ExternalField",
				},
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
							StatusCode: http.StatusOK,
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				querier: &mockExternalQuery{
					sobject: "Account",
					fields: []string{
						"Name",
						"Email",
					},
					id:       "SomeID",
					external: "ExternalField",
				},
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
						resp := `
						{
							"AccountNumber" : "CD656092",
							"BillingPostalCode" : "27215"
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
				querier: &mockExternalQuery{
					sobject: "Account",
					fields: []string{
						"Name",
						"Email",
					},
					id:       "SomeID",
					external: "ExternalField",
				},
			},
			want: testNewRecord([]byte(`
			{
				"AccountNumber" : "CD656092",
				"BillingPostalCode" : "27215"
			}`)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &query{
				session: tt.fields.session,
			}
			got, err := q.externalCallout(tt.args.querier)
			if (err != nil) != tt.wantErr {
				t.Errorf("query.ExternalQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("query.ExternalQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_query_DeletedRecords(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
	}
	type args struct {
		sobject   string
		startDate time.Time
		endDate   time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    DeletedRecords
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
				sobject:   "Account",
				startDate: time.Now(),
				endDate:   time.Now().AddDate(0, 0, 7),
			},
			want:    DeletedRecords{},
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
				sobject:   "Account",
				startDate: time.Now(),
				endDate:   time.Now().AddDate(0, 0, 7),
			},
			want:    DeletedRecords{},
			wantErr: true,
		},
		{
			name: "Check URL",
			fields: fields{
				session: &mockSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {

						if strings.HasPrefix(req.URL.String(), "https://test.salesforce.com/sobjects/Account/deleted/?") == false {
							t.Errorf("Urls do not match %s", req.URL.String())
						}
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
				sobject:   "Account",
				startDate: time.Now(),
				endDate:   time.Now().AddDate(0, 0, 7),
			},
			want:    DeletedRecords{},
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
									"errorCode" : "INVALID_EMAIL_ADDRESS"
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
				sobject:   "Account",
				startDate: time.Now(),
				endDate:   time.Now().AddDate(0, 0, 7),
			},
			want:    DeletedRecords{},
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
							StatusCode: http.StatusOK,
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				sobject:   "Account",
				startDate: time.Now(),
				endDate:   time.Now().AddDate(0, 0, 7),
			},
			want:    DeletedRecords{},
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
							"deletedRecords" : 
							[ 
								{ 
									"id" : "a00D0000008pQRAIA2", 
									"deletedDate" : "2013-05-03T15:57:00.000+0000"
								}
							],
							"earliestDateAvailable" : "2013-05-03T15:57:00.000+0000",
							"latestDateCovered" : "2013-05-08T21:20:00.000+0000"
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
				sobject:   "Account",
				startDate: time.Now(),
				endDate:   time.Now().AddDate(0, 0, 7),
			},
			want: DeletedRecords{
				Records: []deletedRecord{
					{
						ID:             "a00D0000008pQRAIA2",
						DeletedDateStr: "2013-05-03T15:57:00.000+0000",
						DeletedDate:    testSalesforceParseTime("2013-05-03T15:57:00.000+0000"),
					},
				},
				EarliestDateStr: "2013-05-03T15:57:00.000+0000",
				EarliestDate:    testSalesforceParseTime("2013-05-03T15:57:00.000+0000"),
				LatestDateStr:   "2013-05-08T21:20:00.000+0000",
				LatestDate:      testSalesforceParseTime("2013-05-08T21:20:00.000+0000"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &query{
				session: tt.fields.session,
			}
			got, err := q.deletedRecordsCallout(tt.args.sobject, tt.args.startDate, tt.args.endDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("query.DeletedRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("query.DeletedRecords() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_query_UpdatedRecords(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
	}
	type args struct {
		sobject   string
		startDate time.Time
		endDate   time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    UpdatedRecords
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
				sobject:   "Account",
				startDate: time.Now(),
				endDate:   time.Now().AddDate(0, 0, 7),
			},
			want:    UpdatedRecords{},
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
				sobject:   "Account",
				startDate: time.Now(),
				endDate:   time.Now().AddDate(0, 0, 7),
			},
			want:    UpdatedRecords{},
			wantErr: true,
		},
		{
			name: "Check URL",
			fields: fields{
				session: &mockSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {

						if strings.HasPrefix(req.URL.String(), "https://test.salesforce.com/sobjects/Account/updated/?") == false {
							t.Errorf("Urls do not match %s", req.URL.String())
						}
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
				sobject:   "Account",
				startDate: time.Now(),
				endDate:   time.Now().AddDate(0, 0, 7),
			},
			want:    UpdatedRecords{},
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
									"errorCode" : "INVALID_EMAIL_ADDRESS"
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
				sobject:   "Account",
				startDate: time.Now(),
				endDate:   time.Now().AddDate(0, 0, 7),
			},
			want:    UpdatedRecords{},
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
							StatusCode: http.StatusOK,
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				sobject:   "Account",
				startDate: time.Now(),
				endDate:   time.Now().AddDate(0, 0, 7),
			},
			want:    UpdatedRecords{},
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
							"ids" : 
							[ 
								"a00D0000008pQR5IAM", 
								"a00D0000008pQRGIA2", 
								"a00D0000008pQRFIA2"
							],
							"latestDateCovered" : "2013-05-08T21:20:00.000+0000" 
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
				sobject:   "Account",
				startDate: time.Now(),
				endDate:   time.Now().AddDate(0, 0, 7),
			},
			want: UpdatedRecords{
				Records: []string{
					"a00D0000008pQR5IAM",
					"a00D0000008pQRGIA2",
					"a00D0000008pQRFIA2",
				},
				LatestDateStr: "2013-05-08T21:20:00.000+0000",
				LatestDate:    testSalesforceParseTime("2013-05-08T21:20:00.000+0000"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &query{
				session: tt.fields.session,
			}
			got, err := q.updatedRecordsCallout(tt.args.sobject, tt.args.startDate, tt.args.endDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("query.UpdatedRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("query.UpdatedRecords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_query_GetContent(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
	}
	type args struct {
		id      string
		content ContentType
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
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
				id:      "001D000000INjVe",
				content: AttachmentType,
			},
			want:    nil,
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
				id:      "001D000000INjVe",
				content: AttachmentType,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Check Attachment URL",
			fields: fields{
				session: &mockSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {

						if strings.HasPrefix(req.URL.String(), "https://test.salesforce.com/sobjects/Attachment/001D000000INjVe/body") == false {
							t.Errorf("Urls do not match %s", req.URL.String())
						}
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
				id:      "001D000000INjVe",
				content: AttachmentType,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Check Attachment URL",
			fields: fields{
				session: &mockSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {

						if strings.HasPrefix(req.URL.String(), "https://test.salesforce.com/sobjects/Document/001D000000INjVe/body") == false {
							t.Errorf("Urls do not match %s", req.URL.String())
						}
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
				id:      "001D000000INjVe",
				content: DocumentType,
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
						resp := "This is the content body"

						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				id:      "001D000000INjVe",
				content: DocumentType,
			},
			want:    []byte("This is the content body"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &query{
				session: tt.fields.session,
			}
			got, err := q.contentCallout(tt.args.id, tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("query.GetContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("query.GetContent() = %v, want %v", got, tt.want)
			}
		})
	}
}
