# SObject Collection APIs
[back](../../README.md)

The `collection` package is an implementation of `Salesforce APIs` centered on `SObject Collection` operations.  These operations include:
* Create Multiple Records
* Update Multiple Records
* Delete Multiple Records
* Retrieve Multiple Records

As a reference, see `Salesforce API` [documentation](https://developer.salesforce.com/docs/atlas.en-us.api_rest.meta/api_rest/resources_composite_sobjects_collections.htm)

## Examples
The following are examples to access the `APIs`.  It is assumed that a `go-sfdc` [session](../../session/README.md) has been created.

```go
type dml struct {
	sobject       string
	fields        map[string]interface{}
	id            string
}

func (d *dml) SObject() string {
	return d.sobject
}
func (d *dml) Fields() map[string]interface{} {
	return d.fields
}
func (d *dml) ID() string {
	return d.id
}

type query struct {
	sobject  string
	id       string
	fields   []string
}

func (q *query) SObject() string {
	return q.sobject
}
func (q *query) ID() string {
	return q.id
}
func (q *query) Fields() []string {
	return q.fields
}
```
### Create Multiple Records
```go
	// insert some records
	var insertRecords []sobject.Inserter

	acc1 := &dml{
		sobject: "Account",
		fields: map[string]interface{}{
			"Name":  "Collections Demo",
			"Phone": "9045551212",
		},
	}
	insertRecords = append(insertRecords, acc1)
	acc2 := &dml{
		sobject: "Account",
		fields: map[string]interface{}{
			"Name":      "Collections Demo Two",
			"Active__c": "Yes",
		},
	}
	insertRecords = append(insertRecords, acc2)
	con1 := &dml{
		sobject: "Contact",
		fields: map[string]interface{}{
			"LastName":       "Last Name Demo",
			"AssistantPhone": "6065551212",
		},
	}
	insertRecords = append(insertRecords, con1)
	con2 := &dml{
		sobject: "Contact",
		fields: map[string]interface{}{
			"LastName": "Last Name Demo Two",
			"Level__c": "Tertiary",
		},
	}
	insertRecords = append(insertRecords, con2)

	resource := collections.NewResources(session)
	values, err := resource.Insert(true, insertRecords)
	if err != nil {
		fmt.Printf("Collection Error %s", err.Error())
		fmt.Println()
		return
	}

	fmt.Println("Collections Inserted")
	fmt.Println("-------------------")
	for _, value := range values {
		fmt.Printf("%+v\n", value)
	}
	fmt.Println()
```
### Update Multiple Records
```go
	var updateRecords []sobject.Updater

	acc1 := &dml{
		sobject: "Account",
		fields: map[string]interface{}{
			"Name":      "Collections Demo Update",
			"Phone":     "4805551212",
			"Active__c": "No",
		},
		id: "0012E00001pXGh2QAG",
	}
	updateRecords = append(updateRecords, acc1)
	acc2 := &dml{
		sobject: "Account",
		fields: map[string]interface{}{
			"Name":      "Collections Demo Two Update",
			"Phone":     "6065551212",
			"Active__c": "Yes",
		},
		id: "0012ExxxAG",
	}
	updateRecords = append(updateRecords, acc2)
	con1 := &dml{
		sobject: "Contact",
		fields: map[string]interface{}{
			"LastName":       "Last Name Demo Update",
			"AssistantPhone": "6065551212",
			"Level__c":       "Tertiary",
		},
		id: "0032ExxxjQAJ",
	}
	updateRecords = append(updateRecords, con1)
	con2 := &dml{
		sobject: "Contact",
		fields: map[string]interface{}{
			"LastName":       "Last Name Demo Two Update",
			"AssistantPhone": "6025551212",
			"Level__c":       "Tertiary",
		},
		id: "0032ExxxkQAJ",
	}
	updateRecords = append(updateRecords, con2)

	resource := collections.NewResources(session)
	values, err := resource.Update(true, updateRecords)
	if err != nil {
		fmt.Printf("Collection Error %s", err.Error())
		fmt.Println()
		return
	}

	fmt.Println("Collections Updated")
	fmt.Println("-------------------")
	for _, value := range values {
		fmt.Printf("%+v\n", value)
	}
	fmt.Println()
```
### Delete Multiple Records
```go
	deleteRecords := []string{
		"0012Exxx2QAG",
		"0012Exxx3QAG",
		"0032ExxxEjQAJ",
		"0032ExxxEkQAJ",
	}

	resource := collections.NewResources(session)
	values, err := resource.Delete(true, deleteRecords)
	if err != nil {
		fmt.Printf("Collection Error %s", err.Error())
		fmt.Println()
		return
	}

	fmt.Println("Collections Deleted")
	fmt.Println("-------------------")
	for _, value := range values {
		fmt.Printf("%+v\n", value)
	}
	fmt.Println()

```
### Retrieve Multiple Records
```go
	var queryRecords []sobject.Querier

	acc1 := &query{
		sobject: "Account",
		fields: []string{
			"Name",
			"Phone",
		},
		id: "0012Exxx2QAG",
	}
	queryRecords = append(queryRecords, acc1)
	acc2 := &query{
		sobject: "Account",
		fields: []string{
			"Name",
			"Active__c",
		},
		id: "0012Exxx3QAG",
	}
	queryRecords = append(queryRecords, acc2)

	resource := collections.NewResources(session)
	values, err := resource.Query("Account", queryRecords)
	if err != nil {
		fmt.Printf("Collection Error %s", err.Error())
		fmt.Println()
		return
	}

	fmt.Println("Collections Inserted")
	fmt.Println("-------------------")
	for _, value := range values {
		fmt.Printf("%+v\n", *value)
	}
	fmt.Println()

```
