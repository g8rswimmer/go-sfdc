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

type collectionDmlPayload struct {
	AllOrNone bool          `json:"allOrNone"`
	Records   []interface{} `json:"records"`
}

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
	response, err := c.send(ci.session)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var insertErrs []goforce.Error
		err = decoder.Decode(&insertErrs)
		var errMsg error
		if err == nil {
			for _, insertErr := range insertErrs {
				errMsg = fmt.Errorf("insert response err: %s: %s", insertErr.ErrorCode, insertErr.Message)
			}
		} else {
			errMsg = fmt.Errorf("insert response err: %d %s", response.StatusCode, response.Status)
		}
		return nil, errMsg
	}
	var values []sobject.InsertValue
	err = decoder.Decode(&values)
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
	dmlPayload := collectionDmlPayload{
		AllOrNone: allOrNone,
		Records:   records,
	}
	payload, err := json.Marshal(dmlPayload)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(payload), nil
}
