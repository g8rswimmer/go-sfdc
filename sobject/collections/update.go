package collections

import (
	"io"
	"net/http"

	"github.com/g8rswimmer/goforce/session"
	"github.com/g8rswimmer/goforce/sobject"
)

type UpdateValue struct {
	sobject.InsertValue
}

type Update struct {
	session session.ServiceFormatter
	records []sobject.Updater
}

func (u *Update) Callout(allOrNone bool) ([]UpdateValue, error) {
	payload, err := u.payload(allOrNone)
	if err != nil {
		return nil, err
	}
	c := &collection{
		method:   http.MethodPatch,
		body:     payload,
		endpoint: endpoint,
	}
	var values []UpdateValue
	err = c.send(u.session, &values)
	if err != nil {
		return nil, err
	}
	return values, nil
}
func (u *Update) Records(records ...sobject.Updater) {
	if u == nil {
		panic("collections: Collection Update can not be nil")
	}
	u.records = append(u.records, records...)
}
func (u *Update) payload(allOrNone bool) (io.Reader, error) {
	records := make([]interface{}, len(u.records))
	for idx, updater := range u.records {
		rec := map[string]interface{}{
			"attributes": map[string]string{
				"type": updater.SObject(),
			},
		}
		for field, value := range updater.Fields() {
			rec[field] = value
		}
		rec["id"] = updater.ID()
		records[idx] = rec
	}
	return dmlpayload(allOrNone, records)
}
