package collections

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"go-sfdc"
	"go-sfdc/session"
	"go-sfdc/sobject"
)

type collectionQueryPayload struct {
	IDs    []string `json:"ids"`
	Fields []string `json:"fields"`
}

type query struct {
	session session.ServiceFormatter
}

func (q *query) callout(sobject string, records []sobject.Querier) ([]*sfdc.Record, error) {
	if q == nil {
		panic("collections: Collection Query can not be nil")
	}
	payload, err := q.payload(sobject, records)
	if err != nil {
		return nil, err
	}
	c := &collection{
		method:      http.MethodPost,
		body:        payload,
		endpoint:    endpoint + "/" + sobject,
		contentType: jsonContentType,
	}
	var values []*sfdc.Record
	err = c.send(q.session, &values)
	if err != nil {
		return nil, err
	}
	return values, nil
}
func (q *query) payload(sobject string, records []sobject.Querier) (*bytes.Reader, error) {
	fields := make(map[string]interface{})
	ids := make(map[string]interface{})
	for _, querier := range records {
		if sobject != querier.SObject() {
			return nil, fmt.Errorf("sobject collections: sobjects do not match got %s want %s", querier.SObject(), sobject)
		}
		ids[querier.ID()] = nil
		for _, field := range querier.Fields() {
			fields[field] = nil
		}
	}
	queryPayload := collectionQueryPayload{
		IDs:    q.keyArray(ids),
		Fields: q.keyArray(fields),
	}
	payload, err := json.Marshal(queryPayload)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(payload), nil
}
func (q *query) keyArray(m map[string]interface{}) []string {
	array := make([]string, len(m))
	idx := 0
	for k := range m {
		array[idx] = k
		idx++
	}
	return array
}
