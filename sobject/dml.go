package sobject

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/g8rswimmer/goforce/session"
)

// InsertValue is the value that is returned when a
// record is inserted into Salesforce.
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

// Inserter provides the parameters needed insert a record.
//
// SObject is the Salesforce table name.  An example would be Account or Custom__c.
//
// Fields are the fields of the record that are to be inserted.  It is the
// callers responsbility to provide value fields and values.
type Inserter interface {
	SObject() string
	Fields() map[string]interface{}
}

// Updater provides the parameters needed to update a record.
//
// SObject is the Salesforce table name.  An example would be Account or Custom__c.
//
// ID is the Salesforce ID that will be updated.
//
// Fields are the fields of the record that are to be inserted.  It is the
// callers responsbility to provide value fields and values.
type Updater interface {
	SObject() string
	ID() string
	Fields() map[string]interface{}
}

type dml struct {
	session session.Formatter
}

func (d *dml) Insert(inserter Inserter) (InsertValue, error) {
	request, err := d.insertRequest(inserter)

	if err != nil {
		return InsertValue{}, err
	}

	value, err := d.insertResponse(request)

	if err != nil {
		return InsertValue{}, err
	}

	return value, nil
}
func (d *dml) insertRequest(inserter Inserter) (*http.Request, error) {

	url := d.session.ServiceURL() + objectEndpoint + inserter.SObject()

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
	d.session.AuthorizationHeader(request)
	return request, nil

}

func (d *dml) insertResponse(request *http.Request) (InsertValue, error) {
	response, err := d.session.Client().Do(request)

	if err != nil {
		return InsertValue{}, err
	}

	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
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

func (d *dml) Update(updater Updater) error {
	request, err := d.requestUpdate(updater)

	if err != nil {
		return err
	}

	return d.responseUpdate(request)

}

func (d *dml) requestUpdate(updater Updater) (*http.Request, error) {

	url := d.session.ServiceURL() + objectEndpoint + updater.SObject() + "/" + updater.ID()

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
	d.session.AuthorizationHeader(request)
	return request, nil

}

func (d *dml) responseUpdate(request *http.Request) error {
	response, err := d.session.Client().Do(request)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("metadata response err: %d %s", response.StatusCode, response.Status)
	}

	return nil
}
