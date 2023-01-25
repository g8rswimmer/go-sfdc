package sobject

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/g8rswimmer/go-sfdc"
	"github.com/g8rswimmer/go-sfdc/session"
)

// Querier is the interface used to query a SObject from
// Salesforce.
//
// SObject is the table object in Salesforce, like Account.
//
// ID is the Salesforce ID of the table object to retrieve.
//
// Fields is the fields to be returned.  If the field array
// is empty, the all of the fields will be returned.
type Querier interface {
	SObject() string
	ID() string
	Fields() []string
}

// ExternalQuerier is the interface used to query a SObject from
// Salesforce using an external ID.
//
// SObject is the table object in Salesforce, like Account.
//
// ID is the external ID of the table object to retrieve.
//
// Fields is the fields to be returned.  If the field array
// is empty, the all of the fields will be returned.
//
// ExternalField is the external field on the sobject.
type ExternalQuerier interface {
	Querier
	ExternalField() string
}

type deletedRecord struct {
	ID             string    `json:"id"`
	DeletedDateStr string    `json:"deletedDate"`
	DeletedDate    time.Time `json:"-"`
}

// DeletedRecords is the return structure listing the deleted records.
type DeletedRecords struct {
	Records         []deletedRecord `json:"deletedRecords"`
	EarliestDateStr string          `json:"earliestDateAvailable"`
	LatestDateStr   string          `json:"latestDateCovered"`
	EarliestDate    time.Time       `json:"-"`
	LatestDate      time.Time       `json:"-"`
}

// UpdatedRecords is the return structure listing the updated records.
type UpdatedRecords struct {
	Records       []string  `json:"ids"`
	LatestDateStr string    `json:"latestDateCovered"`
	LatestDate    time.Time `json:"-"`
}

// ContentType is indicator of the content type in Salesforce blob.
type ContentType string

const (
	// AttachmentType is the content blob from the Salesforce Attachment record.
	AttachmentType ContentType = "Attachment"
	// DocumentType is the content blob from the Salesforce Document record.
	DocumentType ContentType = "Document"
)

const deletedRoute = "deleted"
const updatedRoute = "updated"
const contentBody = "body"

type query struct {
	session session.ServiceFormatter
}

func (q *query) callout(querier Querier) (*sfdc.Record, error) {
	request, err := q.queryRequest(querier)

	if err != nil {
		return nil, err
	}

	value, err := q.queryResponse(request)

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

func (q *query) queryResponse(request *http.Request) (*sfdc.Record, error) {
	response, err := q.session.Client().Do(request)

	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var queryErrs []sfdc.Error
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

	var record sfdc.Record
	err = decoder.Decode(&record)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (q *query) externalCallout(querier ExternalQuerier) (*sfdc.Record, error) {
	request, err := q.externalQueryRequest(querier)

	if err != nil {
		return nil, err
	}

	value, err := q.queryResponse(request)

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
func (q *query) deletedRecordsCallout(sobject string, startDate, endDate time.Time) (DeletedRecords, error) {
	request, err := q.operationRequest(sobject, deletedRoute, startDate, endDate)

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

	for idx, record := range records.Records {
		date, err := sfdc.ParseTime(record.DeletedDateStr)
		if err != nil {
			return DeletedRecords{}, err
		}
		records.Records[idx].DeletedDate = date
	}
	var date time.Time
	date, err = sfdc.ParseTime(records.EarliestDateStr)
	if err != nil {
		return DeletedRecords{}, err
	}
	records.EarliestDate = date
	date, err = sfdc.ParseTime(records.LatestDateStr)
	if err != nil {
		return DeletedRecords{}, err
	}
	records.LatestDate = date

	return records, nil
}

func (q *query) updatedRecordsCallout(sobject string, startDate, endDate time.Time) (UpdatedRecords, error) {
	request, err := q.operationRequest(sobject, updatedRoute, startDate, endDate)

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

	date, err := sfdc.ParseTime(records.LatestDateStr)
	if err != nil {
		return UpdatedRecords{}, err
	}
	records.LatestDate = date

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

func (q *query) contentCallout(id string, content ContentType) ([]byte, error) {
	request, err := q.contentRequest(id, content)

	if err != nil {
		return nil, err
	}

	return q.contentResponse(request)
}
func (q *query) contentRequest(id string, content ContentType) (*http.Request, error) {

	queryURL := q.session.ServiceURL() + objectEndpoint + string(content) + "/" + id + "/" + contentBody

	request, err := http.NewRequest(http.MethodGet, queryURL, nil)

	if err != nil {
		return nil, err
	}

	q.session.AuthorizationHeader(request)
	return request, nil
}

func (q *query) contentResponse(request *http.Request) ([]byte, error) {
	response, err := q.session.Client().Do(request)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("deleted records response err: %d %s", response.StatusCode, response.Status)
	}

	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	return body, err
}
