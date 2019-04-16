package collections

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/g8rswimmer/goforce"
	"github.com/g8rswimmer/goforce/session"
)

func Test_collection_send(t *testing.T) {
	type fields struct {
		method   string
		endpoint string
		values   *url.Values
		body     io.Reader
	}
	type args struct {
		session session.ServiceFormatter
		value   interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Collection get with values",
			fields: fields{
				method:   http.MethodGet,
				endpoint: "/some/cool/endpoint",
				values: &url.Values{
					"one": []string{"this is fun"},
					"two": []string{"whatever,"},
				},
				body: nil,
			},
			args: args{
				session: &mockSessionFormatter{
					url: "something.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {

						if strings.HasPrefix(req.URL.String(), "something.com/some/cool/endpoint") == false {
							return &http.Response{
								StatusCode: 500,
								Status:     "Bad URL: " + req.URL.String(),
								Body:       ioutil.NopCloser(strings.NewReader("resp")),
								Header:     make(http.Header),
							}
						}

						if req.Method != http.MethodGet {
							return &http.Response{
								StatusCode: 500,
								Status:     "Bad Method",
								Body:       ioutil.NopCloser(strings.NewReader("resp")),
								Header:     make(http.Header),
							}
						}

						values := req.URL.Query()
						if _, ok := values["one"]; ok == false {

							return &http.Response{
								StatusCode: 500,
								Status:     "No one value",
								Body:       ioutil.NopCloser(strings.NewReader("resp")),
								Header:     make(http.Header),
							}
						}
						if _, ok := values["two"]; ok == false {

							return &http.Response{
								StatusCode: 500,
								Status:     "No two value",
								Body:       ioutil.NopCloser(strings.NewReader("resp")),
								Header:     make(http.Header),
							}
						}
						resp := `
							{
								"message" : "Email: invalid email address: Not a real email address",
								"errorCode" : "INVALID_EMAIL_ADDRESS",
								"fields" : [ "Email" ]
							}`

						return &http.Response{
							StatusCode: http.StatusOK,
							Status:     "Some Status",
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
				value: &goforce.Error{},
			},
			wantErr: false,
		},
		{
			name: "Collection with errors",
			fields: fields{
				method:   http.MethodGet,
				endpoint: "/some/cool/endpoint",
				values: &url.Values{
					"one": []string{"this is fun"},
					"two": []string{"whatever,"},
				},
				body: nil,
			},
			args: args{
				session: &mockSessionFormatter{
					url: "something.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {

						if strings.HasPrefix(req.URL.String(), "something.com/some/cool/endpoint") == false {
							return &http.Response{
								StatusCode: 500,
								Status:     "Bad URL: " + req.URL.String(),
								Body:       ioutil.NopCloser(strings.NewReader("resp")),
								Header:     make(http.Header),
							}
						}

						if req.Method != http.MethodGet {
							return &http.Response{
								StatusCode: 500,
								Status:     "Bad Method",
								Body:       ioutil.NopCloser(strings.NewReader("resp")),
								Header:     make(http.Header),
							}
						}

						values := req.URL.Query()
						if _, ok := values["one"]; ok == false {

							return &http.Response{
								StatusCode: 500,
								Status:     "No one value",
								Body:       ioutil.NopCloser(strings.NewReader("resp")),
								Header:     make(http.Header),
							}
						}
						if _, ok := values["two"]; ok == false {

							return &http.Response{
								StatusCode: 500,
								Status:     "No two value",
								Body:       ioutil.NopCloser(strings.NewReader("resp")),
								Header:     make(http.Header),
							}
						}
						resp := `
						[
							{
								"message" : "Email: invalid email address: Not a real email address",
								"errorCode" : "INVALID_EMAIL_ADDRESS",
								"fields" : [ "Email" ]
							}
						]`

						return &http.Response{
							StatusCode: http.StatusConflict,
							Status:     "Some Status",
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
				value: &goforce.Error{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &collection{
				method:   tt.fields.method,
				endpoint: tt.fields.endpoint,
				values:   tt.fields.values,
				body:     tt.fields.body,
			}
			if err := c.send(tt.args.session, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("collection.send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
