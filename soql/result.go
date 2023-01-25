package soql

import "errors"

// QueryResult is returned from the SOQL query.  This will
// allow for retrieving all of the records and query the
// next round of records if available.
type QueryResult struct {
	response queryResponse
	records  []*QueryRecord
	resource *Resource
}

func newQueryResult(response queryResponse, resource *Resource) (*QueryResult, error) {
	result := &QueryResult{
		response: response,
		records:  make([]*QueryRecord, len(response.Records)),
		resource: resource,
	}

	for idx, record := range response.Records {
		qr, err := newQueryRecord(record, resource)
		if err != nil {
			return nil, err
		}
		result.records[idx] = qr
	}
	return result, nil
}

// Done will indicate if the result does not contain any more records.
func (result *QueryResult) Done() bool {
	return result.response.Done
}

// TotalSize is the total size of the query result.  This may
// or may not be the size of the records in the result.
func (result *QueryResult) TotalSize() int {
	return result.response.TotalSize
}

// MoreRecords will indicate if the remaining records require another
// Saleforce service callout.
func (result *QueryResult) MoreRecords() bool {
	return result.response.NextRecordsURL != ""
}

// Records returns the records from the query request.
func (result *QueryResult) Records() []*QueryRecord {
	return result.records
}

// Next will query the next set of records.
func (result *QueryResult) Next() (*QueryResult, error) {
	if !result.MoreRecords() {
		return nil, errors.New("soql query result: no more records to query")
	}
	return result.resource.next(result.response.NextRecordsURL)
}
