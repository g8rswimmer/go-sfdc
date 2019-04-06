package goforce

import (
	"time"
)

// SalesforceDateTime is the format returned by Salesforce TimeDate field type.
const SalesforceDateTime = "2006-01-02T15:04:05.000+0000"

// SalesforceDate is the format returned by the Salesforce Date field type.
const SalesforceDate = "2006-01-02"

// ParseTime attempts to parse a JSON time string from Salesforce.  It will attempt
// to parse the time using RFC 3339, then Salesforce DateTime format and lastly Salesforce
// Date format.
func ParseTime(salesforceTime string) (time.Time, error) {
	var date time.Time
	var err error
	if date, err = time.Parse(time.RFC3339, salesforceTime); err == nil {
		return date, nil
	} else if date, err = time.Parse(SalesforceDateTime, salesforceTime); err == nil {
		return date, nil
	} else if date, err = time.Parse(SalesforceDate, salesforceTime); err == nil {
		return date, nil
	} else {
		return time.Time{}, err
	}
}
