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

// Resource is the structure for the SObject Collections API.
type Resource struct {
	update *update
	query  *query
	insert *insert
	remove *remove
}

// NewResources forms the Salesforce SObject Collections resource structure.  The
// session formatter is required to form the proper URLs and authorization
// header.
func NewResources(session session.ServiceFormatter) *Resource {
	return &Resource{
		update: &update{
			session: session,
		},
		query: &query{
			session: session,
		},
		insert: &insert{
			session: session,
		},
		remove: &remove{
			session: session,
		},
	}
}

// Insert will create a group of records in the Salesforce org.  The records do not need to be
// the same SObject.  It is the responsibility of the caller to properly chunck the records.
func (r *Resource) Insert(allOrNone bool, records []sobject.Inserter) ([]sobject.InsertValue, error) {
	if r == nil {
		panic("collections resource: Resource can not be nil")
	}
	if r.insert == nil {
		return nil, errors.New("collections resource: collections may not have been initialized properly")
	}
	if records == nil {
		return nil, errors.New("collections resource: insert records can not be nil")
	}
	return r.insert.callout(allOrNone, records)
}

// Delete will remove a group of records in the Salesforce org.  The records do not need to
// be the same SObject.
func (r *Resource) Delete(allOrNone bool, records []string) ([]DeleteValue, error) {
	if r == nil {
		panic("collections resource: Resource can not be nil")
	}
	if r.remove == nil {
		return nil, errors.New("collections resource: collections may not have been initialized properly")
	}
	if records == nil {
		return nil, errors.New("collections resource: delete records can not be nil")
	}
	return r.remove.callout(allOrNone, records)
}

// Update will update a group of records in the Salesforce org.  The records do not need to be
// the same SObject.  It is the responsibility of the caller to properly chunck the records.
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

// Query will retrieve a group of records from the Salesforce org.  The records to retrieve must
// be the same SObject.
func (r *Resource) Query(sobject string, records []sobject.Querier) ([]*goforce.Record, error) {
	if r == nil {
		panic("collections resource: Resource can not be nil")
	}
	if r.query == nil {
		return nil, errors.New("collections resource: collections may not have been initialized properly")
	}
	if records == nil {
		return nil, errors.New("collections resource: update records can not be nil")
	}

	matching, err := regexp.MatchString(`\w`, sobject)
	if err != nil {
		return nil, err
	}

	if matching == false {
		return nil, fmt.Errorf("collection resource: %s is not a valid sobject", sobject)
	}

	return r.query.callout(sobject, records)
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
