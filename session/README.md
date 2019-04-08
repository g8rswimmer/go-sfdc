# Session
[back](../README.md)

The session is used to authenticate with `Salesforce` and retrieve the org's information.  The session will be used by the API packages to properly access the `Salesforce org`.  The [credentials](../credentials/README.md) will need to be will be used to properly create the session.

## Example
The following example demonstrates how to create a session.
```go
creds := credentials.PasswordCredentails{
  URL:          "https://login.salesforce.com",
  Username:     "my.user@name.com",
  Password:     "greatpassword",
  ClientID:     "asdfnapodfnavppe",
  ClientSecret: "12312573857105",
}

config := goforce.Configuration{
  Credentials: credentials.NewPasswordCredentials(creds),
  Client:      http.DefaultClient,
  Version:     44,
}

session, err := session.Open(config)

if err != nil {
  fmt.Printf("Error %s", err.Error())
  fmt.Println()
  return
}

// access Salesforce APIs
```
