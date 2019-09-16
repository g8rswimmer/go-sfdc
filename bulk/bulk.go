package bulk

import (
	"fmt"
	"io"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/g8rswimmer/go-sfdc/session"
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
	if session == nil {
		return nil, fmt.Errorf("bulk: session can not be nil")
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
func (r *Resource) QueryJobsResults(jobs []*Job, writers []io.Writer, parameters Parameters, waitMaxDuration time.Duration, maxRecords int) (map[string]error, error) {

	if len(jobs) != len(writers) {
		return map[string]error{}, fmt.Errorf("len(jobs) %d != len(writes) %d", len(jobs), len(writers))
	}

	errsMap := make(map[string]error)
	jobsMap := make(map[string]*Job)
	writersMap := make(map[string]io.Writer)
	for i, j := range jobs {
		jobsMap[j.info.ID] = jobs[i]
		writersMap[j.info.ID] = writers[i]
		errsMap[j.info.ID] = nil
	}

	return errsMap, wait.ExponentialBackoff(wait.Backoff{
		Duration: 100 * time.Millisecond,
		Jitter:   0.5,
		Factor:   1.5,
		Cap:      30 * time.Second,
		Steps:    min(1, int(waitMaxDuration.Seconds()/30.0)),
	}, func() (bool, error) {
		jobsResp, err := r.JobsInfo(parameters)
		if err != nil {
			return false, err
		}
		for _, res := range jobsResp {
			if jobsMap[res.ID] != nil && State(res.State) == JobComplete {
				errsMap[res.ID] = jobsMap[res.ID].QueryResults(writersMap[res.ID], maxRecords, "")
				delete(jobsMap, res.ID)
			}
		}
		if len(jobsMap) == 0 {
			return true, nil
		}
		return false, nil
	})
}
