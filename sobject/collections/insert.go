package collections

import (
	"bytes"
	"net/http"

	"github.com/g8rswimmer/go-sfdc/session"
	"github.com/g8rswimmer/go-sfdc/sobject"
)

type insert struct {
	session session.ServiceFormatter
}

func (i *insert) callout(allOrNone bool, records []sobject.Inserter) ([]sobject.InsertValue, error) {
	payload, err := i.payload(allOrNone, records)
	if err != nil {
		return nil, err
	}
	c := &collection{
		method:      http.MethodPost,
		body:        payload,
		endpoint:    endpoint,
		contentType: jsonContentType,
	}
	var values []sobject.InsertValue
	err = c.send(i.session, &values)
	if err != nil {
		return nil, err
	}
	return values, nil
}
func (i *insert) payload(allOrNone bool, records []sobject.Inserter) (*bytes.Reader, error) {
	recs := make([]interface{}, len(records))
	for idx, inserter := range records {
		rec := map[string]interface{}{
			"attributes": map[string]string{
				"type": inserter.SObject(),
			},
		}
		for field, value := range inserter.Fields() {
			rec[field] = value
		}
		recs[idx] = rec
	}
	return dmlpayload(allOrNone, recs)
}
