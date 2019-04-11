package soql

import "github.com/g8rswimmer/goforce"

type QueryRecord struct {
	record     *goforce.Record
	subresults map[string]*QueryResult
}

func newQueryRecord(jsonMap map[string]interface{}) (*QueryRecord, error) {
	rec, err := goforce.RecordFromJSONMap(jsonMap)
	if err != nil {
		return nil, err
	}
	subresults := make(map[string]*QueryResult)
	for k, v := range jsonMap {
		if sub, has := v.(map[string]interface{}); has {
			if k != goforce.RecordAttributes {
				resp, err := newQueryResponseJSON(sub)
				if err != nil {
					return nil, err
				}
				result, err := newQueryResult(resp)
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
func (rec *QueryRecord) Record() *goforce.Record {
	return rec.record
}
func (rec *QueryRecord) Subresults() map[string]*QueryResult {
	return rec.subresults
}
func (rec *QueryRecord) Subresult(sub string) (*QueryResult, bool) {
	result, has := rec.subresults[sub]
	return result, has
}
