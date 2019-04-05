package sobject

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/g8rswimmer/goforce"
	"github.com/g8rswimmer/goforce/session"
)

type Querier interface {
	SObject() string
	ID() string
	Fields() []string
}

type ExternalQuerier interface {
	Querier
	ExternalField() string
}

type queryError struct {
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
}

type deletedRecord struct {
	ID          string    `json:"id"`
	DeletedDate time.Time `json:"deletedDate"`
}

type DeletedRecords struct {
	Records      []deletedRecord `json:"deletedRecords"`
	EarliestDate time.Time       `json:"earliestDateAvailable"`
	LatestDate   time.Time       `json:"latestDateCovered"`
}

type UpdatedRecords struct {
	Records    []string  `json:"ids"`
	LatestDate time.Time `json:"latestDateCovered"`
}

const deletedRecords = "deleted"
const updatedRecords = "updated"

type query struct {
	session session.Formatter
}

func (q *query) Query(querier Querier) (*goforce.Record, error) {
	request, err := q.queryRequest(querier)

	if err != nil {
		return nil, err
	}

	value, err := q.response(request)

	if err != nil {
		return nil, err
	}

	return value, nil
}
func (q *query) queryRequest(querier Querier) (*http.Request, error) {

	queryURL := q.session.ServiceURL() + objectEndpoint + querier.SObject() + "/" + querier.ID()

	if len(querier.Fields()) > 0 {
		fields := strings.Join(querier.Fields(), ",")
		form := url.Values{}
		form.Add("fields", fields)
		queryURL += "?" + form.Encode()
	}

	request, err := http.NewRequest(http.MethodGet, queryURL, nil)

	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/json")
	q.session.AuthorizationHeader(request)
	return request, nil

}

func (q *query) response(request *http.Request) (*goforce.Record, error) {
	response, err := q.session.Client().Do(request)

	if err != nil {
		return nil, err
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

		return nil, errMsg
	}

	var record goforce.Record
	err = decoder.Decode(&record)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (q *query) ExternalQuery(querier ExternalQuerier) (*goforce.Record, error) {
	request, err := q.externalQueryRequest(querier)

	if err != nil {
		return nil, err
	}

	value, err := q.response(request)

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (q *query) externalQueryRequest(querier ExternalQuerier) (*http.Request, error) {

	queryURL := q.session.ServiceURL() + objectEndpoint + querier.SObject() + "/" + querier.ExternalField() + "/" + querier.ID()

	if len(querier.Fields()) > 0 {
		fields := strings.Join(querier.Fields(), ",")
		form := url.Values{}
		form.Add("fields", fields)
		queryURL += "?" + form.Encode()
	}

	request, err := http.NewRequest(http.MethodGet, queryURL, nil)

	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/json")
	q.session.AuthorizationHeader(request)
	return request, nil

}
func (q *query) DeletedRecords(sobject string, startDate, endDate time.Time) (DeletedRecords, error) {
	request, err := q.operationRequest(sobject, deletedRecords, startDate, endDate)

	if err != nil {
		return DeletedRecords{}, err
	}

	value, err := q.deletedRecordsResponse(request)

	if err != nil {
		return DeletedRecords{}, err
	}

	return value, nil
}

func (q *query) deletedRecordsResponse(request *http.Request) (DeletedRecords, error) {
	response, err := q.session.Client().Do(request)

	if err != nil {
		return DeletedRecords{}, err
	}

	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return DeletedRecords{}, fmt.Errorf("deleted records response err: %d %s", response.StatusCode, response.Status)
	}

	var records DeletedRecords
	err = decoder.Decode(&records)
	if err != nil {
		return DeletedRecords{}, err
	}

	return records, nil
}

func (q *query) UpdatedRecords(sobject string, startDate, endDate time.Time) (UpdatedRecords, error) {
	request, err := q.operationRequest(sobject, updatedRecords, startDate, endDate)

	if err != nil {
		return UpdatedRecords{}, err
	}

	value, err := q.updatedRecordsResponse(request)

	if err != nil {
		return UpdatedRecords{}, err
	}

	return value, nil
}

func (q *query) updatedRecordsResponse(request *http.Request) (UpdatedRecords, error) {
	response, err := q.session.Client().Do(request)

	if err != nil {
		return UpdatedRecords{}, err
	}

	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return UpdatedRecords{}, fmt.Errorf("deleted records response err: %d %s", response.StatusCode, response.Status)
	}

	var records UpdatedRecords
	err = decoder.Decode(&records)
	if err != nil {
		return UpdatedRecords{}, err
	}

	return records, nil
}

func (q *query) operationRequest(sobject, operation string, startDate, endDate time.Time) (*http.Request, error) {

	form := url.Values{}
	form.Add("start", startDate.Format(time.RFC3339))
	form.Add("end", endDate.Format(time.RFC3339))
	dateRange := "?" + form.Encode()

	queryURL := q.session.ServiceURL() + objectEndpoint + sobject + "/" + operation + "/" + dateRange

	request, err := http.NewRequest(http.MethodGet, queryURL, nil)

	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/json")
	q.session.AuthorizationHeader(request)
	return request, nil

}
