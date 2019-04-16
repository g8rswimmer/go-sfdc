package goforce

// Error is the error structure defined by the Salesforce API.
type Error struct {
	StatusCode string   `json:"statusCode"`
	ErrorCode  string   `json:"errorCode"`
	Message    string   `json:"message"`
	Fields     []string `json:"fields"`
}
