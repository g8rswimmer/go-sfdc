package soql

import (
	"reflect"
	"testing"

	"github.com/g8rswimmer/goforce"
)

func testQueryRecord(jsonMap map[string]interface{}) *goforce.Record {
	if rec, err := goforce.RecordFromJSONMap(jsonMap); err == nil {
		return rec
	}
	return nil
}
func Test_newQueryRecord(t *testing.T) {
	type args struct {
		jsonMap map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *QueryRecord
		wantErr bool
	}{
		{
			name: "No sub results",
			args: args{
				jsonMap: map[string]interface{}{
					"attributes": map[string]interface{}{
						"type": "Account",
						"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
					},
					"Name": "Test 1",
				},
			},
			want: &QueryRecord{
				record: testQueryRecord(map[string]interface{}{
					"attributes": map[string]interface{}{
						"type": "Account",
						"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
					},
					"Name": "Test 1",
				}),
				subresults: make(map[string]*QueryResult),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newQueryRecord(tt.args.jsonMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("newQueryRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newQueryRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}
