package tree

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/g8rswimmer/go-sfdc"

	"github.com/g8rswimmer/go-sfdc/session"
)

// Inserter is used to define the SObject and it's records for the
// composite tree API.
type Inserter interface {
	SObject() string
	Records() []*Record
}

// InsertValue is the return value for each record.
type InsertValue struct {
	ReferenceID string       `json:"referenceId"`
	ID          string       `json:"id"`
	Errors      []sfdc.Error `json:"errors"`
}

// Value is the return value from the API call.
type Value struct {
	HasErrors bool          `json:"hasErrors"`
	Results   []InsertValue `json:"results"`
}

// Resource is the composite tree API resource.
type Resource struct {
	session session.ServiceFormatter
}

const objectEndpoint = "/composite/tree/"

// NewResource creates a new composite tree resource from the session.
func NewResource(session session.ServiceFormatter) (*Resource, error) {
	if session == nil {
		return nil, errors.New("sobject tree: session can not be nil")
	}
	return &Resource{
		session: session,
	}, nil
}

// Insert will call the composite tree API.
func (r *Resource) Insert(inserter Inserter) (*Value, error) {
	if inserter == nil {
		return nil, errors.New("tree resourse: inserter can not be nil")
	}
	sobject := inserter.SObject()
	matching, err := regexp.MatchString(`\w`, sobject)
	if err != nil {
		return nil, err
	}
	if matching == false {
		return nil, fmt.Errorf("tree resourse: %s is not a valid sobject", sobject)
	}

	return r.callout(inserter)
}
func (r *Resource) callout(inserter Inserter) (*Value, error) {

	request, err := r.request(inserter)

	if err != nil {
		return nil, err
	}

	value, err := r.response(request)

	if err != nil {
		return nil, err
	}

	return &value, nil
}
func (r *Resource) request(inserter Inserter) (*http.Request, error) {

	url := r.session.ServiceURL() + objectEndpoint + inserter.SObject()

	body, err := r.payload(inserter)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url, body)

	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	r.session.AuthorizationHeader(request)
	return request, nil

}
func (r *Resource) payload(inserter Inserter) (*bytes.Reader, error) {
	records := struct {
		Records []*Record `json:"records"`
	}{
		Records: inserter.Records(),
	}
	payload, err := json.Marshal(records)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(payload), nil
}
func (r *Resource) response(request *http.Request) (Value, error) {
	response, err := r.session.Client().Do(request)

	if err != nil {
		return Value{}, err
	}

	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	var value Value
	err = decoder.Decode(&value)
	if err != nil {
		return Value{}, err
	}

	if response.StatusCode != http.StatusCreated {
		return value, fmt.Errorf("insert response err: %d %s", response.StatusCode, response.Status)
	}

	return value, nil
}
