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

type CollectionQuery struct {
	session session.ServiceFormatter
	records []sobject.Querier
	sobject string
}

func (cq *CollectionQuery) Query() ([]*goforce.Record, error) {
	payload, err := cq.payload()
	if err != nil {
		return nil, err
	}
	c := &collection{
		method:   http.MethodPost,
		body:     payload,
		endpoint: cq.session.ServiceURL() + endpoint + "/" + cq.sobject,
	}
	var values []*goforce.Record
	err = c.send(cq.session, &values)
	if err != nil {
		return nil, err
	}
	return values, nil
}
func (cq *CollectionQuery) Records(records ...sobject.Querier) {
	if cq == nil {
		panic("collections: Collection Query can not be nil")
	}
	cq.records = append(cq.records, records...)
}
func (cq *CollectionQuery) payload() (io.Reader, error) {
	fields := make(map[string]interface{})
	ids := make(map[string]interface{})
	for _, querier := range cq.records {
		if cq.sobject == "" {
			cq.sobject = querier.SObject()
		} else {
			if cq.sobject != querier.SObject() {
				return nil, fmt.Errorf("sobject collections: sobjects do not match got %s want %s", querier.SObject(), cq.sobject)
			}
		}
		ids[querier.ID()] = nil
		for _, field := range querier.Fields() {
			fields[field] = nil
		}
	}
	queryPayload := collectionQueryPayload{
		IDs:    cq.keyArray(ids),
		Fields: cq.keyArray(fields),
	}
	payload, err := json.Marshal(queryPayload)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(payload), nil
}
func (cq *CollectionQuery) keyArray(m map[string]interface{}) []string {
	array := make([]string, len(m))
	idx := 0
	for k := range m {
		array[idx] = k
		idx++
	}
	return array
}
