# goforce
This is a `golang` library for interfacing with `Salesforce` APIs.

## Sessions
Before calling a `Salesforce` API, a session will need to be created.  This can be done with `OAuth 2.0`.

### Example
```
	
    creds := credentials.PasswordCredentails{
		URL:          "https://test.salesforce.com",
		Username:     "user@somename.com",
		Password:     "myfunpassword",
		ClientID:     "Some ACSII stuff",
		ClientSecret: "Some other numbers",
	}

	config := goforce.Configuration{
        Credentials: crendentials.NewPasswordCredentials(creds),
        Client: myHttpClient,
    }

    session, err := session.NewPasswordSession(config)

	if err != nil {
        // handle the session error...
        return
	}

    // can start accessing Salesforce APIs

```
