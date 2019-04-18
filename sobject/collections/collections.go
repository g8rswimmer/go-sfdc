package collections

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"

	"github.com/g8rswimmer/goforce"
	"github.com/g8rswimmer/goforce/session"
	"github.com/g8rswimmer/goforce/sobject"
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

type Resource struct {
	session session.ServiceFormatter
	update  *update
}

func NewResource(session session.ServiceFormatter) *Resource {
	return &Resource{
		session: session,
		update: &update{
			session: session,
		},
	}
}

func (r *Resource) NewInsert() *Insert {
	if r == nil {
		panic("collections resource: Resource can not be nil")
	}
	return &Insert{
		session: r.session,
	}
}

func (r *Resource) NewDelete() *Delete {
	if r == nil {
		panic("collections resource: Resource can not be nil")
	}
	return &Delete{
		session: r.session,
	}
}
func (r *Resource) Update(allOrNone bool, records []sobject.Updater) ([]UpdateValue, error) {
	if r == nil {
		panic("collections resource: Resource can not be nil")
	}
	if r.update == nil {
		return nil, errors.New("collections resource: collections may not have been initialized properly")
	}
	if records == nil {
		return nil, errors.New("collections resource: update records can not be nil")
	}
	return r.update.callout(allOrNone, records)
}
func (r *Resource) NewQuery(sobject string) (*Query, error) {
	if r == nil {
		panic("collections resource: Resource can not be nil")
	}

	matching, err := regexp.MatchString(`\w`, sobject)
	if err != nil {
		return nil, err
	}

	if matching == false {
		return nil, fmt.Errorf("collection resource: %s is not a valid sobject", sobject)
	}

	return &Query{
		session: r.session,
		sobject: sobject,
	}, nil
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
				errMsg = fmt.Errorf("insert response err: %s: %s", insertErr.ErrorCode, insertErr.Message)
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
