package collections

import (
	"io"
	"net/http"

	"github.com/g8rswimmer/goforce/session"
	"github.com/g8rswimmer/goforce/sobject"
)

type CollectionInsert struct {
	session session.ServiceFormatter
	records []sobject.Inserter
}

func (ci *CollectionInsert) Insert(allOrNone bool) ([]sobject.InsertValue, error) {
	payload, err := ci.payload(allOrNone)
	if err != nil {
		return nil, err
	}
	c := &collection{
		method:   http.MethodPost,
		body:     payload,
		endpoint: ci.session.ServiceURL() + endpoint,
	}
	var values []sobject.InsertValue
	err = c.send(ci.session, &values)
	if err != nil {
		return nil, err
	}
	return values, nil
}
func (ci *CollectionInsert) Records(records ...sobject.Inserter) {
	if ci == nil {
		panic("collections: Collection Insert can not be nil")
	}
	ci.records = append(ci.records, records...)
}
func (ci *CollectionInsert) payload(allOrNone bool) (io.Reader, error) {
	records := make([]interface{}, len(ci.records))
	for idx, inserter := range ci.records {
		rec := map[string]interface{}{
			"attributes": map[string]string{
				"type": inserter.SObject(),
			},
		}
		for field, value := range inserter.Fields() {
			rec[field] = value
		}
		records[idx] = rec
	}
	return dmlpayload(allOrNone, records)
}
