# goforce
This is a `golang` library for interfacing with `Salesforce` APIs.

## Usage
To use this library, the following will need to be done.
* Create `Salesforce` [credentials](./credentials/README.md) to properly authenticate with the `Salesforce org`
* Configure
* Open a [session](./session/README.md)
* Use the `APIs`
  - [SObject APIs](./sobject/README.md)

## Configuration
The configuration defines several parameters that can be used by the library.  The configuration is a defined per [session](./session/README.md).
* `Credentails` - this is an implementation of the `credentials.Provider` interface
* `Client` - the the HTTP client used by the `APIs`
* `Version` - is the `Salesforce` version.  Please refer to `Salesforce` documentation to make sure that `APIs` are supported in the version that is specified.
### Example
```
config := goforce.Configuration{
	Credentials: credentials.NewPasswordCredentials(creds),
	Client:      salesforceHTTPClient,
	Version:     44,
}
```
