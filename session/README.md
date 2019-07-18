# Session
[back](../README.md)

The session is used to authenticate with `Salesforce` and retrieve the org's information.  The session will be used by the API packages to properly access the `Salesforce org`.  The [credentials](../credentials/README.md) will need to be will be used to properly create the session.

## Example
The following example demonstrates how to create a session.
```go
pwdCreds, err := credentials.NewPasswordCredentials(credentials.PasswordCredentials{
	URL:          "https://login.salesforce.com",
	Username:     "my.user@name.com",
	Password:     "greatpassword",
	ClientID:     "asdfnapodfnavppe",
	ClientSecret: "12312573857105",
})

if err != nil {
    fmt.Printf("error %v\n", err)
    return
}

config := sfdc.Configuration{
	Credentials: pwdCreds,
	Client:      http.DefaultClient,
	Version:     44,
}

session, err := session.Open(config)

if err != nil {
	fmt.Printf("Error %v\n", err)
	return
}

// access Salesforce APIs
```
