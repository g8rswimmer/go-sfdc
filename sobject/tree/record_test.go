package tree

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestRecord_addAttributes(t *testing.T) {
	type fields struct {
		Attributes Attributes
		Fields     map[string]interface{}
		Records    map[string][]*Record
	}
	type args struct {
		jsonMap map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]interface{}
	}{
		{
			name: "attributes",
			fields: fields{
				Attributes: Attributes{
					Type:        "Account",
					ReferenceID: "AccountReference",
				},
			},
			args: args{
				jsonMap: make(map[string]interface{}),
			},
			want: map[string]interface{}{
				"attributes": map[string]interface{}{
					"type":        "Account",
					"referenceId": "AccountReference",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Record{
				Attributes: tt.fields.Attributes,
				Fields:     tt.fields.Fields,
				Records:    tt.fields.Records,
			}
			r.addAttributes(tt.args.jsonMap)
			if !reflect.DeepEqual(tt.args.jsonMap, tt.want) {
				t.Errorf("Record.addAttributes() = %v, want %v", tt.args.jsonMap, tt.want)
			}
		})
	}
}

func TestRecord_addFields(t *testing.T) {
	type fields struct {
		Attributes Attributes
		Fields     map[string]interface{}
		Records    map[string][]*Record
	}
	type args struct {
		jsonMap map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]interface{}
	}{
		{
			name: "fields",
			fields: fields{
				Fields: map[string]interface{}{
					"name":              "SampleAccount",
					"phone":             "1234567890",
					"website":           "www.salesforce.com",
					"numberOfEmployees": 100,
					"industry":          "Banking",
				},
			},
			args: args{
				jsonMap: make(map[string]interface{}),
			},
			want: map[string]interface{}{
				"name":              "SampleAccount",
				"phone":             "1234567890",
				"website":           "www.salesforce.com",
				"numberOfEmployees": 100,
				"industry":          "Banking",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Record{
				Attributes: tt.fields.Attributes,
				Fields:     tt.fields.Fields,
				Records:    tt.fields.Records,
			}
			r.addFields(tt.args.jsonMap)
			if !reflect.DeepEqual(tt.args.jsonMap, tt.want) {
				t.Errorf("Record.addFields() = %v, want %v", tt.args.jsonMap, tt.want)
			}
		})
	}
}

func TestRecord_addRecords(t *testing.T) {
	type fields struct {
		Attributes Attributes
		Fields     map[string]interface{}
		Records    map[string][]*Record
	}
	type args struct {
		jsonMap map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]interface{}
	}{
		{
			name: "adding records",
			fields: fields{
				Records: map[string][]*Record{
					"Contacts": []*Record{
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
			},
			args: args{
				jsonMap: make(map[string]interface{}),
			},
			want: map[string]interface{}{
				"Contacts": map[string]interface{}{
					"records": []interface{}{
						map[string]interface{}{
							"attributes": map[string]interface{}{
								"type":        "Contact",
								"referenceId": "ref2",
							},
							"lastname": "Smith",
							"title":    "President",
							"email":    "sample@salesforce.com",
						},
						map[string]interface{}{
							"attributes": map[string]interface{}{
								"type":        "Contact",
								"referenceId": "ref3",
							},
							"lastname": "Evans",
							"title":    "Vice President",
							"email":    "sample@salesforce.com",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Record{
				Attributes: tt.fields.Attributes,
				Fields:     tt.fields.Fields,
				Records:    tt.fields.Records,
			}
			r.addRecords(tt.args.jsonMap)
			if !reflect.DeepEqual(tt.args.jsonMap, tt.want) {
				t.Errorf("Record.addRecords() = %v, want %v", tt.args.jsonMap, tt.want)
			}
		})
	}
}

func TestRecord_createMap(t *testing.T) {
	type fields struct {
		Attributes Attributes
		Fields     map[string]interface{}
		Records    map[string][]*Record
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]interface{}
	}{
		{
			name: "JSON Map",
			fields: fields{
				Attributes: Attributes{
					Type:        "Account",
					ReferenceID: "AccountReference",
				},
				Fields: map[string]interface{}{
					"name":              "SampleAccount",
					"phone":             "1234567890",
					"website":           "www.salesforce.com",
					"numberOfEmployees": 100,
					"industry":          "Banking",
				},
				Records: map[string][]*Record{
					"Contacts": []*Record{
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
			},
			want: map[string]interface{}{
				"attributes": map[string]interface{}{
					"type":        "Account",
					"referenceId": "AccountReference",
				},
				"name":              "SampleAccount",
				"phone":             "1234567890",
				"website":           "www.salesforce.com",
				"numberOfEmployees": 100,
				"industry":          "Banking",
				"Contacts": map[string]interface{}{
					"records": []interface{}{
						map[string]interface{}{
							"attributes": map[string]interface{}{
								"type":        "Contact",
								"referenceId": "ref2",
							},
							"lastname": "Smith",
							"title":    "President",
							"email":    "sample@salesforce.com",
						},
						map[string]interface{}{
							"attributes": map[string]interface{}{
								"type":        "Contact",
								"referenceId": "ref3",
							},
							"lastname": "Evans",
							"title":    "Vice President",
							"email":    "sample@salesforce.com",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Record{
				Attributes: tt.fields.Attributes,
				Fields:     tt.fields.Fields,
				Records:    tt.fields.Records,
			}
			if got := r.createMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Record.createMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecord_MarshalJSON(t *testing.T) {
	type fields struct {
		Attributes Attributes
		Fields     map[string]interface{}
		Records    map[string][]*Record
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "marshaling",
			fields: fields{
				Attributes: Attributes{
					Type:        "Account",
					ReferenceID: "AccountReference",
				},
				Fields: map[string]interface{}{
					"name":              "SampleAccount",
					"phone":             "1234567890",
					"website":           "www.salesforce.com",
					"numberOfEmployees": 100,
					"industry":          "Banking",
				},
				Records: map[string][]*Record{
					"Contacts": []*Record{
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
			},
			want: map[string]interface{}{
				"attributes": map[string]interface{}{
					"type":        "Account",
					"referenceId": "AccountReference",
				},
				"name":              "SampleAccount",
				"phone":             "1234567890",
				"website":           "www.salesforce.com",
				"numberOfEmployees": 100,
				"industry":          "Banking",
				"Contacts": map[string]interface{}{
					"records": []interface{}{
						map[string]interface{}{
							"attributes": map[string]interface{}{
								"type":        "Contact",
								"referenceId": "ref2",
							},
							"lastname": "Smith",
							"title":    "President",
							"email":    "sample@salesforce.com",
						},
						map[string]interface{}{
							"attributes": map[string]interface{}{
								"type":        "Contact",
								"referenceId": "ref3",
							},
							"lastname": "Evans",
							"title":    "Vice President",
							"email":    "sample@salesforce.com",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Record{
				Attributes: tt.fields.Attributes,
				Fields:     tt.fields.Fields,
				Records:    tt.fields.Records,
			}
			got, err := r.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Record.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var jsonMap map[string]interface{}
			err = json.Unmarshal(got, &jsonMap)
			if err != nil {
				t.Errorf("Record.MarshalJSON() error = %v", err)
				return
			}
		})
	}
}
