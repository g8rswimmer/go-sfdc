package bulk

import (
	"fmt"
	"io"

	"github.com/TuSKan/go-sfdc/session"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func bulk2Endpoint(jobType JobType) string {
	if jobType == "V2Query" {
		return "/jobs/query"
	}
	return "/jobs/ingest"
}

// Resource is the structure that can be used to create bulk 2.0 jobs.
type Resource struct {
	session session.ServiceFormatter
}

// NewResource creates a new bulk 2.0 REST resource.  If the session is nil
// an error will be returned.
func NewResource(session session.ServiceFormatter) (*Resource, error) {
	if err := session.Validade(); err != nil {
		return nil, fmt.Errorf("bulk: %v", err)
	}
	return &Resource{
		session: session,
	}, nil
}

// CreateJob will create a new bulk 2.0 job from the options that where passed.
// The Job that is returned can be used to upload object data to the Salesforce org.
func (r *Resource) CreateJob(options Options) (*Job, error) {
	job := &Job{
		session: r.session,
	}
	if err := job.create(options); err != nil {
		return nil, err
	}

	return job, nil
}

// JobsInfo will retrieve all of the bulk 2.0 jobs.
func (r *Resource) JobsInfo(parameters Parameters) ([]Response, error) {
	jobs, err := jobsInfo(r.session, parameters)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

// QueryJobsResults ...
func (r *Resource) QueryJobsResults(jobs []*Job, writers []io.Writer, parameters Parameters, bo Backoff, maxRecords int) (map[string]int, map[string]error, error) {

	if len(jobs) != len(writers) {
		return map[string]int{}, map[string]error{}, fmt.Errorf("len(jobs) %d != len(writes) %d", len(jobs), len(writers))
	}

	nrecMap := make(map[string]int)
	errsMap := make(map[string]error)
	jobsMap := make(map[string]*Job)
	writersMap := make(map[string]io.Writer)
	for i, j := range jobs {
		jobsMap[j.info.ID] = jobs[i]
		writersMap[j.info.ID] = writers[i]
	}

	return nrecMap, errsMap, Retry(bo, func() (bool, error) {
		jobsResp, err := r.JobsInfo(parameters)
		if err != nil {
			return false, err
		}
		for _, res := range jobsResp {
			if jobsMap[res.ID] != nil && State(res.State) == JobComplete {
				n, err := jobsMap[res.ID].QueryResults(writersMap[res.ID], maxRecords, "")
				if err != nil {
					errsMap[res.ID] = err
				} else {
					nrecMap[res.ID] = n
				}
				delete(jobsMap, res.ID)
				//delete(writersMap, res.ID)
			}
		}
		if len(jobsMap) == 0 {
			return true, nil
		}
		return false, nil
	})
}
