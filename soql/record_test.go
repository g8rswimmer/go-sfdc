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
		{
			name: "Sub results",
			args: args{
				jsonMap: map[string]interface{}{
					"attributes": map[string]interface{}{
						"type": "Account",
						"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
					},
					"Name": "Test 1",
					"Contacts": map[string]interface{}{
						"done":      true,
						"totalSize": float64(2),
						"records": []interface{}{
							map[string]interface{}{
								"attributes": map[string]interface{}{
									"type": "Contact",
									"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
								},
								"LastName": "Test 1",
							},
							map[string]interface{}{
								"attributes": map[string]interface{}{
									"type": "Contact",
									"url":  "/services/data/v20.0/sobjects/Account/001D000000IomazIAB",
								},
								"LastName": "Test 2",
							},
						},
					},
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
				subresults: map[string]*QueryResult{
					"Contacts": &QueryResult{
						response: queryResponse{
							Done:      true,
							TotalSize: 2,
							Records: []map[string]interface{}{
								{
									"attributes": map[string]interface{}{
										"type": "Contact",
										"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
									},
									"LastName": "Test 1",
								},
								{
									"attributes": map[string]interface{}{
										"type": "Contact",
										"url":  "/services/data/v20.0/sobjects/Account/001D000000IomazIAB",
									},
									"LastName": "Test 2",
								},
							},
						},
						records: []*QueryRecord{
							{
								record: testQueryRecord(map[string]interface{}{
									"attributes": map[string]interface{}{
										"type": "Contact",
										"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
									},
									"LastName": "Test 1",
								}),
								subresults: make(map[string]*QueryResult),
							},
							{
								record: testQueryRecord(map[string]interface{}{
									"attributes": map[string]interface{}{
										"type": "Contact",
										"url":  "/services/data/v20.0/sobjects/Account/001D000000IomazIAB",
									},
									"LastName": "Test 2",
								}),
								subresults: make(map[string]*QueryResult),
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newQueryRecord(tt.args.jsonMap, nil)
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

func TestQueryRecord_Record(t *testing.T) {
	type fields struct {
		record     *goforce.Record
		subresults map[string]*QueryResult
	}
	tests := []struct {
		name   string
		fields fields
		want   *goforce.Record
	}{
		{
			name: "Get Record",
			fields: fields{
				record: testQueryRecord(map[string]interface{}{
					"attributes": map[string]interface{}{
						"type": "Account",
						"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
					},
					"Name": "Test 1",
				}),
				subresults: make(map[string]*QueryResult),
			},
			want: testQueryRecord(map[string]interface{}{
				"attributes": map[string]interface{}{
					"type": "Account",
					"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
				},
				"Name": "Test 1",
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := &QueryRecord{
				record:     tt.fields.record,
				subresults: tt.fields.subresults,
			}
			if got := rec.Record(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryRecord.Record() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryRecord_Subresults(t *testing.T) {
	type fields struct {
		record     *goforce.Record
		subresults map[string]*QueryResult
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]*QueryResult
	}{
		{
			name: "Sub Results",
			fields: fields{
				record: testQueryRecord(map[string]interface{}{
					"attributes": map[string]interface{}{
						"type": "Account",
						"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
					},
					"Name": "Test 1",
				}),
				subresults: map[string]*QueryResult{
					"Contacts": &QueryResult{
						response: queryResponse{
							Done:      true,
							TotalSize: 2,
							Records: []map[string]interface{}{
								{
									"attributes": map[string]interface{}{
										"type": "Contact",
										"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
									},
									"LastName": "Test 1",
								},
								{
									"attributes": map[string]interface{}{
										"type": "Contact",
										"url":  "/services/data/v20.0/sobjects/Account/001D000000IomazIAB",
									},
									"LastName": "Test 2",
								},
							},
						},
						records: []*QueryRecord{
							{
								record: testQueryRecord(map[string]interface{}{
									"attributes": map[string]interface{}{
										"type": "Contact",
										"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
									},
									"LastName": "Test 1",
								}),
								subresults: make(map[string]*QueryResult),
							},
							{
								record: testQueryRecord(map[string]interface{}{
									"attributes": map[string]interface{}{
										"type": "Contact",
										"url":  "/services/data/v20.0/sobjects/Account/001D000000IomazIAB",
									},
									"LastName": "Test 2",
								}),
								subresults: make(map[string]*QueryResult),
							},
						},
					},
				},
			},
			want: map[string]*QueryResult{
				"Contacts": &QueryResult{
					response: queryResponse{
						Done:      true,
						TotalSize: 2,
						Records: []map[string]interface{}{
							{
								"attributes": map[string]interface{}{
									"type": "Contact",
									"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
								},
								"LastName": "Test 1",
							},
							{
								"attributes": map[string]interface{}{
									"type": "Contact",
									"url":  "/services/data/v20.0/sobjects/Account/001D000000IomazIAB",
								},
								"LastName": "Test 2",
							},
						},
					},
					records: []*QueryRecord{
						{
							record: testQueryRecord(map[string]interface{}{
								"attributes": map[string]interface{}{
									"type": "Contact",
									"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
								},
								"LastName": "Test 1",
							}),
							subresults: make(map[string]*QueryResult),
						},
						{
							record: testQueryRecord(map[string]interface{}{
								"attributes": map[string]interface{}{
									"type": "Contact",
									"url":  "/services/data/v20.0/sobjects/Account/001D000000IomazIAB",
								},
								"LastName": "Test 2",
							}),
							subresults: make(map[string]*QueryResult),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := &QueryRecord{
				record:     tt.fields.record,
				subresults: tt.fields.subresults,
			}
			if got := rec.Subresults(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryRecord.Subresults() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryRecord_Subresult(t *testing.T) {
	type fields struct {
		record     *goforce.Record
		subresults map[string]*QueryResult
	}
	type args struct {
		sub string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *QueryResult
		want1  bool
	}{
		{
			name: "Actual Sub Result",
			fields: fields{
				record: testQueryRecord(map[string]interface{}{
					"attributes": map[string]interface{}{
						"type": "Account",
						"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
					},
					"Name": "Test 1",
				}),
				subresults: map[string]*QueryResult{
					"Contacts": &QueryResult{
						response: queryResponse{
							Done:      true,
							TotalSize: 2,
							Records: []map[string]interface{}{
								{
									"attributes": map[string]interface{}{
										"type": "Contact",
										"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
									},
									"LastName": "Test 1",
								},
								{
									"attributes": map[string]interface{}{
										"type": "Contact",
										"url":  "/services/data/v20.0/sobjects/Account/001D000000IomazIAB",
									},
									"LastName": "Test 2",
								},
							},
						},
						records: []*QueryRecord{
							{
								record: testQueryRecord(map[string]interface{}{
									"attributes": map[string]interface{}{
										"type": "Contact",
										"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
									},
									"LastName": "Test 1",
								}),
								subresults: make(map[string]*QueryResult),
							},
							{
								record: testQueryRecord(map[string]interface{}{
									"attributes": map[string]interface{}{
										"type": "Contact",
										"url":  "/services/data/v20.0/sobjects/Account/001D000000IomazIAB",
									},
									"LastName": "Test 2",
								}),
								subresults: make(map[string]*QueryResult),
							},
						},
					},
				},
			},
			args: args{
				sub: "Contacts",
			},
			want: &QueryResult{
				response: queryResponse{
					Done:      true,
					TotalSize: 2,
					Records: []map[string]interface{}{
						{
							"attributes": map[string]interface{}{
								"type": "Contact",
								"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
							},
							"LastName": "Test 1",
						},
						{
							"attributes": map[string]interface{}{
								"type": "Contact",
								"url":  "/services/data/v20.0/sobjects/Account/001D000000IomazIAB",
							},
							"LastName": "Test 2",
						},
					},
				},
				records: []*QueryRecord{
					{
						record: testQueryRecord(map[string]interface{}{
							"attributes": map[string]interface{}{
								"type": "Contact",
								"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
							},
							"LastName": "Test 1",
						}),
						subresults: make(map[string]*QueryResult),
					},
					{
						record: testQueryRecord(map[string]interface{}{
							"attributes": map[string]interface{}{
								"type": "Contact",
								"url":  "/services/data/v20.0/sobjects/Account/001D000000IomazIAB",
							},
							"LastName": "Test 2",
						}),
						subresults: make(map[string]*QueryResult),
					},
				},
			},
			want1: true,
		},
		{
			name: "Actual Sub Result",
			fields: fields{
				record: testQueryRecord(map[string]interface{}{
					"attributes": map[string]interface{}{
						"type": "Account",
						"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
					},
					"Name": "Test 1",
				}),
				subresults: map[string]*QueryResult{
					"Contacts": &QueryResult{
						response: queryResponse{
							Done:      true,
							TotalSize: 2,
							Records: []map[string]interface{}{
								{
									"attributes": map[string]interface{}{
										"type": "Contact",
										"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
									},
									"LastName": "Test 1",
								},
								{
									"attributes": map[string]interface{}{
										"type": "Contact",
										"url":  "/services/data/v20.0/sobjects/Account/001D000000IomazIAB",
									},
									"LastName": "Test 2",
								},
							},
						},
						records: []*QueryRecord{
							{
								record: testQueryRecord(map[string]interface{}{
									"attributes": map[string]interface{}{
										"type": "Contact",
										"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
									},
									"LastName": "Test 1",
								}),
								subresults: make(map[string]*QueryResult),
							},
							{
								record: testQueryRecord(map[string]interface{}{
									"attributes": map[string]interface{}{
										"type": "Contact",
										"url":  "/services/data/v20.0/sobjects/Account/001D000000IomazIAB",
									},
									"LastName": "Test 2",
								}),
								subresults: make(map[string]*QueryResult),
							},
						},
					},
				},
			},
			args: args{
				sub: "Something",
			},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := &QueryRecord{
				record:     tt.fields.record,
				subresults: tt.fields.subresults,
			}
			got, got1 := rec.Subresult(tt.args.sub)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryRecord.Subresult() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("QueryRecord.Subresult() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
