package bulk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	sfdc "github.com/g8rswimmer/go-sfdc"
	"github.com/g8rswimmer/go-sfdc/session"
)

// Parameters to query all of the bulk jobs.
//
// IsPkChunkingEnabled will filter jobs with PK chunking enabled.
//
// JobType will filter jobs based on job type.
type Parameters struct {
	IsPkChunkingEnabled bool
	JobType             JobType
	ConcurrencyMode     concurrencyMode
	QueryLocator        string
}

type jobResponse struct {
	Done           bool       `json:"done"`
	Records        []Response `json:"records"`
	NextRecordsURL string     `json:"nextRecordsUrl"`
}

// Jobs presents the response from the all jobs request.
type Jobs struct {
	session  session.ServiceFormatter
	response jobResponse
}

func newJobs(session session.ServiceFormatter, parameters Parameters) (*Jobs, error) {
	j := &Jobs{
		session: session,
	}
	url := session.ServiceURL() + bulk2Endpoint(parameters.JobType)
	request, err := j.request(url)
	if err != nil {
		return nil, err
	}
	q := request.URL.Query()
	q.Add("isPkChunkingEnabled", strconv.FormatBool(parameters.IsPkChunkingEnabled))
	q.Add("jobType", string(parameters.JobType))
	request.URL.RawQuery = q.Encode()

	response, err := j.do(request)
	if err != nil {
		return nil, err
	}
	j.response = response
	return j, nil
}

// Done indicates whether there are more jobs to get.
func (j *Jobs) Done() bool {
	return j.response.Done
}

// Records contains the information for each retrieved job.
func (j *Jobs) Records() []Response {
	return j.response.Records
}

// Next will retrieve the next batch of job information.
func (j *Jobs) Next() (*Jobs, error) {
	if j.Done() == true {
		return nil, errors.New("jobs: there is no more records")
	}
	request, err := j.request(j.response.NextRecordsURL)
	if err != nil {
		return nil, err
	}
	response, err := j.do(request)
	if err != nil {
		return nil, err
	}
	return &Jobs{
		session:  j.session,
		response: response,
	}, nil
}
func (j *Jobs) request(url string) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", "application/json")
	j.session.AuthorizationHeader(request)
	return request, nil
}
func (j *Jobs) do(request *http.Request) (jobResponse, error) {
	response, err := j.session.Client().Do(request)
	if err != nil {
		return jobResponse{}, err
	}

	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var jobsErrs []sfdc.Error
		err = decoder.Decode(&jobsErrs)
		var errMsg error
		if err == nil {
			for _, jobErr := range jobsErrs {
				errMsg = fmt.Errorf("insert response err: %s: %s", jobErr.ErrorCode, jobErr.Message)
			}
		} else {
			errMsg = fmt.Errorf("insert response err: %d %s", response.StatusCode, response.Status)
		}

		return jobResponse{}, errMsg
	}

	var value jobResponse
	err = decoder.Decode(&value)
	if err != nil {
		return jobResponse{}, err
	}
	return value, nil

}
