package session

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/g8rswimmer/goforce"

	"github.com/g8rswimmer/goforce/credentials"
)

func TestPasswordSessionRequest(t *testing.T) {

	scenarios := []struct {
		desc  string
		creds credentials.PasswordCredentails
		err   error
	}{
		{
			desc: "Passing HTTP request",
			creds: credentials.PasswordCredentails{
				URL:          "http://test.password.session",
				Username:     "myusername",
				Password:     "12345",
				ClientID:     "some client id",
				ClientSecret: "shhhh its a secret",
			},
			err: nil,
		},
		{
			desc: "Bad URL",
			creds: credentials.PasswordCredentails{
				URL:          "123://something.com",
				Username:     "myusername",
				Password:     "12345",
				ClientID:     "some client id",
				ClientSecret: "shhhh its a secret",
			},
			err: errors.New("parse 123://something.com/services/oauth2/token: first path segment in URL cannot contain colon"),
		},
	}

	for _, scenario := range scenarios {

		passwordCreds := credentials.NewPasswordCredentials(scenario.creds)
		request, err := passwordSessionRequest(passwordCreds)

		if err != nil && scenario.err == nil {
			t.Errorf("%s Error was not expected %s", scenario.desc, err.Error())
		} else if err == nil && scenario.err != nil {
			t.Errorf("%s Error was expected %s", scenario.desc, scenario.err.Error())
		} else {
			if err != nil {
				if err.Error() != scenario.err.Error() {
					t.Errorf("%s Error %s :: %s", scenario.desc, err.Error(), scenario.err.Error())
				}
			} else {
				if request.Method != http.MethodPost {
					t.Errorf("%s HTTP request method needs to be POST not %s", scenario.desc, request.Method)
				}

				if request.URL.String() != scenario.creds.URL+oauthEndppoint {
					t.Errorf("%s URL not matching %s :: %s", scenario.desc, scenario.creds.URL+oauthEndppoint, request.URL.String())
				}

				buf, err := ioutil.ReadAll(request.Body)
				request.Body.Close()
				if err != nil {
					t.Fatal(err.Error())
				}
				reader, err := passwordCreds.Retrieve()
				if err != nil {
					t.Fatal(err.Error())
				}
				body, err := ioutil.ReadAll(reader)
				if err != nil {
					t.Fatal(err.Error())
				}

				if string(body) != string(buf) {
					t.Errorf("%s Form data %s :: %s", scenario.desc, string(buf), string(body))
				}
			}
		}
	}

}

func TestPasswordSessionResponse(t *testing.T) {
	scenarios := []struct {
		desc     string
		url      string
		client   *http.Client
		response *sessionPasswordResponse
		err      error
	}{
		{
			desc: "Passing Response",
			url:  "http://example.com/foo",
			client: mockHTTPClient(func(req *http.Request) *http.Response {
				resp := `
				{
					"access_token": "token",
					"instance_url": "https://some.salesforce.instance.com",
					"id": "https://test.salesforce.com/id/123456789",
					"token_type": "Bearer",
					"issued_at": "1553568410028",
					"signature": "hello"
				}`

				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(resp)),
					Header:     make(http.Header),
				}
			}),
			response: &sessionPasswordResponse{
				AccessToken: "token",
				InstanceURL: "https://some.salesforce.instance.com",
				ID:          "https://test.salesforce.com/id/123456789",
				TokenType:   "Bearer",
				IssuedAt:    "1553568410028",
				Signature:   "hello",
			},
			err: nil,
		},
		{
			desc: "Failed Response",
			url:  "http://example.com/foo",
			client: mockHTTPClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     "Some status",
					Body:       ioutil.NopCloser(strings.NewReader("")),
					Header:     make(http.Header),
				}
			}),
			response: &sessionPasswordResponse{},
			err:      fmt.Errorf("session response error: %d %s", http.StatusInternalServerError, "Some status"),
		},
		{
			desc: "Response Decode Error",
			url:  "http://example.com/foo",
			client: mockHTTPClient(func(req *http.Request) *http.Response {
				resp := `
				{
					"access_token": "token",
					"instance_url": "https://some.salesforce.instance.com",
					"id": "https://test.salesforce.com/id/123456789",
					"token_type": "Bearer",
					"issued_at": "1553568410028",
					"signature": "hello",
				}`

				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(resp)),
					Header:     make(http.Header),
				}
			}),
			response: &sessionPasswordResponse{},
			err:      errors.New("invalid character '}' looking for beginning of object key string"),
		},
	}

	for _, scenario := range scenarios {

		request, err := http.NewRequest(http.MethodPost, scenario.url, nil)

		if err != nil {
			t.Fatal(err.Error())
		}

		response, err := passwordSessionResponse(request, scenario.client)

		if err != nil && scenario.err == nil {
			t.Errorf("%s Error was not expected %s", scenario.desc, err.Error())
		} else if err == nil && scenario.err != nil {
			t.Errorf("%s Error was expected %s", scenario.desc, scenario.err.Error())
		} else {
			if err != nil {
				if err.Error() != scenario.err.Error() {
					t.Errorf("%s Error %s :: %s", scenario.desc, err.Error(), scenario.err.Error())
				}
			} else {
				if response.AccessToken != scenario.response.AccessToken {
					t.Errorf("%s Access Tokens %s %s", scenario.desc, scenario.response.AccessToken, response.AccessToken)
				}

				if response.InstanceURL != scenario.response.InstanceURL {
					t.Errorf("%s Instance URL %s %s", scenario.desc, scenario.response.InstanceURL, response.InstanceURL)
				}

				if response.ID != scenario.response.ID {
					t.Errorf("%s ID %s %s", scenario.desc, scenario.response.ID, response.ID)
				}

				if response.TokenType != scenario.response.TokenType {
					t.Errorf("%s Token Type %s %s", scenario.desc, scenario.response.TokenType, response.TokenType)
				}

				if response.IssuedAt != scenario.response.IssuedAt {
					t.Errorf("%s Issued At %s %s", scenario.desc, scenario.response.IssuedAt, response.IssuedAt)
				}

				if response.Signature != scenario.response.Signature {
					t.Errorf("%s Signature %s %s", scenario.desc, scenario.response.Signature, response.Signature)
				}

			}
		}

	}
}

func TestNewPasswordSession(t *testing.T) {
	scenarios := []struct {
		desc    string
		config  goforce.Configuration
		session *Session
		err     error
	}{
		{
			desc: "Passing",
			config: goforce.Configuration{
				Credentials: credentials.NewPasswordCredentials(credentials.PasswordCredentails{
					URL:          "http://test.password.session",
					Username:     "myusername",
					Password:     "12345",
					ClientID:     "some client id",
					ClientSecret: "shhhh its a secret",
				}),
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					resp := `
					{
						"access_token": "token",
						"instance_url": "https://some.salesforce.instance.com",
						"id": "https://test.salesforce.com/id/123456789",
						"token_type": "Bearer",
						"issued_at": "1553568410028",
						"signature": "hello"
					}`

					return &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(strings.NewReader(resp)),
						Header:     make(http.Header),
					}
				}),
			},
			session: &Session{
				response: &sessionPasswordResponse{
					AccessToken: "token",
					InstanceURL: "https://some.salesforce.instance.com",
					ID:          "https://test.salesforce.com/id/123456789",
					TokenType:   "Bearer",
					IssuedAt:    "1553568410028",
					Signature:   "hello",
				},
			},
			err: nil,
		},
		{
			desc: "Error Request",
			config: goforce.Configuration{
				Credentials: credentials.NewPasswordCredentials(credentials.PasswordCredentails{
					URL:          "123://test.password.session",
					Username:     "myusername",
					Password:     "12345",
					ClientID:     "some client id",
					ClientSecret: "shhhh its a secret",
				}),
				Client: nil,
			},
			session: nil,
			err:     errors.New("parse 123://test.password.session/services/oauth2/token: first path segment in URL cannot contain colon"),
		},
		{
			desc: "Error Response",
			config: goforce.Configuration{
				Credentials: credentials.NewPasswordCredentials(credentials.PasswordCredentails{
					URL:          "http://test.password.session",
					Username:     "myusername",
					Password:     "12345",
					ClientID:     "some client id",
					ClientSecret: "shhhh its a secret",
				}),
				Client: mockHTTPClient(func(req *http.Request) *http.Response {

					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Status:     "Some status",
						Body:       ioutil.NopCloser(strings.NewReader("")),
						Header:     make(http.Header),
					}
				}),
			},
			session: nil,
			err:     fmt.Errorf("session response error: %d %s", http.StatusInternalServerError, "Some status"),
		},
	}

	for _, scenario := range scenarios {

		session, err := NewPasswordSession(scenario.config)

		if err != nil && scenario.err == nil {
			t.Errorf("%s Error was not expected %s", scenario.desc, err.Error())
		} else if err == nil && scenario.err != nil {
			t.Errorf("%s Error was expected %s", scenario.desc, scenario.err.Error())
		} else {
			if err != nil {
				if err.Error() != scenario.err.Error() {
					t.Errorf("%s Error %s :: %s", scenario.desc, err.Error(), scenario.err.Error())
				}
			} else {
				if session.response.AccessToken != scenario.session.response.AccessToken {
					t.Errorf("%s Access Tokens %s %s", scenario.desc, scenario.session.response.AccessToken, session.response.AccessToken)
				}

				if session.response.InstanceURL != scenario.session.response.InstanceURL {
					t.Errorf("%s Instance URL %s %s", scenario.desc, scenario.session.response.InstanceURL, session.response.InstanceURL)
				}

				if session.response.ID != scenario.session.response.ID {
					t.Errorf("%s ID %s %s", scenario.desc, scenario.session.response.ID, session.response.ID)
				}

				if session.response.TokenType != scenario.session.response.TokenType {
					t.Errorf("%s Token Type %s %s", scenario.desc, scenario.session.response.TokenType, session.response.TokenType)
				}

				if session.response.IssuedAt != scenario.session.response.IssuedAt {
					t.Errorf("%s Issued At %s %s", scenario.desc, scenario.session.response.IssuedAt, session.response.IssuedAt)
				}

				if session.response.Signature != scenario.session.response.Signature {
					t.Errorf("%s Signature %s %s", scenario.desc, scenario.session.response.Signature, session.response.Signature)
				}

			}
		}

	}
}
