package goforce

import (
	"time"
)

const SalesforceDateTime = "2006-01-02T15:04:05.000+0000"
const SalesforceDate = "2006-01-02"

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
