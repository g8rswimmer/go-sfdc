package soql

import "errors"

type queryResponse struct {
	Done           bool                     `json:"done"`
	TotalSize      int                      `json:"totalSize"`
	NextRecordsURL string                   `json:"nextRecordsUrl"`
	Records        []map[string]interface{} `json:"records"`
}

func newQueryResponseJSON(jsonMap map[string]interface{}) (queryResponse, error) {
	response := queryResponse{}
	if d, has := jsonMap["done"]; has {
		if done, ok := d.(bool); ok {
			response.Done = done
		} else {
			return queryResponse{}, errors.New("query response: done is not a bool")
		}
	} else {
		return queryResponse{}, errors.New("query response: done is not present")
	}
	if ts, has := jsonMap["totalSize"]; has {
		if totalSize, ok := ts.(float64); ok {
			response.TotalSize = int(totalSize)
		} else {
			return queryResponse{}, errors.New("query response: totalSize is not a number")
		}
	} else {
		return queryResponse{}, errors.New("query response: totalSize is not present")
	}
	if nru, has := jsonMap["nextRecordsUrl"]; has {
		if nextRecordsURL, ok := nru.(string); ok {
			response.NextRecordsURL = nextRecordsURL
		} else {
			return queryResponse{}, errors.New("query response: nextRecordsUrl is not a string")
		}
	}
	if r, has := jsonMap["records"]; has {
		if records, ok := r.([]map[string]interface{}); ok {
			response.Records = records
		} else {
			return queryResponse{}, errors.New("query response: records is not an array")
		}
	} else {
		return queryResponse{}, errors.New("query response: records is not present")
	}
	return response, nil
}
