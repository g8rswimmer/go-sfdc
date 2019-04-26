package credentials

import (
	"io"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestNewPasswordCredentials(t *testing.T) {
	type args struct {
		creds PasswordCredentails
	}
	tests := []struct {
		name    string
		args    args
		want    *Credentials
		wantErr bool
	}{
		{
			name: "Password Credentials",
			args: args{
				creds: PasswordCredentails{
					URL:          "http://test.password.session",
					Username:     "myusername",
					Password:     "12345",
					ClientID:     "some client id",
					ClientSecret: "shhhh its a secret",
				},
			},
			want: &Credentials{
				provider: &passwordProvider{
					creds: PasswordCredentails{
						URL:          "http://test.password.session",
						Username:     "myusername",
						Password:     "12345",
						ClientID:     "some client id",
						ClientSecret: "shhhh its a secret",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "No URL",
			args: args{
				creds: PasswordCredentails{
					URL:          "",
					Username:     "myusername",
					Password:     "12345",
					ClientID:     "some client id",
					ClientSecret: "shhhh its a secret",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "No Username",
			args: args{
				creds: PasswordCredentails{
					URL:          "http://test.password.session",
					Username:     "",
					Password:     "12345",
					ClientID:     "some client id",
					ClientSecret: "shhhh its a secret",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "No password",
			args: args{
				creds: PasswordCredentails{
					URL:          "http://test.password.session",
					Username:     "myusername",
					Password:     "",
					ClientID:     "some client id",
					ClientSecret: "shhhh its a secret",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "No client ID",
			args: args{
				creds: PasswordCredentails{
					URL:          "http://test.password.session",
					Username:     "myusername",
					Password:     "12345",
					ClientID:     "",
					ClientSecret: "shhhh its a secret",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "No client secret",
			args: args{
				creds: PasswordCredentails{
					URL:          "http://test.password.session",
					Username:     "myusername",
					Password:     "12345",
					ClientID:     "some client id",
					ClientSecret: "",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := NewPasswordCredentials(tt.args.creds)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPasswordCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPasswordCredentials() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCredentials(t *testing.T) {
	type args struct {
		provider Provider
	}
	tests := []struct {
		name    string
		args    args
		want    *Credentials
		wantErr bool
	}{
		{
			name: "New Credentials",
			args: args{
				provider: &passwordProvider{
					creds: PasswordCredentails{
						URL:          "http://test.password.session",
						Username:     "myusername",
						Password:     "12345",
						ClientID:     "some client id",
						ClientSecret: "shhhh its a secret",
					},
				},
			},
			want: &Credentials{
				provider: &passwordProvider{
					creds: PasswordCredentails{
						URL:          "http://test.password.session",
						Username:     "myusername",
						Password:     "12345",
						ClientID:     "some client id",
						ClientSecret: "shhhh its a secret",
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "No Provider",
			args:    args{},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCredentials(tt.args.provider)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCredentials() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCredentials_URL(t *testing.T) {
	type fields struct {
		provider Provider
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
		{
			name: "Credential URL",
			fields: fields{
				provider: &passwordProvider{
					creds: PasswordCredentails{
						URL:          "http://test.password.session",
						Username:     "myusername",
						Password:     "12345",
						ClientID:     "some client id",
						ClientSecret: "shhhh its a secret",
					},
				},
			},
			want: "http://test.password.session",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creds := &Credentials{
				provider: tt.fields.provider,
			}
			if got := creds.URL(); got != tt.want {
				t.Errorf("Credentials.URL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockCredentialsRetriveReader(creds PasswordCredentails) io.Reader {
	form := url.Values{}
	form.Add("grant_type", string(passwordGrantType))
	form.Add("username", creds.Username)
	form.Add("password", creds.Password)
	form.Add("client_id", creds.ClientID)
	form.Add("client_secret", creds.ClientSecret)

	return strings.NewReader(form.Encode())
}

func TestCredentials_Retrieve(t *testing.T) {
	type fields struct {
		provider Provider
	}
	tests := []struct {
		name    string
		fields  fields
		want    io.Reader
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Credential Retrieve",
			fields: fields{
				provider: &passwordProvider{
					creds: PasswordCredentails{
						URL:          "http://test.password.session",
						Username:     "myusername",
						Password:     "12345",
						ClientID:     "some client id",
						ClientSecret: "shhhh its a secret",
					},
				},
			},
			want: mockCredentialsRetriveReader(PasswordCredentails{
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
			creds := &Credentials{
				provider: tt.fields.provider,
			}
			got, err := creds.Retrieve()
			if (err != nil) != tt.wantErr {
				t.Errorf("Credentials.Retrieve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Credentials.Retrieve() = %v, want %v", got, tt.want)
			}
		})
	}
}
