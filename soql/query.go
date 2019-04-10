package soql

import (
	"github.com/g8rswimmer/goforce"
)

type queryResponse struct {
	Done           bool                     `json:"done"`
	TotalSize      int                      `json:"totalSize"`
	NextRecordsURL string                   `json:"nextRecordsUrl"`
	Records        []map[string]interface{} `json:"records"`
}

type QueryResult struct {
	response queryResponse
	records  []QueryRecord
}

type QueryRecord struct {
	record     *goforce.Record
	subresults map[string]*QueryResult
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
func (result *QueryResult) Records() []QueryRecord {
	return result.records
}

func (rec *QueryRecord) Record() *goforce.Record {
	return rec.record
}
func (rec *QueryRecord) Subresults(sub string) (*QueryResult, bool) {
	result, has := rec.subresults[sub]
	return result, has
}
