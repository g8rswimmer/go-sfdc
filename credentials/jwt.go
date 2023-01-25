package credentials

import (
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

const (
	JwtExpiration = 3 * time.Minute
)

type jwtProvider struct {
	creds JwtCredentials
}

type claims struct {
	jwt.RegisteredClaims
}

func (provider *jwtProvider) Retrieve() (io.Reader, error) {
	expirationTime := provider.GetAppropriateExpirationTime()
	tokenString, err := provider.BuildClaimsToken(expirationTime, provider.creds.URL, provider.creds.ClientId, provider.creds.ClientUsername)
	if err != nil {
		return nil, fmt.Errorf("jwtProvider.Retrieve() error: %w", err)
	}

	form := url.Values{}
	form.Add("grant_type", string(jwtGrantType))
	form.Add("assertion", tokenString)

	return strings.NewReader(form.Encode()), nil
}

func (provider *jwtProvider) URL() string {
	return provider.creds.URL
}

func (provider *jwtProvider) GetAppropriateExpirationTime() int64 {
	return time.Now().Add(JwtExpiration).Unix()
}

// builds the actual claims token required for authentication
func (provider *jwtProvider) BuildClaimsToken(expirationTime int64, url string, clientId string, clientUsername string) (string, error) {
	claims := &claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(time.Unix(expirationTime, 0)),
			Audience:  []string{url},  // "https://test.salesforce.com" || "https://login.salesforce.com"
			Issuer:    clientId,       // consumer key of the connected app, hardcoded
			Subject:   clientUsername, // username of the salesforce user, whose profile is added to the connected app
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, error := token.SignedString(provider.creds.ClientKey)
	return tokenString, error
}
