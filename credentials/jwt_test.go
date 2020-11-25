package credentials

import (
	"github.com/dgrijalva/jwt-go"
	"io"
	"net/url"
	"strings"
	"testing"
)

const SampleKey = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAwCPdYprMz3AMh4CwK8fPdArUL63RMVYoXXYfzFdluW5XYE5m\n0a5PuNpMoc33i7+JYGOCS1T+ZhoAM2AHO3/BbC2sB5qNNj48ToR7RADgy+pKyaUa\niks4hWXI2fqzAZR11xFEMkCKl0S7Zn4t/oZkFlXbgI+fxt+ab8+9rXa770pL7yCO\nlh5HLIQ1VUPWJN7JeBiKfSnBowGLuelQ8ot7YJmEhohBUN++5ZrfSqPedeLlDYPV\nZLYlEaZE6Xtg0lI+prsJ6wiv0IlTwH7yYECc2XE8MjyWAlNEoObK6kbfD0oIQqU1\noSXkRCp21MHoJ9ZTFJvd+2GArbYTzz0KL/r7TQIDAQABAoIBAQCfXmAnhIyy1pad\n4gC+H5qT/tNmxL6KNJOAihTv8eH/P2WcDQu9id64TeFYKDXWpUU2PPN6toHYgGKA\nOntlP59Ysj1JhUjxoAd3fO2dRzkuCiSEQrzTznaQNw+0tfu6KMDhZYHySJRryefC\nqJBP2Hq2B/rsFLULSLaZXW9PrPdPDxijnq+Mnok8t+1F1LRkhdXAiTAAryLT7V4I\neK6uMd0bHK776dQy7A0hR55B5NOW/1U5iYHhNMNCw31Tct9Ula5Dt5U+oe69xMd5\ntog7UaESglsovusN65GpXjwsBUN/a4qYXXUEa3ZWhDsuH3b2ekcMRgfsx5QLSDqH\n5nFtX3kFAoGBAPLa4oBmLerzZ68GMU+uKQscN/C8o671UTS7jCbtWcBMUZOGrN3q\n+smpB0YB9W4kALSxYM4LsTQ4n8qkfJz1vlMu6iATcPnG+KlEhVITUHXB/ek5aWfZ\nN1uZDGUlg0sgWSlubnNs5xl6J8tYGRrk84g5i6QCSYxesoJ+M/1P7yozAoGBAMqK\nP/PYkbJWq/gh3KcrbbiQhjz6EoPlypcjBzdfQnPamJ94voh2YYNETlDXTSr5bATP\n+dooSaw6lkoDIzZg9IZrq9FDOwXjHptpakpIkYKXKxLVcBl6PrD/hv7jznawdLrP\nyWr9nkqIVHvJxMGvjg7ONgJhCuCHmecrO50p4sR/AoGAa/8aqq7FzK3hddvzIdP5\nPI+X8N5yi+Nb8W9VrBnwx6sou8owJZ/RVsxsB53nXstz5ObcfcSFUQu9Q4hSQhqm\nQKekRg9fNjRdcCiggRdFuJhEKer2DNBz5a/x6yj7cfU4sUwCoiHTw2inOa47u9IE\n2pd8mbrKqjmSeKVWyVc6rDECgYAlrp0BYByTQn7SNnKYA4NxYCopdBk3wuvzPIge\nLDHv3g6hNNS2DNhNlMrBTZ1EzozjRFJm3TH/whKuCHFnr5gu3h9kWo7DpKLQJUeq\nNGAmHLvd0CoAA3dgdNoH2BhUirXc/8WoizEFCuI0+bAKnP/gD0uLG8TrSy8+DBQW\nRHG1PwKBgBAsnnjH4KplKrzfMycTczHEM1pll/wWBe38TbA7YjrOYLGJkac0UVVQ\nGqhoj3JfpSlWoUbMrOlyY7FlIptmj71P+xPNThKTcc42KMzYJCPhdMllXWktwWKo\nhQZsXUbv/2dzOsyQZWcWM/k+kVArS4+Q3eStBNxaDl0aNQC9CUUg\n-----END RSA PRIVATE KEY-----"



func mockJWTRetriveReader(creds JwtCredentials) io.Reader {
	form := url.Values{}
	form.Add("grant_type", string(jwtGrantType))


	return strings.NewReader(form.Encode())
}
func Test_jwtProvider_Retrieve(t *testing.T) {

	pemData := []byte(SampleKey)

	signKey, _ := jwt.ParseRSAPrivateKeyFromPEM(pemData)

	type fields struct {
		creds JwtCredentials
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Token Retriever",
			fields: fields{
				creds: JwtCredentials{
					URL:          "http://test.password.session",
					ClientId:     "12345",
					ClientUsername:     "myusername",
					ClientKey: signKey,
				},
			},
			want: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJodHRwOi8vdGVzdC5wYXNzd29yZC5zZXNzaW9uIiwiZXhwIjoxMjMyMTMsImlzcyI6IjEyMzQ1Iiwic3ViIjoibXl1c2VybmFtZSJ9.rMwIVIG9IEHuZ9O5Na3VZHCRTw0z7iBRugWPlZeq6VFjcrHimDvBxN91DQtQcOqULSORWvo0HGzEHo_LOnMopj6dyL-wxgEp3Fu9owv-69HNP_uYYfy0_W93-1BVlOgCGuF6NdP2JHb3EsUmEPf-DJcjYmiymUMrMV6vNhkbl3eb_V93xA4-sid4reyMyBDVznKqyrQOOiyEqEq3tF9IZKcu7Cr1KDJa2br8S9Rq_xwv6PZPm4Pm97Bpd5t95Fo6zyTIOLu-zaXRzwKBRF6StuXsisfXo66-jDnSs0R-fx_RaI6GRez900lEzqZK13ggsLmB3PqpNcZkWvKvExBSHg",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &jwtProvider{
				creds: tt.fields.creds,
			}
			token, err := provider.BuildClaimsToken(123213, provider.creds.URL, provider.creds.ClientId, provider.creds.ClientUsername)

			if err != nil && !tt.wantErr {
				t.Error()
			}

			if token != tt.want {
				t.Error("tokens do not match")
			}
		})
	}
}
