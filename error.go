package goforce

type Error struct {
	ErrorCode string   `json:"errorCode"`
	Message   string   `json:"message"`
	Fields    []string `json:"fields"`
}
