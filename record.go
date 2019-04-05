package goforce

import (
	"encoding/json"
)

const recordAttributes = "attributes"
const recordAttrType = "type"
const recordAttrURL = "url"

// Record is a representation of a Salesforce
// record.
type Record struct {
	sobject string
	url     string
	fields  map[string]interface{}
}

func (r *Record) UnmarshalJSON(data []byte) error {
	var jsonMap map[string]interface{}
	err := json.Unmarshal(data, &jsonMap)
	if err != nil {
		return err
	}

	r.fields = make(map[string]interface{})

	for k, v := range jsonMap {
		if k == recordAttributes {
			if attr, ok := v.(map[string]interface{}); ok {
				if obj, ok := attr[recordAttrType]; ok {
					if sobj, ok := obj.(string); ok {
						r.sobject = sobj
					}
				}
				if obj, ok := attr[recordAttrURL]; ok {
					if url, ok := obj.(string); ok {
						r.url = url
					}
				}
			}
		} else {
			r.fields[k] = v
		}
	}

	return nil
}

func (r *Record) SObject() string {
	return r.sobject
}

func (r *Record) URL() string {
	return r.url
}

func (r *Record) FieldValue(field string) (interface{}, bool) {
	value, has := r.fields[field]
	return value, has
}

func (r *Record) Fields() map[string]interface{} {
	fields := make(map[string]interface{})
	for k, v := range r.fields {
		fields[k] = v
	}
	return fields
}
