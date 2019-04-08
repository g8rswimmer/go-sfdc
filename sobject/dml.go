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

// UpsertValue is the value that is return when a
// record as been upserted into Salesforce.
//
// Upsert will return two types of values, which
// are indicated by Inserted.  If Inserted is true,
// then the InsertValue is popluted.
type UpsertValue struct {
	Inserted bool
	InsertValue
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
	Inserter
	ID() string
}

// Upserter provides the parameters needed to upsert a record.
//
// SObject is the Salesforce table name.  An example would be Account or Custom__c.
//
// ID is the External ID that will be updated.
//
// ExternalField is the external ID field.
//
// Fields are the fields of the record that are to be inserted.  It is the
// callers responsbility to provide value fields and values.
type Upserter interface {
	Updater
	ExternalField() string
}

// Deleter provides the parameters needed to delete a record.
//
// SObject is the Salesforce table name.  An example would be Account or Custom__c.
//
// ID is the Salesforce ID to be deleted.
type Deleter interface {
	SObject() string
	ID() string
}

type dml struct {
	session session.Formatter
}

func (d *dml) insertCallout(inserter Inserter) (InsertValue, error) {
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
		var errMsg error
		if err == nil {
			for _, insertErr := range insertErrs {
				errMsg = fmt.Errorf("insert response err: %s: %s", insertErr.ErrorCode, insertErr.Message)
			}
		} else {
			errMsg = fmt.Errorf("insert response err: %d %s", response.StatusCode, response.Status)
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

func (d *dml) updateCallout(updater Updater) error {
	request, err := d.updateRequest(updater)

	if err != nil {
		return err
	}

	return d.updateResponse(request)

}

func (d *dml) updateRequest(updater Updater) (*http.Request, error) {

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

func (d *dml) updateResponse(request *http.Request) error {
	response, err := d.session.Client().Do(request)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("metadata response err: %d %s", response.StatusCode, response.Status)
	}

	return nil
}
func (d *dml) upsertCallout(upserter Upserter) (UpsertValue, error) {
	request, err := d.upsertRequest(upserter)

	if err != nil {
		return UpsertValue{}, err
	}

	value, err := d.upsertResponse(request)

	if err != nil {
		return UpsertValue{}, err
	}

	return value, nil

}
func (d *dml) upsertRequest(upserter Upserter) (*http.Request, error) {

	url := d.session.ServiceURL() + objectEndpoint + upserter.SObject() + "/" + upserter.ExternalField() + "/" + upserter.ID()

	body, err := json.Marshal(upserter.Fields())
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

func (d *dml) upsertResponse(request *http.Request) (UpsertValue, error) {
	response, err := d.session.Client().Do(request)

	if err != nil {
		return UpsertValue{}, err
	}

	decoder := json.NewDecoder(response.Body)

	var isInsert bool
	var value UpsertValue

	switch response.StatusCode {
	case http.StatusCreated:
		defer response.Body.Close()
		isInsert = true
		err = decoder.Decode(&value)
		if err != nil {
			return UpsertValue{}, err
		}
	case http.StatusNoContent:
		isInsert = false
	default:
		defer response.Body.Close()
		var insertErrs []insertError
		err = decoder.Decode(&insertErrs)
		errMsg := fmt.Errorf("upsert response err: %d %s", response.StatusCode, response.Status)
		if err == nil {
			for _, insertErr := range insertErrs {
				errMsg = fmt.Errorf("upsert response err: %s: %s", insertErr.ErrorCode, insertErr.Message)
			}
		}
		return UpsertValue{}, errMsg
	}

	value.Inserted = isInsert

	return value, nil
}
func (d *dml) deleteCallout(deleter Deleter) error {

	request, err := d.deleteRequest(deleter)

	if err != nil {
		return err
	}

	return d.deleteResponse(request)
}
func (d *dml) deleteRequest(deleter Deleter) (*http.Request, error) {

	url := d.session.ServiceURL() + objectEndpoint + deleter.SObject() + "/" + deleter.ID()

	request, err := http.NewRequest(http.MethodDelete, url, nil)

	if err != nil {
		return nil, err
	}

	d.session.AuthorizationHeader(request)
	return request, nil

}

func (d *dml) deleteResponse(request *http.Request) error {
	response, err := d.session.Client().Do(request)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete has failed %d %s", response.StatusCode, response.Status)
	}

	return nil
}
