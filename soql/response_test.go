package soql

import (
	"reflect"
	"testing"
)

func Test_newQueryResponseJSON(t *testing.T) {
	type args struct {
		jsonMap map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    queryResponse
		wantErr bool
	}{
		{
			name: "Passing Decode",
			args: args{
				jsonMap: map[string]interface{}{
					"done":      true,
					"totalSize": float64(2),
					"records": []map[string]interface{}{
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
			want: queryResponse{
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
			wantErr: false,
		},
		{
			name: "Passing Decode with next",
			args: args{
				jsonMap: map[string]interface{}{
					"done":           true,
					"totalSize":      float64(2),
					"nextRecordsUrl": "/services/data/v20.0/query/01gD0000002HU6KIAW-2000",
					"records": []map[string]interface{}{
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
			want: queryResponse{
				Done:           true,
				TotalSize:      2,
				NextRecordsURL: "/services/data/v20.0/query/01gD0000002HU6KIAW-2000",
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
			wantErr: false,
		},
		{
			name: "Passing Decode with done not bool",
			args: args{
				jsonMap: map[string]interface{}{
					"done": "true",
				},
			},
			want:    queryResponse{},
			wantErr: true,
		},
		{
			name: "Passing Decode with done not present",
			args: args{
				jsonMap: map[string]interface{}{},
			},
			want:    queryResponse{},
			wantErr: true,
		},
		{
			name: "Passing Decode with totalSize not int",
			args: args{
				jsonMap: map[string]interface{}{
					"done":      true,
					"totalSize": "float64(2)",
				},
			},
			want:    queryResponse{},
			wantErr: true,
		},
		{
			name: "Passing Decode with totalSize not present",
			args: args{
				jsonMap: map[string]interface{}{
					"done": true,
				},
			},
			want:    queryResponse{},
			wantErr: true,
		},
		{
			name: "Passing Decode with nextRecordsUrl not string",
			args: args{
				jsonMap: map[string]interface{}{
					"done":           true,
					"totalSize":      float64(2),
					"nextRecordsUrl": 22,
				},
			},
			want:    queryResponse{},
			wantErr: true,
		},
		{
			name: "Passing Decode with records not an array",
			args: args{
				jsonMap: map[string]interface{}{
					"done":           true,
					"totalSize":      float64(2),
					"nextRecordsUrl": "/services/data/v20.0/query/01gD0000002HU6KIAW-2000",
					"records":        "something",
				},
			},
			want:    queryResponse{},
			wantErr: true,
		},
		{
			name: "Passing Decode with totalSize not present",
			args: args{
				jsonMap: map[string]interface{}{
					"done":           true,
					"totalSize":      float64(2),
					"nextRecordsUrl": "/services/data/v20.0/query/01gD0000002HU6KIAW-2000",
				},
			},
			want:    queryResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newQueryResponseJSON(tt.args.jsonMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("newQueryResponseJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newQueryResponseJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
