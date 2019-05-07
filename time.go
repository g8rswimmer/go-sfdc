package sfdc

import (
	"errors"
	"time"
)

// SalesforceDateTime is the format returned by Salesforce TimeDate field type.
const SalesforceDateTime = "2006-01-02T15:04:05.000+0000"

// SalesforceDate is the format returned by the Salesforce Date field type.
const SalesforceDate = "2006-01-02"

var layouts = []string{
	time.RFC3339,
	SalesforceDateTime,
	SalesforceDate,
}

// ParseTime attempts to parse a JSON time string from Salesforce.  It will attempt
// to parse the time using RFC 3339, then Salesforce DateTime format and lastly Salesforce
// Date format.
func ParseTime(salesforceTime string) (time.Time, error) {
	if salesforceTime == "" {
		return time.Time{}, errors.New("parse time: time string to decode can not be empty")
	}
	var err error
	for _, layout := range layouts {
		var date time.Time
		if date, err = time.Parse(layout, salesforceTime); err == nil {
			return date, nil
		}
	}
	return time.Time{}, err
}
