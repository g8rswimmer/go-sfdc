package soql

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/g8rswimmer/goforce"
	"github.com/g8rswimmer/goforce/session"
)


type QueryResult struct {
	response queryResponse
	records  []*QueryRecord
}

type QueryRecord struct {
	record     *goforce.Record
	subresults map[string]*QueryResult
}

func newQueryResult(response queryResponse) (*QueryResult, error) {
	result := &QueryResult{
		response: response,
		records:  make([]*QueryRecord, len(response.Records)),
	}

	for idx, record := range response.Records {
		qr, err := newQueryRecord(record)
		if err != nil {
			return nil, err
		}
		result.records[idx] = qr
	}
	return result, nil
}
func (result *QueryResult) Done() bool {
	return result.response.Done
}
func (result *QueryResult) TotalSize() int {
	return result.response.TotalSize
}
func (result *QueryResult) MoreRecords() bool {
	return result.response.NextRecordsURL != ""
}
func (result *QueryResult) Records() []*QueryRecord {
	return result.records
}

func newQueryRecord(jsonMap map[string]interface{}) (*QueryRecord, error) {
	rec, err := goforce.RecordFromJSONMap(jsonMap)
	if err != nil {
		return nil, err
	}
	subresults := make(map[string]*QueryResult)
	for k, v := range jsonMap {
		if sub, has := v.(map[string]interface{}); has {
			if k != goforce.RecordAttributes {
				resp, err := newQueryResponseJSON(sub)
				if err != nil {
					return nil, err
				}
				result, err := newQueryResult(resp)
				if err != nil {
					return nil, err
				}
				subresults[k] = result
			}
		}
	}
	qr := &QueryRecord{
		record:     rec,
		subresults: subresults,
	}
	return qr, nil
}
func (rec *QueryRecord) Record() *goforce.Record {
	return rec.record
}
func (rec *QueryRecord) Subresults(sub string) (*QueryResult, bool) {
	result, has := rec.subresults[sub]
	return result, has
}

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

	result, err := newQueryResult(response)
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
