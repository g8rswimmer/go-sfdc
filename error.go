package goforce

import (
	"encoding/json"
	"errors"
)

// Error is the error structure defined by the Salesforce API.
type Error struct {
	ErrorCode string   `json:"errorCode"`
	Message   string   `json:"message"`
	Fields    []string `json:"fields"`
}

func (e *Error) UnmarshalJSON(data []byte) error {
	if e == nil {
		return errors.New("record: can't unmarshal to a nil struct")
	}

	var jsonMap map[string]interface{}
	err := json.Unmarshal(data, &jsonMap)
	if err != nil {
		return err
	}

	if code, ok := jsonMap["statusCode"]; ok {
		e.ErrorCode = code.(string)
	}
	if code, ok := jsonMap["errorCode"]; ok {
		e.ErrorCode = code.(string)
	}
	if message, ok := jsonMap["message"]; ok {
		e.Message = message.(string)
	}
	if fields, ok := jsonMap["fields"]; ok {
		if array, has := fields.([]interface{}); has {
			e.Fields = make([]string, len(array))
			for idx, element := range array {
				if field, ok := element.(string); ok {
					e.Fields[idx] = field
				}
			}
		}
	}
	return nil
}
