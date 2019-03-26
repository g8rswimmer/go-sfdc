package goforce

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestPasswordSessionRequest(t *testing.T) {

	scenarios := []struct {
		desc  string
		creds SessionPasswordCredentials
		err   error
	}{
		{
			desc: "Passing HTTP request",
			creds: SessionPasswordCredentials{
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
			creds: SessionPasswordCredentials{
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

		request, err := passwordSessionRequest(scenario.creds)

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

				if request.URL.String() != scenario.creds.URL+oauthService {
					t.Errorf("%s URL not matching %s :: %s", scenario.desc, scenario.creds.URL+oauthService, request.URL.String())
				}

				buf, err := ioutil.ReadAll(request.Body)
				request.Body.Close()
				if err != nil {
					t.Fatal(err.Error())
				}
				form := url.Values{}
				form.Add("grant_type", string(passwordGrantType))
				form.Add("username", scenario.creds.Username)
				form.Add("password", scenario.creds.Password)
				form.Add("client_id", scenario.creds.ClientID)
				form.Add("client_secret", scenario.creds.ClientSecret)

				if form.Encode() != string(buf) {
					t.Errorf("%s Form data %s :: %s", scenario.desc, string(buf), form.Encode())
				}
			}
		}
	}

}

func TestPasswordSessionResponse(t *testing.T) {
	t.Fatal("Not yet implemented")
}

func TestNewPasswordSession(t *testing.T) {
	t.Fatal("Not yet implemented")
}
