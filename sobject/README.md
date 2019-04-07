# SObject APIs
The `sobject` package is an implementation of `Salesforce APIs` centered on `SObject` operations.  This operations include:
* Metadata
* Describe
* DML
  - Insert
  - Update
  - Upsert
  - Delete
* Query
  - With `Salesforce` ID
  - With external ID
* List of Deleted records
* List of Updated records
* Get `Attachment` body
* Get `Document` body

As a reference, see `Salesforce API` [documentation](https://developer.salesforce.com/docs/atlas.en-us.api_rest.meta/api_rest/intro_what_is_rest_api.htm)

## Examples
The following are examples to access the `APIs`.  It is assumed that a `goforce` session has been created.
### Metadata
```
sobjectAPI := sobject.NewSalesforceAPI(session)

metadata, err := sobjectAPI.Metadata("Account")

if err != nil {
  fmt.Printf("Error %s", err.Error())
  fmt.Println()
  return
}

fmt.Println("Account Metadata")
fmt.Println("-------------------")
fmt.Printf("%+v\n", metadata)
fmt.Println()
```
### Describe
```
sobjectAPI := sobject.NewSalesforceAPI(session)

describe, err := sobjectAPI.Describe("Account")

if err != nil {
  fmt.Printf("Error %s", err.Error())
  fmt.Println()
  return
}

fmt.Println("Account Describe")
fmt.Println("-------------------")
fmt.Printf("%+v\n", describe)
fmt.Println()
```
### DML Insert
```
type dml struct {
	sobject       string
	fields        map[string]interface{}
}

func (d *dml) SObject() string {
	return d.sobject
}
func (d *dml) Fields() map[string]interface{} {
	return d.fields
}


sobjectAPI := sobject.NewSalesforceAPI(session)

dml := &dml{
  sobject: "Account",
  fields: map[string]interface{}{
    "Name":            "Test Account",
    "MyUID__c":        "AccExID222",
    "MyCustomText__c": "My fun text",
    "Phone":           "9045551212",
  },
}

insertValue, err := sobjectAPI.Insert(dml)

if err != nil {
  fmt.Printf("Insert Error %s", err.Error())
  fmt.Println()
  return
}

fmt.Println("Account")
fmt.Println("-------------------")
fmt.Printf("%+v\n", insertValue)
fmt.Println()
```
### DML Update
```
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


sobjectAPI := sobject.NewSalesforceAPI(session)

dml := &dml{
  sobject: "Account",
}

dml.id = "Account Salesforce ID"
dml.fields["Phone"] = "6065551212"
dml.fields["MyCustomText__c"] = "updated text"

err = sobjectAPI.Update(dml)

if err != nil {
  fmt.Printf("Update Error %s", err.Error())
  fmt.Println()
  return
}

fmt.Println("Account Updated")
fmt.Println("-------------------")
fmt.Println()

```
### DML Upsert
```
type dml struct {
	sobject       string
	fields        map[string]interface{}
	id            string
	externalField string
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
func (d *dml) ExternalField() string {
	return d.externalField
}

sobjectAPI := sobject.NewSalesforceAPI(session)

dml := &dml{
  sobject: "Account",
}
dml.id = "AccExID345"
dml.externalField = "MyUID__c"
dml.fields["Name"] = "Upsert Update"

upsertValue, err := sobjectAPI.Upsert(dml)

if err != nil {
  fmt.Printf("Upsert Error %s", err.Error())
  fmt.Println()
  return
}

fmt.Println("Account Upsert")
fmt.Println("-------------------")
fmt.Printf("%+v\n", upsertValue)
fmt.Println()
```
### DML Delete
```
type dml struct {
	sobject       string
	id            string
}

func (d *dml) SObject() string {
	return d.sobject
}
func (d *dml) ID() string {
	return d.id
}

sobjectAPI := sobject.NewSalesforceAPI(session)

dml := &dml{
  sobject: "Account",
  id:      "0012E00001oHQDNQA4",
}

err = sobjectAPI.Delete(dml)

if err != nil {
  fmt.Printf("Upsert Error %s", err.Error())
  fmt.Println()
  return
}

fmt.Println("Account Deleted")

```
