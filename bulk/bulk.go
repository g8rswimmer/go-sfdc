package bulk

import "github.com/g8rswimmer/go-sfdc/session"

type Resource struct {
	session session.ServiceFormatter
}

func (r *Resource) CreateJob(options Options) (*Job, error) {
	job := &Job{
		session: r.session,
	}
	if err := job.create(options); err != nil {
		return nil, err
	}

	return job, nil
}
