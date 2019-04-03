package sobject

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/g8rswimmer/goforce/session"
)

type InsertValue struct {
	Success bool     `json:"success"`
	ID      string   `json:"id"`
	Errors  []string `json:"errors"`
}

type insertError struct {
	Message   string   `json:"message"`
	ErrorCode string   `json:"errorCode"`
	Fields    []string `json:"fields"`
}

type Inserter interface {
	SObject() string
	Fields() map[string]interface{}
}

type insert struct {
	session session.Formatter
}

func (i *insert) Insert(inserter Inserter) (InsertValue, error) {
	request, err := i.request(inserter)

	if err != nil {
		return InsertValue{}, err
	}

	value, err := i.response(request)

	if err != nil {
		return InsertValue{}, err
	}

	return value, nil
}
func (i *insert) request(inserter Inserter) (*http.Request, error) {

	url := i.session.ServiceURL() + objectEndpoint + inserter.SObject()

	body, err := json.Marshal(inserter.Fields())
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	i.session.AuthorizationHeader(request)
	return request, nil

}

func (i *insert) response(request *http.Request) (InsertValue, error) {
	response, err := i.session.Client().Do(request)

	if err != nil {
		return InsertValue{}, err
	}

	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var insertErrs []insertError
		err = decoder.Decode(&insertErrs)
		errMsg := fmt.Errorf("insert response err: %d %s", response.StatusCode, response.Status)
		if err == nil {
			for _, insertErr := range insertErrs {
				errMsg = fmt.Errorf("insert response err: %s: %s", insertErr.ErrorCode, insertErr.Message)
			}
		}

		return InsertValue{}, errMsg
	}

	var value InsertValue
	err = decoder.Decode(&value)
	if err != nil {
		return InsertValue{}, err
	}

	return value, nil
}
