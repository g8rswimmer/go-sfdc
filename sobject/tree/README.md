# SObject Tree API
[back](../../README.md)

The `tree` package is an implementation of `Salesforce APIs` centered on `SObject Tree` operations.  These operations include:
* Create Multiple Records with Children

As a reference, see `Salesforce API` [documentation](https://developer.salesforce.com/docs/atlas.en-us.api_rest.meta/api_rest/resources_composite_sobject_tree.htm)

## Examples
The following are examples to access the `APIs`.  It is assumed that a `go-sfdc` [session](../../session/README.md) has been created.

### Builder
```go
type treeBuilder struct {
	sobject     string
	fields      map[string]interface{}
	referenceID string
}

func (b *treeBuilder) SObject() string {
	return b.sobject
}
func (b *treeBuilder) Fields() map[string]interface{} {
	return b.fields
}
func (b *treeBuilder) ReferenceID() string {
	return b.referenceID
}

	// build some records
	accountRef1Builder := &treeBuilder{
		sobject:     "Account",
		referenceID: "ref1",
		fields: map[string]interface{}{
			"name":              "SampleAccount11",
			"phone":             "1234567890",
			"website":           "www.salesforce.com",
			"numberOfEmployees": 100,
			"industry":          "Banking",
		},
	}
	accountRef2Builder := &treeBuilder{
		sobject:     "Account",
		referenceID: "ref4",
		fields: map[string]interface{}{
			"name":              "SampleAccount112",
			"Phone":             "1234567890",
			"website":           "www.salesforce2.com",
			"numberOfEmployees": 100,
			"industry":          "Banking",
		},
	}
	contactRef3Builder := &treeBuilder{
		sobject:     "Contact",
		referenceID: "ref2",
		fields: map[string]interface{}{
			"lastname": "Smith11",
			"title":    "President",
			"email":    "sample@salesforce.com",
		},
	}
	contactRef4Builder := &treeBuilder{
		sobject:     "Contact",
		referenceID: "ref3",
		fields: map[string]interface{}{
			"lastname": "Evans11",
			"title":    "Vice President",
			"email":    "sample@salesforce.com",
		},
	}

	account1RecordBuilder, err := tree.NewRecordBuilder(accountRef1Builder)
	if err != nil {
		fmt.Printf("NewRecordBuilder Error %s", err.Error())
		fmt.Println()
		return
	}
	contact1RecordBuilder, err := tree.NewRecordBuilder(contactRef3Builder)
	if err != nil {
		fmt.Printf("NewRecordBuilder Error %s", err.Error())
		fmt.Println()
		return
	}
	contact2RecordBuilder, err := tree.NewRecordBuilder(contactRef4Builder)
	if err != nil {
		fmt.Printf("NewRecordBuilder Error %s", err.Error())
		fmt.Println()
		return
	}
	account1RecordBuilder.SubRecords("contacts", contact1RecordBuilder.Build(), contact2RecordBuilder.Build())

	account2RecordBuilder, err := tree.NewRecordBuilder(accountRef2Builder)
	if err != nil {
		fmt.Printf("NewRecordBuilder Error %s", err.Error())
		fmt.Println()
		return
	}

```

### Create Accounts with Children
```go
	inserter := &treeInserter{
		sobject: "Account",
		records: []*tree.Record{
			account1RecordBuilder.Build(),
			account2RecordBuilder.Build(),
		},
	}
	resource := tree.NewResource(session)
	value, err := resource.Insert(inserter)
	if err != nil {
		fmt.Printf("resource.Insert Error %s", err.Error())
		fmt.Println()
		return
	}
	fmt.Printf("%+v\n", *value)
```
