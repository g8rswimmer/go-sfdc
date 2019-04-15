package collections

import (
	"io"
	"net/http"
	"net/url"

	"github.com/g8rswimmer/goforce/session"
)

const endpoint = "/composite/sobjects"

type collection struct {
	method   string
	endpoint string
	values   *url.Values
	body     io.Reader
}

func (c *collection) send(session session.ServiceFormatter) (*http.Response, error) {
	collectionURL := session.ServiceURL() + c.endpoint
	if c.values != nil {
		collectionURL += "?" + c.values.Encode()
	}
	request, err := http.NewRequest(c.method, collectionURL, c.body)
	if err != nil {
		return nil, err
	}
	return session.Client().Do(request)
}
