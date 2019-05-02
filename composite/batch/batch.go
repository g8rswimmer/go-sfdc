package batch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/g8rswimmer/goforce"
	"github.com/g8rswimmer/goforce/session"
)

// Subrequester provides the composite batch API requests.  The
// order of the array is the order in which the subrequests are
// placed in the composite body.
type Subrequester interface {
	URL() string
	Method() string
	BinaryPartName() string
	BinaryPartNameAlias() string
	RichInput() map[string]interface{}
}

// Value is the returned structure from the composite batch API response.
type Value struct {
	HasErrors bool       `json:"hasErrors"`
	Results   []Subvalue `json:"results"`
}

// Subvalue is the subresponses to the composite batch API.  Using the
// referende id, one will be able to match the response with the request.
type Subvalue struct {
	Result     interface{} `json:"result"`
	StatusCode int         `json:"statusCode"`
}

const endpoint = "/composite/batch"

var validMethods = map[string]struct{}{
	"PUT":    {},
	"POST":   {},
	"PATCH":  {},
	"GET":    {},
	"DELETE": {},
}

// Resource is the structure that can be just to call composite batch APIs.
type Resource struct {
	session session.ServiceFormatter
}

// NewResource creates a new resourse with the session.  If the session is
// nil an error will be returned.
func NewResource(session session.ServiceFormatter) (*Resource, error) {
	if session == nil {
		return nil, errors.New("composite: session can not be nil")
	}
	return &Resource{
		session: session,
	}, nil
}

// Retrieve will retrieve the responses to a composite batch requests.
func (r *Resource) Retrieve(haltOnError bool, requesters []Subrequester) (Value, error) {
	if requesters == nil {
		return Value{}, errors.New("composite subrequests: requesters can not nil")
	}
	err := r.validateSubrequests(requesters)
	if err != nil {
		return Value{}, err
	}

	body, err := r.payload(haltOnError, requesters)
	if err != nil {
		return Value{}, err
	}

	url := r.session.ServiceURL() + endpoint

	request, err := http.NewRequest(http.MethodPost, url, body)

	if err != nil {
		return Value{}, err
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	r.session.AuthorizationHeader(request)

	response, err := r.session.Client().Do(request)

	if err != nil {
		return Value{}, err
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

		return Value{}, errMsg
	}

	var value Value
	err = decoder.Decode(&value)
	if err != nil {
		return Value{}, err
	}
	return value, nil
}
func (r *Resource) validateSubrequests(requesters []Subrequester) error {
	for _, requester := range requesters {
		if requester.URL() == "" {
			return errors.New("composite subrequest: must contain an url")
		}
		if _, has := validMethods[requester.Method()]; has == false {
			return errors.New("composite subrequest: empty or invalid method " + requester.Method())
		}
	}
	return nil
}
func (r *Resource) payload(haltOnError bool, requesters []Subrequester) (*bytes.Reader, error) {
	subRequests := make([]interface{}, len(requesters))
	for idx, requester := range requesters {
		subRequest := map[string]interface{}{
			"url":    requester.URL(),
			"method": requester.Method(),
		}
		if requester.BinaryPartName() != "" {
			subRequest["binaryPartName"] = requester.BinaryPartName()
		}
		if requester.BinaryPartNameAlias() != "" {
			subRequest["binaryPartNameAlias"] = requester.BinaryPartNameAlias()
		}
		if requester.RichInput() != nil {
			subRequest["richInput"] = requester.RichInput()
		}
		subRequests[idx] = subRequest
	}
	payload := map[string]interface{}{
		"haltOnError":   haltOnError,
		"batchRequests": subRequests,
	}
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(jsonBody), nil
}
