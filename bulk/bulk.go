package bulk

import (
	"errors"

	"github.com/g8rswimmer/go-sfdc/session"
)

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
		return nil, errors.New("bulk: session can not be nil")
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

// AllJobs will retrieve all of the bulk 2.0 jobs.
func (r *Resource) AllJobs(parameters Parameters) (*Jobs, error) {
	jobs, err := newJobs(r.session, parameters)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

// WaitJobs - Wait jobs
func (r *Resource) WaitJobs(parameters Parameters) error {
	return return wait.ExponentialBackoff(wait.Backoff{
		Duration: 100 * time.Millisecond,
		Jitter:   0.5,
		Factor:   1.5,
		Cap:    60*time.Second,
		Steps: 10,
	}, func() (bool, error) {
		jobs, err := AllJobs(parameters)
		if err != nil {
			return false, err
		}
		jobs.Response
		if State(I.Response.State) == JobComplete {
			return true, nil
		}
		return false, nil
	})
}