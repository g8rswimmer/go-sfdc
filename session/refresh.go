package session

import (
	"errors"
	"sync"
	"time"

	"github.com/g8rswimmer/go-sfdc"
)

// Refresher is a structure will allow for a session token to be
// re-authenticated based on a time interval.
type Refresher struct {
	session     *Session
	sessionLock sync.RWMutex
	err         error
	stop        chan struct{}
}

// DefaultRefreshPeriod is the default refresh periods for the session.
// One hour for normal refresh, if and error every 5 seonds.
var DefaultRefreshPeriod = RefreshPeriod{
	Normal: 3600,
	Error:  5,
}

// RefreshPeriod is the time delays before requesting a new session.
type RefreshPeriod struct {
	Normal int
	Error  int
}

// OpenRefresh will open the first session and set up a re-occuring refresh of the session token.
func OpenRefresh(config sfdc.Configuration, period RefreshPeriod) (*Refresher, error) {
	if period.Normal <= 0 || period.Error <= 0 {
		return nil, errors.New("session refresh: refresh time can not be less than or equal to zero")
	}
	session, err := Open(config)
	if err != nil {
		return nil, err
	}
	refresher := &Refresher{
		session: session,
		stop:    make(chan struct{}),
	}

	go func() {
		sleep := time.Duration(period.Normal) * time.Second
		for {
			select {
			case <-refresher.stop:
				return
			default:
			}
			time.Sleep(sleep)
			if err := refresher.refresh(); err != nil {
				refresher.err = err
				sleep = time.Duration(period.Error) * time.Second
			} else {
				refresher.err = nil
				sleep = time.Duration(period.Normal) * time.Second
			}
		}
	}()

	return refresher, nil
}

// Session will return the current session
func (r *Refresher) Session() *Session {
	r.sessionLock.RLock()
	defer r.sessionLock.RUnlock()
	return r.session
}

func (r *Refresher) Error() error {
	return r.err
}

// Shutdown will close the refreshing of the seesion token.
func (r *Refresher) Shutdown() {
	close(r.stop)
}
func (r *Refresher) refresh() error {
	session, err := Open(r.session.config)
	if err != nil {
		return err
	}
	r.sessionLock.Lock()
	defer r.sessionLock.Unlock()
	r.session = session
	return nil
}
