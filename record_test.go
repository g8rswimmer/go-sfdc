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
				lookUps: map[string]*Record{},
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
				lookUps: map[string]*Record{},
			},
			wantErr: false,
		},
		{
			name:   "Successfull Decode with look up",
			fields: fields{},
			args: args{
				data: []byte(`
				{
					"attributes" : {
					  "type" : "Customer__x",
					  "url" : "/services/data/v32.0/sobjects/Customer__x/x01D0000000002RIAQ"
					},
					"Country__c" : "Argentina",
					"Id" : "x01D0000000002RIAQ",
					"SomeLookup__r": {
						"attributes": {
							"type": "SomeLookup__c",
							"url": "/services/data/v44.0/sobjects/SomeLookup__c/0012E00001q0KijQAE"
						},
						"Name": "Salesforce"
					}
				  }`),
			},
			want: &Record{
				sobject: "Customer__x",
				url:     "/services/data/v32.0/sobjects/Customer__x/x01D0000000002RIAQ",
				fields: map[string]interface{}{
					"Country__c": "Argentina",
					"Id":         "x01D0000000002RIAQ",
				},
				lookUps: map[string]*Record{
					"SomeLookup__r": &Record{
						sobject: "SomeLookup__c",
						url:     "/services/data/v44.0/sobjects/SomeLookup__c/0012E00001q0KijQAE",
						fields: map[string]interface{}{
							"Name": "Salesforce",
						},
						lookUps: map[string]*Record{},
					},
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
				lookUps: map[string]*Record{},
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

func TestRecord_LookUp(t *testing.T) {
	type fields struct {
		sobject string
		url     string
		fields  map[string]interface{}
		lookUps map[string]*Record
	}
	type args struct {
		lookUp string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Record
		want1  bool
	}{
		{
			name: "has lookup",
			fields: fields{
				lookUps: map[string]*Record{
					"SomeObject": &Record{
						sobject: "SomeObject",
						fields: map[string]interface{}{
							"Here": true,
						},
					},
				},
			},
			args: args{
				lookUp: "SomeObject",
			},
			want: &Record{
				sobject: "SomeObject",
				fields: map[string]interface{}{
					"Here": true,
				},
			},
			want1: true,
		},
		{
			name: "nope",
			fields: fields{
				lookUps: map[string]*Record{
					"SomeObject": &Record{
						sobject: "SomeObject",
						fields: map[string]interface{}{
							"Here": true,
						},
					},
				},
			},
			args: args{
				lookUp: "SomeOther",
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
				lookUps: tt.fields.lookUps,
			}
			got, got1 := r.LookUp(tt.args.lookUp)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Record.LookUp() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Record.LookUp() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestRecord_LookUps(t *testing.T) {
	type fields struct {
		sobject string
		url     string
		fields  map[string]interface{}
		lookUps map[string]*Record
	}
	tests := []struct {
		name   string
		fields fields
		want   []*Record
	}{
		{
			name: "the look ups",
			fields: fields{
				lookUps: map[string]*Record{
					"LookUps": &Record{
						sobject: "LookUps",
					},
				},
			},
			want: []*Record{
				&Record{
					sobject: "LookUps",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Record{
				sobject: tt.fields.sobject,
				url:     tt.fields.url,
				fields:  tt.fields.fields,
				lookUps: tt.fields.lookUps,
			}
			if got := r.LookUps(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Record.LookUps() = %v, want %v", got, tt.want)
			}
		})
	}
}
