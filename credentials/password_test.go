package credentials

import (
	"io"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func mockPasswordRetriveReader(creds PasswordCredentails) io.Reader {
	form := url.Values{}
	form.Add("grant_type", string(passwordGrantType))
	form.Add("username", creds.Username)
	form.Add("password", creds.Password)
	form.Add("client_id", creds.ClientID)
	form.Add("client_secret", creds.ClientSecret)

	return strings.NewReader(form.Encode())
}
func Test_passwordProvider_Retrieve(t *testing.T) {
	type fields struct {
		creds PasswordCredentails
	}
	tests := []struct {
		name    string
		fields  fields
		want    io.Reader
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Password Retriever",
			fields: fields{
				creds: PasswordCredentails{
					URL:          "http://test.password.session",
					Username:     "myusername",
					Password:     "12345",
					ClientID:     "some client id",
					ClientSecret: "shhhh its a secret",
				},
			},
			want: mockPasswordRetriveReader(PasswordCredentails{
				URL:          "http://test.password.session",
				Username:     "myusername",
				Password:     "12345",
				ClientID:     "some client id",
				ClientSecret: "shhhh its a secret",
			}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &passwordProvider{
				creds: tt.fields.creds,
			}
			got, err := provider.Retrieve()
			if (err != nil) != tt.wantErr {
				t.Errorf("passwordProvider.Retrieve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("passwordProvider.Retrieve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_passwordProvider_URL(t *testing.T) {
	type fields struct {
		creds PasswordCredentails
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Password URL",
			fields: fields{
				creds: PasswordCredentails{
					URL:          "http://test.password.session",
					Username:     "myusername",
					Password:     "12345",
					ClientID:     "some client id",
					ClientSecret: "shhhh its a secret",
				},
			},
			want: "http://test.password.session",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &passwordProvider{
				creds: tt.fields.creds,
			}
			if got := provider.URL(); got != tt.want {
				t.Errorf("passwordProvider.URL() = %v, want %v", got, tt.want)
			}
		})
	}
}
