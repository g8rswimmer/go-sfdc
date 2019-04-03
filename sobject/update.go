package sobject

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/g8rswimmer/goforce/session"
)

type Updater interface {
	SObject() string
	ID() string
	Fields() map[string]interface{}
}

type update struct {
	session session.Formatter
}

func (u *update) Update(updater Updater) error {
	request, err := u.request(updater)

	if err != nil {
		return err
	}

	return u.response(request)

}

func (u *update) request(updater Updater) (*http.Request, error) {

	url := u.session.ServiceURL() + objectEndpoint + updater.SObject() + "/" + updater.ID()

	body, err := json.Marshal(updater.Fields())
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	u.session.AuthorizationHeader(request)
	return request, nil

}

func (u *update) response(request *http.Request) error {
	response, err := u.session.Client().Do(request)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("metadata response err: %d %s", response.StatusCode, response.Status)
	}

	return nil
}
