# Composite API
[back](../README.md)

The `composite` package is an implementation of `Salesforce APIs` centered on `Composite` operations.  These operations include:
* SObject Resources
* Query Resource
* Query All Resource
* SObject Collections Resource

As a reference, see `Salesforce API` [documentation](https://developer.salesforce.com/docs/atlas.en-us.api_rest.meta/api_rest/intro_what_is_rest_api.htm)

## Examples
The following are examples to access the `APIs`.  It is assumed that a `sfdc` [session](../session/README.md) has been created.
### Subrequest
```go
type compositeSubRequest struct {
	url         string
	body        map[string]interface{}
	method      string
	httpHeaders http.Header
	referenceID string
}

func (c *compositeSubRequest) URL() string {
	return c.url
}
func (c *compositeSubRequest) ReferenceID() string {
	return c.referenceID
}
func (c *compositeSubRequest) Method() string {
	return c.method
}
func (c *compositeSubRequest) HTTPHeaders() http.Header {
	return c.httpHeaders
}
func (c *compositeSubRequest) Body() map[string]interface{} {
	return c.body
}
```
### Composite
```go
	subRequests := []composite.Subrequester{
		&compositeSubRequest{
			url:         "/services/data/v44.0/sobjects/Account",
			method:      http.MethodPost,
			referenceID: "NewAccount",
			body: map[string]interface{}{
				"Name":          "Salesforce",
				"BillingStreet": "Landmark @ 1 Market Street",
				"BillingCity":   "San Francisco",
				"BillingState":  "California",
				"Industry":      "Technology",
			},
		},
		&compositeSubRequest{
			url:         "/services/data/v44.0/sobjects/Account/@{NewAccount.id}",
			method:      http.MethodGet,
			referenceID: "NewAccountInfo",
		},
		&compositeSubRequest{
			url:         "/services/data/v44.0/sobjects/Contact",
			method:      http.MethodPost,
			referenceID: "NewContact",
			body: map[string]interface{}{
				"lastname":      "John Doe",
				"Title":         "CTO of @{NewAccountInfo.Name}",
				"MailingStreet": "@{NewAccountInfo.BillingStreet}",
				"MailingCity":   "@{NewAccountInfo.BillingAddress.city}",
				"MailingState":  "@{NewAccountInfo.BillingState}",
				"AccountId":     "@{NewAccountInfo.Id}",
				"Email":         "jdoe@salesforce.com",
				"Phone":         "1234567890",
			},
		},
		&compositeSubRequest{
			url:         "/services/data/v44.0/sobjects/Contact/@{NewContact.id}",
			method:      http.MethodGet,
			referenceID: "NewContactInfo",
		},
	}

	resource, err := composite.NewResource(session)
	if err != nil {
		fmt.Printf("Composite Error %s\n", err.Error())
		return
	}
	value, err := resource.Retrieve(false, subRequests)
	if err != nil {
		fmt.Printf("Composite Error %s\n", err.Error())
		return
	}

	fmt.Printf("%+v\n", value)
```
