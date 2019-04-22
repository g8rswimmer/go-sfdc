package tree

import (
	"reflect"
	"testing"
)

type mockBuilder struct {
	sobject     string
	fields      map[string]interface{}
	referenceID string
}

func (mock *mockBuilder) SObject() string {
	return mock.sobject
}
func (mock *mockBuilder) Fields() map[string]interface{} {
	return mock.fields
}
func (mock *mockBuilder) ReferenceID() string {
	return mock.referenceID
}

func TestNewRecordBuilder(t *testing.T) {
	type args struct {
		builder Builder
	}
	tests := []struct {
		name    string
		args    args
		want    *RecordBuilder
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				builder: &mockBuilder{
					sobject: "Account",
					fields: map[string]interface{}{
						"name":              "SampleAccount",
						"phone":             "1234567890",
						"website":           "www.salesforce.com",
						"numberOfEmployees": 100,
						"industry":          "Banking",
					},
					referenceID: "ref1",
				},
			},
			want: &RecordBuilder{
				record: Record{
					Attributes: Attributes{
						Type:        "Account",
						ReferenceID: "ref1",
					},
					Fields: map[string]interface{}{
						"name":              "SampleAccount",
						"phone":             "1234567890",
						"website":           "www.salesforce.com",
						"numberOfEmployees": 100,
						"industry":          "Banking",
					},
					Records: make(map[string][]*Record),
				},
			},
			wantErr: false,
		},
		{
			name: "sobject wrong",
			args: args{
				builder: &mockBuilder{
					sobject: "",
					fields: map[string]interface{}{
						"name":              "SampleAccount",
						"phone":             "1234567890",
						"website":           "www.salesforce.com",
						"numberOfEmployees": 100,
						"industry":          "Banking",
					},
					referenceID: "ref1",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no reference id",
			args: args{
				builder: &mockBuilder{
					sobject: "Account",
					fields: map[string]interface{}{
						"name":              "SampleAccount",
						"phone":             "1234567890",
						"website":           "www.salesforce.com",
						"numberOfEmployees": 100,
						"industry":          "Banking",
					},
					referenceID: "",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRecordBuilder(tt.args.builder)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRecordBuilder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRecordBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordBuilder_SubRecords(t *testing.T) {
	type fields struct {
		record Record
	}
	type args struct {
		sobjects string
		records  []*Record
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]int
	}{
		{
			name: "Add sub-records",
			fields: fields{
				record: Record{
					Records: make(map[string][]*Record),
				},
			},
			args: args{
				sobjects: "Contacts",
				records: []*Record{
					{
						Attributes: Attributes{
							Type:        "Contact",
							ReferenceID: "ref2",
						},
						Fields: map[string]interface{}{
							"lastname": "Smith",
							"title":    "President",
							"email":    "sample@salesforce.com",
						},
					},
					{
						Attributes: Attributes{
							Type:        "Contact",
							ReferenceID: "ref3",
						},
						Fields: map[string]interface{}{
							"lastname": "Evans",
							"title":    "Vice President",
							"email":    "sample@salesforce.com",
						},
					},
				},
			},
			want: map[string]int{
				"Contacts": 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := &RecordBuilder{
				record: tt.fields.record,
			}
			rb.SubRecords(tt.args.sobjects, tt.args.records...)
			if len(rb.record.Records) != len(tt.want) {
				t.Errorf("RecordBuilder.SubRecords() = %v, want %v", len(rb.record.Records), len(tt.want))
			}
			for k, v := range tt.want {
				if records, ok := rb.record.Records[k]; ok {
					if len(records) != v {
						t.Errorf("RecordBuilder.SubRecords() = %v, want %v", len(records), v)
					}
				} else {
					t.Errorf("RecordBuilder.SubRecords() want %v", k)
				}
			}
		})
	}
}

func TestRecordBuilder_Build(t *testing.T) {
	type fields struct {
		record Record
	}
	tests := []struct {
		name   string
		fields fields
		want   *Record
	}{
		{
			name: "Success",
			fields: fields{
				record: Record{
					Attributes: Attributes{
						Type:        "Account",
						ReferenceID: "ref1",
					},
					Fields: map[string]interface{}{
						"name":              "SampleAccount",
						"phone":             "1234567890",
						"website":           "www.salesforce.com",
						"numberOfEmployees": 100,
						"industry":          "Banking",
					},
				},
			},
			want: &Record{
				Attributes: Attributes{
					Type:        "Account",
					ReferenceID: "ref1",
				},
				Fields: map[string]interface{}{
					"name":              "SampleAccount",
					"phone":             "1234567890",
					"website":           "www.salesforce.com",
					"numberOfEmployees": 100,
					"industry":          "Banking",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := &RecordBuilder{
				record: tt.fields.record,
			}
			if got := rb.Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RecordBuilder.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}
