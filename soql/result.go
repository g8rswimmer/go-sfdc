package soql

type QueryResult struct {
	response queryResponse
	records  []*QueryRecord
}

func newQueryResult(response queryResponse) (*QueryResult, error) {
	result := &QueryResult{
		response: response,
		records:  make([]*QueryRecord, len(response.Records)),
	}

	for idx, record := range response.Records {
		qr, err := newQueryRecord(record)
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
