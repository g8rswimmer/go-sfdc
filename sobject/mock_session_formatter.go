package sobject

import "net/http"

type mockSessionFormatter struct {
	url    string
	client *http.Client
}

func (mock *mockSessionFormatter) ServiceURL() string {
	return mock.url
}
func (mock *mockSessionFormatter) AuthorizationHeader(*http.Request) {}

func (mock *mockSessionFormatter) Client() *http.Client {
	return mock.client
}
