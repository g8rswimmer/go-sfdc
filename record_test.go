package sfdc

import (
	"reflect"
	"testing"
)

func TestRecord_UnmarshalJSON(t *testing.T) {
	type fields struct {
		sobject string
		url     string
		fields  map[string]interface{}
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Record
		wantErr bool
	}{
		{
			name:   "Successfull Decode",
			fields: fields{},
			args: args{
				data: []byte(`
				{
					"attributes" : {
					  "type" : "Customer__x",
					  "url" : "/services/data/v32.0/sobjects/Customer__x/x01D0000000002RIAQ"
					},
					"Country__c" : "Argentina",
					"Id" : "x01D0000000002RIAQ"
				  }`),
			},
			want: &Record{
				sobject: "Customer__x",
				url:     "/services/data/v32.0/sobjects/Customer__x/x01D0000000002RIAQ",
				fields: map[string]interface{}{
					"Country__c": "Argentina",
					"Id":         "x01D0000000002RIAQ",
				},
			},
			wantErr: false,
		},
		{
			name:   "Successfull Decode too",
			fields: fields{},
			args: args{
				data: []byte(`
				{
					"AccountNumber" : "CD656092",
					"BillingPostalCode" : "27215"
				}`),
			},
			want: &Record{
				fields: map[string]interface{}{
					"AccountNumber":     "CD656092",
					"BillingPostalCode": "27215",
				},
			},
			wantErr: false,
		},
		{
			name:   "Error Decode",
			fields: fields{},
			args: args{
				data: []byte(`
				{
					"AccountNumber" : "CD656092",
					"BillingPostalCode" : "27215",
				}`),
			},
			want:    &Record{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Record{
				sobject: tt.fields.sobject,
				url:     tt.fields.url,
				fields:  tt.fields.fields,
			}
			if err := r.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Record.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(r, tt.want) {
				t.Errorf("Record.UnmarshalJSON() = %v, want %v", r, tt.want)
			}

		})
	}
}

func TestRecord_SObject(t *testing.T) {
	type fields struct {
		sobject string
		url     string
		fields  map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Sobject",
			fields: fields{
				sobject: "Account",
			},
			want: "Account",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Record{
				sobject: tt.fields.sobject,
				url:     tt.fields.url,
				fields:  tt.fields.fields,
			}
			if got := r.SObject(); got != tt.want {
				t.Errorf("Record.SObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecord_URL(t *testing.T) {
	type fields struct {
		sobject string
		url     string
		fields  map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "URL",
			fields: fields{
				url: "/services/data/v32.0/sobjects/Customer__x/CACTU",
			},
			want: "/services/data/v32.0/sobjects/Customer__x/CACTU",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Record{
				sobject: tt.fields.sobject,
				url:     tt.fields.url,
				fields:  tt.fields.fields,
			}
			if got := r.URL(); got != tt.want {
				t.Errorf("Record.URL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecord_FieldValue(t *testing.T) {
	type fields struct {
		sobject string
		url     string
		fields  map[string]interface{}
	}
	type args struct {
		field string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
		want1  bool
	}{
		{
			name: "Field Contains",
			fields: fields{
				fields: map[string]interface{}{
					"Country__c": "Argentina",
					"Id":         "x01D0000000002RIAQ",
				},
			},
			args: args{
				field: "Country__c",
			},
			want:  "Argentina",
			want1: true,
		},
		{
			name: "Field Missing",
			fields: fields{
				fields: map[string]interface{}{
					"Country__c": "Argentina",
					"Id":         "x01D0000000002RIAQ",
				},
			},
			args: args{
				field: "Nope",
			},
			want:  nil,
			want1: false,
		},
		{
			name: "Nested Field",
			fields: fields{
				fields: map[string]interface{}{
					"Id": "x01D0000000002RIAQ",
					"Account": map[string]interface{}{
						"Name": "Test Account",
						"Id":   "x01D0000000002RIAQ",
					},
				},
			},
			args: args{
				field: "Account.Name",
			},
			want:  "Test Account",
			want1: true,
		},
		{
			name: "Nested Field 2",
			fields: fields{
				fields: map[string]interface{}{
					"Id": "x01D0000000002RIAQ",
					"Account": map[string]interface{}{
						"Name": "Test Account",
						"Id":   "x01D0000000002RIAQ",
					},
				},
			},
			args: args{
				field: "Account",
			},
			want: map[string]interface{}{
				"Name": "Test Account",
				"Id":   "x01D0000000002RIAQ",
			},
			want1: true,
		},
		{
			name: "Nested Field Empty",
			fields: fields{
				fields: map[string]interface{}{
					"Id": "x01D0000000002RIAQ",
				},
			},
			args: args{
				field: "Account",
			},
			want:  nil,
			want1: false,
		},
		{
			name: "Not Nested",
			fields: fields{
				fields: map[string]interface{}{
					"Country__c": "Argentina",
					"Id":         "x01D0000000002RIAQ",
				},
			},
			args: args{
				field: "Country__c.Robots__c",
			},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Record{
				sobject: tt.fields.sobject,
				url:     tt.fields.url,
				fields:  tt.fields.fields,
			}
			got, got1 := r.FieldValue(tt.args.field)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Record.FieldValue() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Record.FieldValue() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestRecord_Fields(t *testing.T) {
	type fields struct {
		sobject string
		url     string
		fields  map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]interface{}
	}{
		{
			name: "Get those fields",
			fields: fields{
				fields: map[string]interface{}{
					"Country__c": "Argentina",
					"Id":         "x01D0000000002RIAQ",
				},
			},
			want: map[string]interface{}{
				"Country__c": "Argentina",
				"Id":         "x01D0000000002RIAQ",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Record{
				sobject: tt.fields.sobject,
				url:     tt.fields.url,
				fields:  tt.fields.fields,
			}
			if got := r.Fields(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Record.Fields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordFromJSONMap(t *testing.T) {
	type args struct {
		jsonMap map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *Record
		wantErr bool
	}{
		{
			name: "JSON map",
			args: args{
				jsonMap: map[string]interface{}{
					"attributes": map[string]interface{}{
						"type": "Account",
						"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
					},
					"Name": "Test 1",
					"Contacts": map[string]interface{}{
						"done":      "true",
						"totalSize": 14,
						"records": []map[string]interface{}{
							{
								"attributes": map[string]interface{}{
									"type": "Contact",
									"url":  "/services/data/v20.0/sobjects/Contact/001D000000IRFmaIAH",
								},
								"LastName": "Test 1",
							},
						},
					},
				},
			},
			want: &Record{
				sobject: "Account",
				url:     "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
				fields: map[string]interface{}{
					"Name": "Test 1",
				},
			},
			wantErr: false,
		},{
			name: "JSON map with nested",
			args: args{
				jsonMap: map[string]interface{}{
					"attributes": map[string]interface{}{
						"type": "Account",
						"url":  "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
					},
					"Name": "Test 1",
					"Account": map[string]interface{}{
						"Name": "Test Account",
						"Id":   "x01D0000000002RIAQ",
					},
					"Contacts": map[string]interface{}{
						"done":      "true",
						"totalSize": 14,
						"records": []map[string]interface{}{
							{
								"attributes": map[string]interface{}{
									"type": "Contact",
									"url":  "/services/data/v20.0/sobjects/Contact/001D000000IRFmaIAH",
								},
								"LastName": "Test 1",
							},
						},
					},
				},
			},
			want: &Record{
				sobject: "Account",
				url:     "/services/data/v20.0/sobjects/Account/001D000000IRFmaIAH",
				fields: map[string]interface{}{
					"Name": "Test 1",
					"Account": map[string]interface{}{
						"Name": "Test Account",
						"Id":   "x01D0000000002RIAQ",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "JSON Map error",
			args: args{
				jsonMap: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RecordFromJSONMap(tt.args.jsonMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordFromJSONMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RecordFromJSONMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
