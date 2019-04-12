package soql

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/g8rswimmer/goforce/session"
)

type Resource struct {
	session session.Formatter
}

func NewResource(session session.Formatter) *Resource {
	return &Resource{
		session: session,
	}
}

func (r *Resource) Query(querier Querier, all bool) (*QueryResult, error) {
	if r == nil {
		panic("soql resource: the resource can not be nil")
	}
	if querier == nil {
		return nil, errors.New("soql resource query: querier can not be nil")
	}

	request, err := r.queryRequest(querier, all)
	if err != nil {
		return nil, err
	}

	response, err := r.queryResponse(request)
	if err != nil {
		return nil, err
	}

	result, err := newQueryResult(response, r)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Resource) next(recordURL string) (*QueryResult, error) {
	queryURL := r.session.InstanceURL() + recordURL
	request, err := http.NewRequest(http.MethodGet, queryURL, nil)

	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/json")
	r.session.AuthorizationHeader(request)

	response, err := r.queryResponse(request)
	if err != nil {
		return nil, err
	}

	result, err := newQueryResult(response, r)
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (r *Resource) queryRequest(querier Querier, all bool) (*http.Request, error) {
	query, err := querier.Query()
	if err != nil {
		return nil, err
	}

	endpoint := "/query"
	if all {
		endpoint += "All"
	}

	queryURL := r.session.ServiceURL() + endpoint + "/"

	form := url.Values{}
	form.Add("q", query)
	queryURL += "?" + form.Encode()

	request, err := http.NewRequest(http.MethodGet, queryURL, nil)

	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/json")
	r.session.AuthorizationHeader(request)
	return request, nil

}
func (r *Resource) queryResponse(request *http.Request) (queryResponse, error) {
	response, err := r.session.Client().Do(request)

	if err != nil {
		return queryResponse{}, err
	}

	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var queryErrs []queryError
		err = decoder.Decode(&queryErrs)
		var errMsg error
		if err == nil {
			for _, queryErr := range queryErrs {
				errMsg = fmt.Errorf("insert response err: %s: %s", queryErr.ErrorCode, queryErr.Message)
			}
		} else {
			errMsg = fmt.Errorf("insert response err: %d %s", response.StatusCode, response.Status)
		}

		return queryResponse{}, errMsg
	}

	var resp queryResponse
	err = decoder.Decode(&resp)
	if err != nil {
		return queryResponse{}, err
	}

	return resp, nil
}
