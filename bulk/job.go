package bulk

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	sfdc "github.com/g8rswimmer/go-sfdc"
	"github.com/g8rswimmer/go-sfdc/session"
)

// JobType is the bulk job type.
type JobType string

const (
	// BigObjects is the big objects job.
	BigObjects JobType = "BigObjectIngest"
	// Classic is the bulk job 1.0.
	Classic JobType = "Classic"
	// V2Ingest is the bulk job 2.0.
	V2Ingest JobType = "V2Ingest"
	// V2Query is the Query bulk job 2.0.
	V2Query JobType = "V2Query"
)

// ColumnDelimiter is the column delimiter used for CSV job data.
type ColumnDelimiter string

const (
	// Backquote is the (`) character.
	Backquote ColumnDelimiter = "BACKQUOTE"
	// Caret is the (^) character.
	Caret ColumnDelimiter = "CARET"
	// Comma is the (,) character.
	Comma ColumnDelimiter = "COMMA"
	// Pipe is the (|) character.
	Pipe ColumnDelimiter = "PIPE"
	// SemiColon is the (;) character.
	SemiColon ColumnDelimiter = "SEMICOLON"
	// Tab is the (\t) character.
	Tab ColumnDelimiter = "TAB"
)

// ContentType is the format of the data being processed.
type ContentType string

// CSV is the supported content data type.
const CSV ContentType = "CSV"

// LineEnding is the line ending used for the CSV job data.
type LineEnding string

const (
	// Linefeed is the (\n) character.
	Linefeed LineEnding = "LF"
	// CarriageReturnLinefeed is the (\r\n) character.
	CarriageReturnLinefeed LineEnding = "CRLF"
)

// Operation is the processing operation for the job.
type Operation string

const (
	// Insert is the object operation for inserting records.
	Insert Operation = "insert"
	// Delete is the object operation for deleting records.
	Delete Operation = "delete"
	// Update is the object operation for updating records.
	Update Operation = "update"
	// Upsert is the object operation for upserting records.
	Upsert Operation = "upsert"
	// Query returns data that has not been deleted or archived
	Query Operation = "query"
	// QueryAll Rreturns records that have been deleted because of a merge or delete, and returns information about archived Task and Event records
	QueryAll Operation = "queryAll"
)

// State is the current state of processing for the job.
type State string

const (
	// Open the job has been created and job data can be uploaded tothe job.
	Open State = "Open"
	// UpdateComplete all data for the job has been uploaded and the job is ready to be queued and processed.
	UpdateComplete State = "UploadComplete"
	// InProgress Salesforce is processing the job
	InProgress State = "InProgress"
	// Aborted the job has been aborted.
	Aborted State = "Aborted"
	// JobComplete the job was processed by Salesforce.
	JobComplete State = "JobComplete"
	// Failed some records in the job failed.
	Failed State = "Failed"
)

// UnprocessedRecord is the unprocessed records from the job.
type UnprocessedRecord struct {
	Fields map[string]string
}

type concurrencyMode string

const (
	// ConcurrencySerial ...
	ConcurrencySerial concurrencyMode = "serial"
	// ConcurrencyParallel ...
	ConcurrencyParallel concurrencyMode = "parallel"
)

// JobRecord is the record for the job.  Includes the Salesforce ID along with the fields.
type JobRecord struct {
	ID string
	UnprocessedRecord
}

// SuccessfulRecord indicates for the record was created and the data that was uploaded.
type SuccessfulRecord struct {
	Created bool
	JobRecord
}

// FailedRecord indicates why the record failed and the data of the record.
type FailedRecord struct {
	Error string
	JobRecord
}

// Options are the options for the job.
//
// ColumnDelimiter is the delimiter used for the CSV job.  This field is optional.
//
// ContentType is the content type for the job.  This field is optional.
//
// ExternalIDFieldName is the external ID field in the object being updated.  Only needed for
// upsert operations.  This field is required for upsert operations.
//
// LineEnding is the line ending used for the CSV job data.  This field is optional.
//
// Object is the object type for the data bneing processed. This field is required.
//
// Operation is the processing operation for the job. This field is required.
type Options struct {
	ColumnDelimiter     ColumnDelimiter `json:"columnDelimiter"`
	ContentType         ContentType     `json:"contentType"`
	ExternalIDFieldName string          `json:"externalIdFieldName"`
	LineEnding          LineEnding      `json:"lineEnding"`
	Object              string          `json:"object"`
	Operation           Operation       `json:"operation"`
	Query               string          `json:"query"`
}

// Response is the response to job APIs.
type Response struct {
	APIVersion          float32 `json:"apiVersion"`
	ColumnDelimiter     string  `json:"columnDelimiter"`
	ConcurrencyMode     string  `json:"concurrencyMode"`
	ContentType         string  `json:"contentType"`
	ContentURL          string  `json:"contentUrl"`
	CreatedByID         string  `json:"createdById"`
	CreatedDate         string  `json:"createdDate"`
	ExternalIDFieldName string  `json:"externalIdFieldName"`
	ID                  string  `json:"id"`
	JobType             string  `json:"jobType"`
	LineEnding          string  `json:"lineEnding"`
	Object              string  `json:"object"`
	Operation           string  `json:"operation"`
	State               string  `json:"state"`
	SystemModstamp      string  `json:"systemModstamp"`
}

// Info is the response to the job information API.
type Info struct {
	Response
	ApexProcessingTime      int    `json:"apexProcessingTime"`
	APIActiveProcessingTime int    `json:"apiActiveProcessingTime"`
	NumberRecordsFailed     int    `json:"numberRecordsFailed"`
	NumberRecordsProcessed  int    `json:"numberRecordsProcessed"`
	Retries                 int    `json:"retries"`
	TotalProcessingTime     int    `json:"totalProcessingTime"`
	ErrorMessage            string `json:"errorMessage"`
}

// Job is the bulk job.
type Job struct {
	session session.ServiceFormatter
	options Options
	info    Response
	jobType JobType
}

func (j *Job) create(options Options) error {
	j.options = options
	err := j.formatOptions()
	if err != nil {
		return err
	}
	if j.options.Operation == "query" || j.options.Operation == "queryAll" {
		j.jobType = "V2Query"
	} else {
		j.jobType = "V2Ingest"
	}
	j.info, err = j.createCallout()
	if err != nil {
		return err
	}

	return nil
}

func (j *Job) formatOptions() error {
	if j.options.Operation == "" {
		return errors.New("bulk job: operation is required")
	}
	if j.options.Operation == Upsert {
		if j.options.ExternalIDFieldName == "" {
			return errors.New("bulk job: external id field name is required for upsert operation")
		}
	}
	if j.options.Object == "" && (j.options.Operation != "query" && j.options.Operation != "queryAll") {
		return errors.New("bulk job: object is required")
	}
	if j.options.LineEnding == "" {
		j.options.LineEnding = Linefeed
	}
	if j.options.ContentType == "" {
		j.options.ContentType = CSV
	}
	if j.options.ColumnDelimiter == "" {
		j.options.ColumnDelimiter = Comma
	}
	if j.options.Operation == "query" || j.options.Operation == "queryAll" {
		if j.options.Query == "" {
			return errors.New("bulk job: query is required for query or queryAll operations")
		}
	}
	return nil
}

func (j *Job) createCallout() (Response, error) {
	url := j.session.ServiceURL() + bulk2Endpoint(j.jobType)
	body, err := json.Marshal(j.options)
	if err != nil {
		return Response{}, err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return Response{}, err
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	j.session.AuthorizationHeader(request)

	return j.response(request)
}

func (j *Job) response(request *http.Request) (Response, error) {
	response, err := j.session.Client().Do(request)
	if err != nil {
		return Response{}, err
	}

	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var errs []sfdc.Error
		err = decoder.Decode(&errs)
		var errMsg error
		if err == nil {
			for _, err := range errs {
				errMsg = fmt.Errorf("insert response err: %s: %s", err.ErrorCode, err.Message)
			}
		} else {
			errMsg = fmt.Errorf("insert response err: %d %s", response.StatusCode, response.Status)
		}

		return Response{}, errMsg
	}

	var value Response
	err = decoder.Decode(&value)
	if err != nil {
		return Response{}, err
	}
	return value, nil
}

// Info returns the current job information.
func (j *Job) Info() (Info, error) {
	url := j.session.ServiceURL() + bulk2Endpoint(j.jobType) + "/" + j.info.ID
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Info{}, err
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	j.session.AuthorizationHeader(request)

	return j.infoResponse(request)
}

func (j *Job) infoResponse(request *http.Request) (Info, error) {
	response, err := j.session.Client().Do(request)
	if err != nil {
		return Info{}, err
	}

	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var errs []sfdc.Error
		err = decoder.Decode(&errs)
		var errMsg error
		if err == nil {
			for _, err := range errs {
				errMsg = fmt.Errorf("job err: %s: %s", err.ErrorCode, err.Message)
			}
		} else {
			errMsg = fmt.Errorf("job err: %d %s", response.StatusCode, response.Status)
		}

		return Info{}, errMsg
	}

	var value Info
	err = decoder.Decode(&value)
	if err != nil {
		return Info{}, err
	}
	return value, nil
}

func (j *Job) setState(state State) (Response, error) {
	url := j.session.ServiceURL() + bulk2Endpoint(j.jobType) + "/" + j.info.ID
	jobState := struct {
		State string `json:"state"`
	}{
		State: string(state),
	}
	body, err := json.Marshal(jobState)
	if err != nil {
		return Response{}, err
	}
	request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		return Response{}, err
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	j.session.AuthorizationHeader(request)

	return j.response(request)
}

// Close will close the current job.
func (j *Job) Close() (Response, error) {
	return j.setState(UpdateComplete)
}

// Abort will abort the current job.
func (j *Job) Abort() (Response, error) {
	return j.setState(Aborted)
}

// Delete will delete the current job.
func (j *Job) Delete() error {
	url := j.session.ServiceURL() + bulk2Endpoint(j.jobType) + "/" + j.info.ID
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	j.session.AuthorizationHeader(request)

	response, err := j.session.Client().Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusNoContent {
		return errors.New("job error: unable to delete job")
	}
	return nil
}

// Upload will upload data to processing.
func (j *Job) Upload(body io.Reader) error {
	url := j.session.ServiceURL() + bulk2Endpoint(j.jobType) + "/" + j.info.ID + "/batches"
	request, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "text/csv")
	j.session.AuthorizationHeader(request)

	response, err := j.session.Client().Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusCreated {
		return errors.New("job error: unable to upload job")
	}
	return nil
}

// QueryResults Gets the results for a query job. The job must have the state JobComplete.
func (j *Job) QueryResults(w io.Writer, maxRecords int, locator string) error {
	url := j.session.ServiceURL() + bulk2Endpoint(j.jobType) + "/" + j.info.ID + "/results"
	if locator != "" {
		url += "?locator=" + locator
		if maxRecords > 0 {
			url += "&maxRecords=" + string(maxRecords)
		}
	} else if maxRecords > 0 {
		url += "?maxRecords=" + string(maxRecords)
	}
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	request.Header.Add("Accept", "text/csv")
	j.session.AuthorizationHeader(request)

	response, err := j.session.Client().Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", response.Status)
	}

	r := bufio.NewReader(response.Body)
	defer response.Body.Close()

	if _, err = r.WriteTo(w); err != nil {
		return err
	}
	// response.Header.Get("Sforce-NumberOfRecords")
	if response.Header.Get("Sforce-Locator") != "" {
		return j.QueryResults(w, maxRecords, response.Header.Get("Sforce-Locator"))
	}

	return nil
}

// SuccessfulRecords returns the successful records for the job.
func (j *Job) SuccessfulRecords() ([]SuccessfulRecord, error) {
	url := j.session.ServiceURL() + bulk2Endpoint(j.jobType) + "/" + j.info.ID + "/successfulResults/"
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", "text/csv")
	j.session.AuthorizationHeader(request)

	response, err := j.session.Client().Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		decoder := json.NewDecoder(response.Body)
		defer response.Body.Close()
		var errs []sfdc.Error
		err = decoder.Decode(&errs)
		var errMsg error
		if err == nil {
			for _, err := range errs {
				errMsg = fmt.Errorf("job err: %s: %s", err.ErrorCode, err.Message)
			}
		} else {
			errMsg = fmt.Errorf("job err: %d %s", response.StatusCode, response.Status)
		}
		return nil, errMsg
	}

	scanner := bufio.NewScanner(response.Body)
	defer response.Body.Close()
	scanner.Split(bufio.ScanLines)
	var records []SuccessfulRecord
	delimiter := j.delimiter()
	columns, err := j.recordResultHeader(scanner, delimiter)
	if err != nil {
		return nil, err
	}
	createIdx, err := j.headerPosition(`sf__Created`, columns)
	if err != nil {
		return nil, err
	}
	idIdx, err := j.headerPosition(`sf__Id`, columns)
	if err != nil {
		return nil, err
	}
	fields := j.fields(columns, 2)
	for scanner.Scan() {
		var record SuccessfulRecord
		values := strings.Split(scanner.Text(), delimiter)
		isCreated := strings.Replace(values[createIdx], "\"", "", -1)
		created, err := strconv.ParseBool(isCreated)
		if err != nil {
			return nil, err
		}
		record.Created = created
		record.ID = values[idIdx]
		record.Fields = j.record(fields, values[2:])
		records = append(records, record)
	}

	return records, nil
}

// FailedRecords returns the failed records for the job.
func (j *Job) FailedRecords() ([]FailedRecord, error) {
	url := j.session.ServiceURL() + bulk2Endpoint(j.jobType) + "/" + j.info.ID + "/failedResults/"
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", "text/csv")
	j.session.AuthorizationHeader(request)

	response, err := j.session.Client().Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		decoder := json.NewDecoder(response.Body)
		defer response.Body.Close()
		var errs []sfdc.Error
		err = decoder.Decode(&errs)
		var errMsg error
		if err == nil {
			for _, err := range errs {
				errMsg = fmt.Errorf("job err: %s: %s", err.ErrorCode, err.Message)
			}
		} else {
			errMsg = fmt.Errorf("job err: %d %s", response.StatusCode, response.Status)
		}
		return nil, errMsg
	}

	scanner := bufio.NewScanner(response.Body)
	defer response.Body.Close()
	scanner.Split(bufio.ScanLines)
	var records []FailedRecord
	delimiter := j.delimiter()
	columns, err := j.recordResultHeader(scanner, delimiter)
	if err != nil {
		return nil, err
	}
	errorIdx, err := j.headerPosition(`sf__Error`, columns)
	if err != nil {
		return nil, err
	}
	idIdx, err := j.headerPosition(`sf__Id`, columns)
	if err != nil {
		return nil, err
	}
	fields := j.fields(columns, 2)
	for scanner.Scan() {
		var record FailedRecord
		values := strings.Split(scanner.Text(), delimiter)
		record.Error = values[errorIdx]
		record.ID = values[idIdx]
		record.Fields = j.record(fields, values[2:])
		records = append(records, record)
	}

	return records, nil
}

// UnprocessedRecords returns the unprocessed records for the job.
func (j *Job) UnprocessedRecords() ([]UnprocessedRecord, error) {
	url := j.session.ServiceURL() + bulk2Endpoint(j.jobType) + "/" + j.info.ID + "/unprocessedrecords/"
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", "text/csv")
	j.session.AuthorizationHeader(request)

	response, err := j.session.Client().Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		decoder := json.NewDecoder(response.Body)
		defer response.Body.Close()
		var errs []sfdc.Error
		err = decoder.Decode(&errs)
		var errMsg error
		if err == nil {
			for _, err := range errs {
				errMsg = fmt.Errorf("job err: %s: %s", err.ErrorCode, err.Message)
			}
		} else {
			errMsg = fmt.Errorf("job err: %d %s", response.StatusCode, response.Status)
		}
		return nil, errMsg
	}

	scanner := bufio.NewScanner(response.Body)
	defer response.Body.Close()
	scanner.Split(bufio.ScanLines)
	var records []UnprocessedRecord
	delimiter := j.delimiter()
	columns, err := j.recordResultHeader(scanner, delimiter)
	if err != nil {
		return nil, err
	}
	fields := j.fields(columns, 0)
	for scanner.Scan() {
		var record UnprocessedRecord
		values := strings.Split(scanner.Text(), delimiter)
		record.Fields = j.record(fields, values)
		records = append(records, record)
	}

	return records, nil
}

func (j *Job) recordResultHeader(scanner *bufio.Scanner, delimiter string) ([]string, error) {
	if scanner.Scan() == false {
		return nil, errors.New("job: response needs to have header")
	}
	text := strings.Replace(scanner.Text(), "\"", "", -1)
	return strings.Split(text, delimiter), nil
}
func (j *Job) headerPosition(column string, header []string) (int, error) {
	for idx, col := range header {
		if col == column {
			return idx, nil
		}
	}
	return -1, fmt.Errorf("job header: %s column is not in header", column)
}
func (j *Job) fields(header []string, offset int) []string {
	fields := make([]string, len(header)-offset)
	copy(fields[:], header[offset:])
	return fields
}
func (j *Job) record(fields, values []string) map[string]string {
	record := make(map[string]string)
	for idx, field := range fields {
		record[field] = values[idx]
	}
	return record
}

func (j *Job) delimiter() string {
	switch ColumnDelimiter(j.info.ColumnDelimiter) {
	case Tab:
		return "\t"
	case SemiColon:
		return ";"
	case Pipe:
		return "|"
	case Caret:
		return "^"
	case Backquote:
		return "`"
	default:
		return ","
	}
}

func (j *Job) newline() string {
	switch LineEnding(j.info.LineEnding) {
	case CarriageReturnLinefeed:
		return "\r\n"
	default:
		return "\n"
	}
}
