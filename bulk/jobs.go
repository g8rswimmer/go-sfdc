package bulk

import (
	"encoding/json"
	"fmt"
	"net/http"
	//"strconv"

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

type jobsResponse struct {
	Done           bool       `json:"done"`
	Records        []Response `json:"records"`
	NextRecordsURL string     `json:"nextRecordsUrl"`
}


func jobsInfo(session session.ServiceFormatter, parameters Parameters) ([]Response, error) {
	var responses []Response
	url := session.ServiceURL() + bulk2Endpoint(parameters.JobType)
	request, err := jobsRequest(session, url)
	if err != nil {
		return nil, err
	}
	q := request.URL.Query()
	//q.Add("isPkChunkingEnabled", strconv.FormatBool(parameters.IsPkChunkingEnabled))
	//q.Add("jobType", string(parameters.JobType))
	request.URL.RawQuery = q.Encode()

	jobResp, err := jobsDo(session, request)
	if err != nil {
		return nil, err
	}
	responses = append(responses, jobResp.Records...)
	for jobResp.Done != true {
		var err error
		request, err = jobsRequest(session, jobResp.NextRecordsURL)
		if err != nil {
			return responses, err
		}
		jobResp, err = jobsDo(session, request)
		if err != nil {
			return responses, err
		}
		responses = append(responses, jobResp.Records...)
	}
	return responses, nil
}


func jobsRequest(session session.ServiceFormatter, url string) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", "application/json")
	session.AuthorizationHeader(request)
	return request, nil
}


func jobsDo(session session.ServiceFormatter, request *http.Request) (jobsResponse, error) {
	response, err := session.Client().Do(request)
	if err != nil {
		return jobsResponse{}, err
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

		return jobsResponse{}, errMsg
	}

	var value jobsResponse
	err = decoder.Decode(&value)
	if err != nil {
		return jobsResponse{}, err
	}
	return value, nil
}
