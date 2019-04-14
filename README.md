# goforce
This is a `golang` library for interfacing with `Salesforce` APIs.

## Usage
To use this library, the following will need to be done.
* Create `Salesforce` [credentials](./credentials/README.md) to properly authenticate with the `Salesforce org`
* Configure
* Open a [session](./session/README.md)
* Use the `APIs`
  - [SObject APIs](./sobject/README.md)
  - [SOQL APIs](./soql/README.md)

## Configuration
The configuration defines several parameters that can be used by the library.  The configuration is used per [session](./session/README.md).
* `Credentials` - this is an implementation of the `credentials.Provider` interface
* `Client` - the HTTP client used by the `APIs`
* `Version` - is the `Salesforce` version.  Please refer to [`Salesforce` documentation](https://developer.salesforce.com/docs/atlas.en-us.api_rest.meta/api_rest/intro_what_is_rest_api.htm) to make sure that `APIs` are supported in the version that is specified.
### Example
```go
config := goforce.Configuration{
	Credentials: credentials.NewPasswordCredentials(creds),
	Client:      salesforceHTTPClient,
	Version:     44,
}
```
