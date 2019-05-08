# Composite API
[back](../../README.md)

The `batch` package is an implementation of `Salesforce APIs` centered on `Composite Batch` operations.  These operations include:
* Limits Resources
* SObject Resources
* Query All
* Query
* Search
* Connect Resources
* Chatter Resources

As a reference, see `Salesforce API` [documentation](https://developer.salesforce.com/docs/atlas.en-us.api_rest.meta/api_rest/intro_what_is_rest_api.htm)

## Examples
The following are examples to access the `APIs`.  It is assumed that a `sfdc` [session](../../session/README.md) has been created.
### Subrequest
```go
type batchSubrequester struct {
	url             string
	method          string
	richInput       map[string]interface{}
	binaryPartName  string
	binaryPartAlais string
}

func (b *batchSubrequester) URL() string {
	return b.url
}
func (b *batchSubrequester) Method() string {
	return b.method
}
func (b *batchSubrequester) BinaryPartName() string {
	return b.binaryPartName
}
func (b *batchSubrequester) BinaryPartNameAlias() string {
	return b.binaryPartAlais
}
func (b *batchSubrequester) RichInput() map[string]interface{} {
	return b.richInput
}
```
### Composite Batch
```go
	subRequests := []batch.Subrequester{
		&batchSubrequester{
			url:    "v44.0/sobjects/Account/0012E00001qLpKZQA0",
			method: http.MethodPatch,
			richInput: map[string]interface{}{
				"Name": "NewName",
			},
		},
		&batchSubrequester{
			url:    "v44.0/sobjects/Account/0012E00001qLpKZQA0",
			method: http.MethodGet,
		},
	}

	resource, err := batch.NewResource(session)
	if err != nil {
		fmt.Printf("Batch Composite Error %s", err.Error())
		fmt.Println()
		return
	}
	value, err := resource.Retrieve(false, subRequests)
	if err != nil {
		fmt.Printf("Batch Composite Error %s", err.Error())
		fmt.Println()
		return
	}

	fmt.Printf("%+v", value)
```
