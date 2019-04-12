package soql

import "errors"

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
func (result *QueryResult) Done() bool {
	return result.response.Done
}
func (result *QueryResult) TotalSize() int {
	return result.response.TotalSize
}
func (result *QueryResult) MoreRecords() bool {
	return result.response.NextRecordsURL != ""
}
func (result *QueryResult) Records() []*QueryRecord {
	return result.records
}
func (result *QueryResult) Next() (*QueryResult, error) {
	if result.MoreRecords() == false {
		return nil, errors.New("soql query result: no more records to query")
	}
	return result.resource.next(result.response.NextRecordsURL)
}
