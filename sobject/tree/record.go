package tree

import (
	"encoding/json"
	"errors"
)

// Record is the composite tree SObject.
type Record struct {
	Attributes Attributes
	Fields     map[string]interface{}
	Records    map[string][]*Record
}

// Attributes are the attributes of the composite tree.
type Attributes struct {
	Type        string `json:"type"`
	ReferenceID string `json:"referenceId"`
}

// MarshalJSON will create the JSON byte array.
func (r *Record) MarshalJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("record: can't unmarshal to a nil struct")
	}

	rec := r.createMap()

	return json.Marshal(rec)
}
func (r *Record) createMap() map[string]interface{} {
	rec := make(map[string]interface{})
	r.addAttributes(rec)
	r.addFields(rec)
	r.addRecords(rec)
	return rec
}
func (r *Record) addAttributes(jsonMap map[string]interface{}) {
	attributes := map[string]interface{}{
		"type":        r.Attributes.Type,
		"referenceId": r.Attributes.ReferenceID,
	}
	jsonMap["attributes"] = attributes
}
func (r *Record) addFields(jsonMap map[string]interface{}) {
	for k, v := range r.Fields {
		jsonMap[k] = v
	}
}
func (r *Record) addRecords(jsonMap map[string]interface{}) {
	for name, records := range r.Records {
		recs := make([]interface{}, len(records))
		for idx, record := range records {
			rec := make(map[string]interface{})
			record.addAttributes(rec)
			record.addFields(rec)
			record.addRecords(rec)
			recs[idx] = rec
		}
		subRecords := map[string]interface{}{
			"records": recs,
		}
		jsonMap[name] = subRecords
	}
}
