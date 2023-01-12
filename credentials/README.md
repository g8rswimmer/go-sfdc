# Credentials
[back](../README.md)

To access `Salesforce APIs`, there needs to be authentication between the client and the org.  `go-sfdc` uses [OAuth 2.0](https://help.salesforce.com/articleView?id=remoteaccess_oauth_web_server_flow.htm&type=5) and this package provides the credentials needed to authenticate.

The user is able to use the `Providers` that are part of this package, or implement one of their own.  This allows for extendability beyond what is currently supported.

Currently, this package supports `grant type` of password OAuth flow.  The package may or may not be support other flows in the future.
## Examples
The following are some example(s) of creating credentials to be used when opening a session.
### Password
```go
creds := credentials.PasswordCredentials{
	URL:          "https://login.salesforce.com",
	Username:     "my.user@name.com",
	Password:     "greatpassword",
	ClientID:     "asdfnapodfnavppe",
	ClientSecret: "12312573857105",
}

config := sfdc.Configuration{
	Credentials: credentials.NewPasswordCredentials(creds),
	Client:      salesforceHTTPClient,
	Version:     44,
}
```

### JWT

```go
// read the key file
privateKeyFile, err := os.Open("/path/to/server.key")
if err != nil {
    panic()
}
pemfileinfo, _ := privateKeyFile.Stat()
var size int64 = pemfileinfo.Size()
pembytes := make([]byte, size)
buffer := bufio.NewReader(privateKeyFile)
_, err = buffer.Read(pembytes)
pemData := []byte(pembytes)
privateKeyFile.Close() // close file
signKey, err := jwt.ParseRSAPrivateKeyFromPEM(pemData)

// prepare the credentials
jwtCreds := credentials.JwtCredentials{
    URL: "https://login.salesforce.com",
    ClientId: "glfeoeirjgouwehf",
    ClientUsername: "my.user@name.com",
    ClientKey: signKey,
}

config := sfdc.Configuration{
    Credentials: credentials.NewJWTCredentials(jwtCreds),
    Client:      &http.Client{},
    Version:     44,
}
```
