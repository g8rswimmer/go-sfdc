package bulk

import (
	"errors"

	"github.com/g8rswimmer/go-sfdc/session"
)

type Resource struct {
	session session.ServiceFormatter
}

func NewResource(session session.ServiceFormatter) (*Resource, error) {
	if session == nil {
		return nil, errors.New("bulk: session can not be nil")
	}
	return &Resource{
		session: session,
	}, nil
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
