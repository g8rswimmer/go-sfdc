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

type UpdateValue struct {
	sobject.InsertValue
}

type CollectionUpdate struct {
	session session.ServiceFormatter
	records []sobject.Updater
}

func (cu *CollectionUpdate) Update(allOrNone bool) ([]UpdateValue, error) {
	payload, err := cu.payload(allOrNone)
	if err != nil {
		return nil, err
	}
	c := &collection{
		method:   http.MethodPatch,
		body:     payload,
		endpoint: cu.session.ServiceURL() + endpoint,
	}
	response, err := c.send(cu.session)
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
	var values []UpdateValue
	err = decoder.Decode(&values)
	if err != nil {
		return nil, err
	}
	return values, nil
}
func (cu *CollectionUpdate) Records(records ...sobject.Updater) {
	if cu == nil {
		panic("collections: Collection Update can not be nil")
	}
	cu.records = append(cu.records, records...)
}
func (cu *CollectionUpdate) payload(allOrNone bool) (io.Reader, error) {
	records := make([]interface{}, len(cu.records))
	for idx, updater := range cu.records {
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
