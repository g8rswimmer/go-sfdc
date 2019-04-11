package soql

import (
	"reflect"
	"testing"
)

func testNewQueryRecords(records []map[string]interface{}) []*QueryRecord {
	recs := make([]*QueryRecord, len(records))
	for idx, record := range records {
		rec, err := newQueryRecord(record)
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
			got, err := newQueryResult(tt.args.response)
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
