package common

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"

	"github.com/g8rswimmer/go-sfdc/credentials"
)

// Payload is the structure used in the examples
type Payload struct {
	Credentials credentials.PasswordCredentials `json:"credentials"`
	Version     int                             `json:"version"`
	DML         *DML                            `json:"dml,omitempty"`
}

// RetrievePayload will return the payload if passed by file or object.
func RetrievePayload() Payload {
	filePtr := flag.String("file", "", "json payload file")
	objPtr := flag.String("object", "", "json payload object")
	flag.Parse()
	var data []byte
	var err error

	switch {
	case len(*filePtr) > 0:
		data, err = ioutil.ReadFile(*filePtr)
		Error(err)
	case len(*objPtr) > 0:
		data = []byte(*objPtr)
	default:
		Error(errors.New("An file or object argument must be present"))
	}

	var payload Payload
	err = json.Unmarshal(data, &payload)
	Error(err)
	return payload
}
