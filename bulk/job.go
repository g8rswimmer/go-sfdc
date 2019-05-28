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

type JobType string

const (
	BigObjects JobType = "BigObjectIngest"
	Classic    JobType = "Classic"
	V2Ingest   JobType = "V2Ingest"
)

type ColumnDelimiter string

const (
	Backquote ColumnDelimiter = "BACKQUOTE"
	Caret     ColumnDelimiter = "CARET"
	Comma     ColumnDelimiter = "COMMA"
	Pipe      ColumnDelimiter = "PIPE"
	SemiColon ColumnDelimiter = "SEMICOLON"
	Tab       ColumnDelimiter = "TAB"
)

type ContentType string

const CSV ContentType = "CSV"

type LineEnding string

const (
	Linefeed               LineEnding = "LF"
	CarriageReturnLinefeed LineEnding = "CRLF"
)

type Operation string

const (
	Insert Operation = "insert"
	Delete Operation = "delete"
	Update Operation = "update"
	Upsert Operation = "upsert"
)

type State string

const (
	Open           State = "Open"
	UpdateComplete State = "UploadComplete"
	Aborted        State = "Aborted"
	JobComplete    State = "JobComplete"
	Fialed         State = "Failed"
)

const createEndpoint = "/jobs/ingest"

type JobUnprocessedRecord struct {
	Fields map[string]string
}
type JobRecord struct {
	ID string
	JobUnprocessedRecord
}
type JobSuccessfulRecord struct {
	Created bool
	JobRecord
}
type JobFailedRecord struct {
	Error string
	JobRecord
}
type JobOptions struct {
	ColumnDelimiter     ColumnDelimiter `json:"columnDelimiter"`
	ContentType         ContentType     `json:"contentType"`
	ExternalIDFieldName string          `json:"externalIdFieldName"`
	LineEnding          LineEnding      `json:"lineEnding"`
	Object              string          `json:"object"`
	Operation           Operation       `json:"operation"`
}

type JobResponse struct {
	APIVersion          string `json:"apiVersion"`
	ColumnDelimiter     string `json:"columnDelimiter"`
	ConcurrencyMode     string `json:"concurrencyMode"`
	ContentType         string `json:"contentType"`
	ContentURL          string `json:"contentUrl"`
	CreatedByID         string `json:"createdById"`
	CreatedDate         string `json:"createdDate"`
	ExternalIDFieldName string `json:"externalIdFieldName"`
	ID                  string `json:"id"`
	JobType             string `json:"jobType"`
	LineEnding          string `json:"lineEnding"`
	Object              string `json:"object"`
	Operation           string `json:"operation"`
	State               string `json:"state"`
	SystemModstamp      string `json:"systemModstamp"`
}

type JobInfo struct {
	JobResponse
	ApexProcessingTime      int `json:"apexProcessingTime"`
	APIActiveProcessingTime int `json:"apiActiveProcessingTime"`
	NumberRecordsFailed     int `json:"numberRecordsFailed"`
	NumberRecordsProcessed  int `json:"numberRecordsProcessed"`
	Retries                 int `json:"retries"`
	TotalProcessingTime     int `json:"totalProcessingTime"`
}
type Job struct {
	session session.ServiceFormatter
	info    JobResponse
}

func (j *Job) create(options JobOptions) error {
	err := j.formatOptions(&options)
	if err != nil {
		return err
	}
	j.info, err = j.createCallout(options)
	if err != nil {
		return err
	}

	return nil
}

func (j *Job) formatOptions(options *JobOptions) error {
	if options.Operation == "" {
		return errors.New("bulk job: operation is required")
	}
	if options.Operation == Upsert {
		if options.ExternalIDFieldName == "" {
			return errors.New("bulk job: external id field name is required for upsert operation")
		}
	}
	if options.Object == "" {
		return errors.New("bulk job: object is required")
	}
	if options.LineEnding == "" {
		options.LineEnding = Linefeed
	}
	if options.ContentType == "" {
		options.ContentType = CSV
	}
	if options.ColumnDelimiter == "" {
		options.ColumnDelimiter = Comma
	}
	return nil
}

func (j *Job) createCallout(options JobOptions) (JobResponse, error) {
	url := j.session.ServiceURL() + createEndpoint
	body, err := json.Marshal(options)
	if err != nil {
		return JobResponse{}, err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return JobResponse{}, err
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	j.session.AuthorizationHeader(request)

	return j.response(request)
}

func (j *Job) response(request *http.Request) (JobResponse, error) {
	response, err := j.session.Client().Do(request)
	if err != nil {
		return JobResponse{}, err
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

		return JobResponse{}, errMsg
	}

	var value JobResponse
	err = decoder.Decode(&value)
	if err != nil {
		return JobResponse{}, err
	}
	return value, nil
}
func (j *Job) Info() (JobInfo, error) {
	url := j.session.ServiceURL() + createEndpoint + "/" + j.info.ID
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return JobInfo{}, err
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	j.session.AuthorizationHeader(request)

	return j.infoResponse(request)
}

func (j *Job) infoResponse(request *http.Request) (JobInfo, error) {
	response, err := j.session.Client().Do(request)
	if err != nil {
		return JobInfo{}, err
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

		return JobInfo{}, errMsg
	}

	var value JobInfo
	err = decoder.Decode(&value)
	if err != nil {
		return JobInfo{}, err
	}
	return value, nil
}

func (j *Job) setState(state State) (JobResponse, error) {
	url := j.session.ServiceURL() + createEndpoint + "/" + j.info.ID
	jobState := struct {
		State string
	}{
		State: string(state),
	}
	body, err := json.Marshal(jobState)
	if err != nil {
		return JobResponse{}, err
	}
	request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		return JobResponse{}, err
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	j.session.AuthorizationHeader(request)

	return j.response(request)
}

func (j *Job) Close() (JobResponse, error) {
	return j.setState(UpdateComplete)
}

func (j *Job) Abort() (JobResponse, error) {
	return j.setState(Aborted)
}

func (j *Job) Delete() error {
	url := j.session.ServiceURL() + createEndpoint + "/" + j.info.ID
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

func (j *Job) Upload(body io.Reader) error {
	url := j.session.ServiceURL() + createEndpoint + "/" + j.info.ID + "/batches"
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

func (j *Job) SuccessfulRecords() ([]JobSuccessfulRecord, error) {
	url := j.session.ServiceURL() + createEndpoint + "/" + j.info.ID + "/successfulResults/"
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
	var records []JobSuccessfulRecord
	delimiter := j.delimiter()
	fields, err := j.fields(scanner, delimiter, 2)
	if err != nil {
		return nil, err
	}
	for scanner.Scan() {
		var record JobSuccessfulRecord
		values := strings.Split(scanner.Text(), delimiter)
		created, err := strconv.ParseBool(values[0])
		if err != nil {
			return nil, err
		}
		record.Created = created
		record.ID = values[1]
		record.Fields = j.record(fields, values[2:])
		records = append(records, record)
	}

	return records, nil
}

func (j *Job) FailedRecords() ([]JobFailedRecord, error) {
	url := j.session.ServiceURL() + createEndpoint + "/" + j.info.ID + "/failedResults/"
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
	var records []JobFailedRecord
	delimiter := j.delimiter()
	fields, err := j.fields(scanner, delimiter, 2)
	if err != nil {
		return nil, err
	}
	for scanner.Scan() {
		var record JobFailedRecord
		values := strings.Split(scanner.Text(), delimiter)
		record.Error = values[0]
		record.ID = values[1]
		record.Fields = j.record(fields, values[2:])
		records = append(records, record)
	}

	return records, nil
}

func (j *Job) UnprocessedRecords() ([]JobUnprocessedRecord, error) {
	url := j.session.ServiceURL() + createEndpoint + "/" + j.info.ID + "/failedResults/"
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
	var records []JobUnprocessedRecord
	delimiter := j.delimiter()
	fields, err := j.fields(scanner, delimiter, 0)
	if err != nil {
		return nil, err
	}
	for scanner.Scan() {
		var record JobUnprocessedRecord
		values := strings.Split(scanner.Text(), delimiter)
		record.Fields = j.record(fields, values)
		records = append(records, record)
	}

	return records, nil
}

func (j *Job) fields(scanner *bufio.Scanner, delimiter string, offset int) ([]string, error) {
	if scanner.Scan() == false {
		return nil, errors.New("job: response needs to have header")
	}
	values := strings.Split(scanner.Text(), delimiter)
	fields := make([]string, len(values)-offset)
	copy(values[offset:], fields[:])
	return fields, nil
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
	case Comma:
		fallthrough
	default:
		return ","
	}
}
