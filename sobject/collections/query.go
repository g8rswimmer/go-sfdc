package collections

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/g8rswimmer/goforce"
	"github.com/g8rswimmer/goforce/session"
	"github.com/g8rswimmer/goforce/sobject"
)

type collectionQueryPayload struct {
	IDs    []string `json:"ids"`
	Fields []string `json:"fields"`
}

type Query struct {
	session session.ServiceFormatter
	records []sobject.Querier
	sobject string
}

func (q *Query) Callout() ([]*goforce.Record, error) {
	if q == nil {
		panic("collections: Collection Query can not be nil")
	}
	payload, err := q.payload()
	if err != nil {
		return nil, err
	}
	c := &collection{
		method:   http.MethodPost,
		body:     payload,
		endpoint: endpoint + "/" + q.sobject,
	}
	var values []*goforce.Record
	err = c.send(q.session, &values)
	if err != nil {
		return nil, err
	}
	return values, nil
}
func (q *Query) Records(records ...sobject.Querier) {
	if q == nil {
		panic("collections: Collection Query can not be nil")
	}
	q.records = append(q.records, records...)
}
func (q *Query) payload() (io.Reader, error) {
	fields := make(map[string]interface{})
	ids := make(map[string]interface{})
	for _, querier := range q.records {
		if q.sobject == "" {
			q.sobject = querier.SObject()
		} else {
			if q.sobject != querier.SObject() {
				return nil, fmt.Errorf("sobject collections: sobjects do not match got %s want %s", querier.SObject(), q.sobject)
			}
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
func (q *Query) keyArray(m map[string]interface{}) []string {
	array := make([]string, len(m))
	idx := 0
	for k := range m {
		array[idx] = k
		idx++
	}
	return array
}
