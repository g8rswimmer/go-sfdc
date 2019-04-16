package collections

import (
	"io"
	"net/http"

	"github.com/g8rswimmer/goforce/session"
	"github.com/g8rswimmer/goforce/sobject"
)

type Insert struct {
	session session.ServiceFormatter
	records []sobject.Inserter
}

func (i *Insert) Callout(allOrNone bool) ([]sobject.InsertValue, error) {
	payload, err := i.payload(allOrNone)
	if err != nil {
		return nil, err
	}
	c := &collection{
		method:   http.MethodPost,
		body:     payload,
		endpoint: endpoint,
	}
	var values []sobject.InsertValue
	err = c.send(i.session, &values)
	if err != nil {
		return nil, err
	}
	return values, nil
}
func (i *Insert) Records(records ...sobject.Inserter) {
	if i == nil {
		panic("collections: Collection Insert can not be nil")
	}
	i.records = append(i.records, records...)
}
func (i *Insert) payload(allOrNone bool) (io.Reader, error) {
	records := make([]interface{}, len(i.records))
	for idx, inserter := range i.records {
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
