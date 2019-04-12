package soql

import (
	"reflect"
	"testing"
)

func testNewQueryRecords(records []map[string]interface{}) []*QueryRecord {
	recs := make([]*QueryRecord, len(records))
	for idx, record := range records {
		rec, err := newQueryRecord(record, nil)
		if err != nil {
			return nil
		}
		recs[idx] = rec
	}
	return recs
}
func Test_newQueryResult(t *testing.T) {
	type args struct {
		response queryResponse
	}
	tests := []struct {
		name    string
		args    args
		want    *QueryResult
		wantErr bool
	}{
		{
			name: "No sub results",
			args: args{
				response: queryResponse{
					Done:      true,
					TotalSize: 2,
					Records: []map[string]interface{}{
						{
							"attributes": map[string]interface{}{
								"type": "Account",
								"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
							},
							"Name": "Test 1",
						},
						{
							"attributes": map[string]interface{}{
								"type": "Account",
								"url":  "/services/data/v20.0/sobjects/Account/001D000000IomazIAB",
							},
							"Name": "Test 2",
						},
					},
				},
			},
			want: &QueryResult{
				response: queryResponse{
					Done:      true,
					TotalSize: 2,
					Records: []map[string]interface{}{
						{
							"attributes": map[string]interface{}{
								"type": "Account",
								"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
							},
							"Name": "Test 1",
						},
						{
							"attributes": map[string]interface{}{
								"type": "Account",
								"url":  "/services/data/v20.0/sobjects/Account/001D000000IomazIAB",
							},
							"Name": "Test 2",
						},
					},
				},
				records: testNewQueryRecords([]map[string]interface{}{
					{
						"attributes": map[string]interface{}{
							"type": "Account",
							"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
						},
						"Name": "Test 1",
					},
					{
						"attributes": map[string]interface{}{
							"type": "Account",
							"url":  "/services/data/v20.0/sobjects/Account/001D000000IomazIAB",
						},
						"Name": "Test 2",
					},
				}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newQueryResult(tt.args.response, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("newQueryResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newQueryResult() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryResult_Done(t *testing.T) {
	type fields struct {
		response queryResponse
		records  []*QueryRecord
		resource *Resource
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Done",
			fields: fields{
				response: queryResponse{
					Done: true,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &QueryResult{
				response: tt.fields.response,
				records:  tt.fields.records,
				resource: tt.fields.resource,
			}
			if got := result.Done(); got != tt.want {
				t.Errorf("QueryResult.Done() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryResult_TotalSize(t *testing.T) {
	type fields struct {
		response queryResponse
		records  []*QueryRecord
		resource *Resource
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Total Size",
			fields: fields{
				response: queryResponse{
					TotalSize: 23,
				},
			},
			want: 23,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &QueryResult{
				response: tt.fields.response,
				records:  tt.fields.records,
				resource: tt.fields.resource,
			}
			if got := result.TotalSize(); got != tt.want {
				t.Errorf("QueryResult.TotalSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryResult_MoreRecords(t *testing.T) {
	type fields struct {
		response queryResponse
		records  []*QueryRecord
		resource *Resource
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Has More",
			fields: fields{
				response: queryResponse{
					NextRecordsURL: "The Next URL",
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &QueryResult{
				response: tt.fields.response,
				records:  tt.fields.records,
				resource: tt.fields.resource,
			}
			if got := result.MoreRecords(); got != tt.want {
				t.Errorf("QueryResult.MoreRecords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryResult_Records(t *testing.T) {
	type fields struct {
		response queryResponse
		records  []*QueryRecord
		resource *Resource
	}
	tests := []struct {
		name   string
		fields fields
		want   []*QueryRecord
	}{
		{
			name: "Result Records",
			fields: fields{
				records: testNewQueryRecords([]map[string]interface{}{
					{
						"attributes": map[string]interface{}{
							"type": "Account",
							"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
						},
						"Name": "Test 1",
					},
					{
						"attributes": map[string]interface{}{
							"type": "Account",
							"url":  "/services/data/v20.0/sobjects/Account/001D000000IomazIAB",
						},
						"Name": "Test 2",
					},
				}),
			},
			want: testNewQueryRecords([]map[string]interface{}{
				{
					"attributes": map[string]interface{}{
						"type": "Account",
						"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
					},
					"Name": "Test 1",
				},
				{
					"attributes": map[string]interface{}{
						"type": "Account",
						"url":  "/services/data/v20.0/sobjects/Account/001D000000IomazIAB",
					},
					"Name": "Test 2",
				},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &QueryResult{
				response: tt.fields.response,
				records:  tt.fields.records,
				resource: tt.fields.resource,
			}
			if got := result.Records(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryResult.Records() = %v, want %v", got, tt.want)
			}
		})
	}
}
