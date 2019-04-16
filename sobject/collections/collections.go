package collections

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/g8rswimmer/goforce"
	"github.com/g8rswimmer/goforce/session"
)

const endpoint = "/composite/sobjects"

type collectionDmlPayload struct {
	AllOrNone bool          `json:"allOrNone"`
	Records   []interface{} `json:"records"`
}

type collection struct {
	method   string
	endpoint string
	values   *url.Values
	body     io.Reader
}

func (c *collection) send(session session.ServiceFormatter, value interface{}) error {
	collectionURL := session.ServiceURL() + c.endpoint
	if c.values != nil {
		collectionURL += "?" + c.values.Encode()
	}
	request, err := http.NewRequest(c.method, collectionURL, c.body)
	if err != nil {
		return err
	}
	response, err := session.Client().Do(request)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var insertErrs []goforce.Error
		err = decoder.Decode(&insertErrs)
		var errMsg error
		if err == nil {
			for _, insertErr := range insertErrs {
				errMsg = fmt.Errorf("insert response err: %s: %s", insertErr.StatusCode, insertErr.Message)
			}
		} else {
			errMsg = fmt.Errorf("insert response err: %d %s", response.StatusCode, response.Status)
		}
		return errMsg
	}
	err = decoder.Decode(value)
	if err != nil {
		return err
	}
	return nil
}

func dmlpayload(allOrNone bool, records []interface{}) (io.Reader, error) {
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
