package soql

import "github.com/TuSKan/go-sfdc"

// QueryRecord is the result of the SOQL record.  If
// the query statement contains an inner query, there
// will be a group of subresults.
type QueryRecord struct {
	record     *sfdc.Record
	subresults map[string]*QueryResult
}

func newQueryRecord(jsonMap map[string]interface{}, resource *Resource) (*QueryRecord, error) {
	rec, err := sfdc.RecordFromJSONMap(jsonMap)
	if err != nil {
		return nil, err
	}
	subresults := make(map[string]*QueryResult)
	for k, v := range jsonMap {
		if sub, has := v.(map[string]interface{}); has {
			if isSubQuery(sub) {
				resp, err := newQueryResponseJSON(sub)
				if err != nil {
					return nil, err
				}
				result, err := newQueryResult(resp, resource)
				if err != nil {
					return nil, err
				}
				subresults[k] = result
			}
		}
	}
	qr := &QueryRecord{
		record:     rec,
		subresults: subresults,
	}
	return qr, nil
}

// Record returns the SOQL record.
func (rec *QueryRecord) Record() *sfdc.Record {
	return rec.record
}

// Subresults returns all of the inner query results.
func (rec *QueryRecord) Subresults() map[string]*QueryResult {
	return rec.subresults
}

// Subresult returns a specific inner query result.
func (rec *QueryRecord) Subresult(sub string) (*QueryResult, bool) {
	result, has := rec.subresults[sub]
	return result, has
}

func isSubQuery(jsonMap map[string]interface{}) bool {
	if _, has := jsonMap["totalSize"]; has == false {
		return false
	}
	if _, has := jsonMap["done"]; has == false {
		return false
	}
	if _, has := jsonMap["records"]; has == false {
		return false
	}
	return true
}
