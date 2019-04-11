# SObject APIs
[back](../README.md)

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
The following are examples to access the `APIs`.  It is assumed that a `goforce` [session](../session/README.md) has been created.
### Metadata
```go
sobjResources := sobject.NewResources(session)

metadata, err := sobjResources.Metadata("Account")

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
```go
sobjResources := sobject.NewResources(session)

describe, err := sobjResources.Describe("Account")

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
```go
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


sobjResources := sobject.NewResources(session)

dml := &dml{
  sobject: "Account",
  fields: map[string]interface{}{
    "Name":            "Test Account",
    "MyUID__c":        "AccExID222",
    "MyCustomText__c": "My fun text",
    "Phone":           "9045551212",
  },
}

insertValue, err := sobjResources.Insert(dml)

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


sobjResources := sobject.NewResources(session)

dml := &dml{
  sobject: "Account",
}

dml.id = "Account Salesforce ID"
dml.fields["Phone"] = "6065551212"
dml.fields["MyCustomText__c"] = "updated text"

err = sobjResources.Update(dml)

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
```go
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

sobjResources := sobject.NewResources(session)

dml := &dml{
  sobject: "Account",
}
dml.id = "AccExID345"
dml.externalField = "MyUID__c"
dml.fields["Name"] = "Upsert Update"

upsertValue, err := sobjResources.Upsert(dml)

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
```go
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

sobjResources := sobject.NewResources(session)

dml := &dml{
  sobject: "Account",
  id:      "0012E00001oHQDNQA4",
}

err = sobjResources.Delete(dml)

if err != nil {
  fmt.Printf("Upsert Error %s", err.Error())
  fmt.Println()
  return
}

fmt.Println("Account Deleted")

```
### Query: With Salesforce ID
Return all `SObject` fields.
```go
type query struct {
	sobject string
	id      string
	fields  []string
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

sobjResources := sobject.NewResources(session)

query := &query{
  sobject: "Account",
  id:      "Account Salesforce ID",
}

record, err := sobjResources.Query(query)
if err != nil {
  fmt.Printf("Query Error %s", err.Error())
  fmt.Println()
  return
}

fmt.Println("Account Query")
fmt.Println("-------------------")
fmt.Printf("%+v", record)
fmt.Println()

```
Return specific `SObject` fields.
```go
type query struct {
	sobject string
	id      string
	fields  []string
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

sobjResources := sobject.NewResources(session)

query := &query{
  sobject: "Account",
  id:      "Account Salesforce ID",
  fields: []string{
    "Name",
    "MyUID__c",
    "Phone",
    "MyCustomText__c",
  },
}

record, err := sobjResources.Query(query)
if err != nil {
  fmt.Printf("Query Error %s", err.Error())
  fmt.Println()
  return
}

fmt.Println("Account Query")
fmt.Println("-------------------")
fmt.Printf("%+v", record)
fmt.Println()

```
### Query: With External ID
Return all `SObject` fields.
```go
type query struct {
	sobject  string
	id       string
	external string
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
func (q *query) ExternalField() string {
	return q.external
}

sobjResources := sobject.NewResources(session)

query := &query{
  sobject:  "Account",
  id:       "AccExID234",
  external: "MyUID__c",
}

record, err := sobjResources.ExternalQuery(query)
if err != nil {
  fmt.Printf("Query Error %s", err.Error())
  fmt.Println()
  return
}

fmt.Println("Account Query")
fmt.Println("-------------------")
fmt.Printf("%+v", record)
fmt.Println()
```
Return specific `SObject` fields.
```go
type query struct {
	sobject  string
	id       string
	external string
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
func (q *query) ExternalField() string {
	return q.external
}

sobjResources := sobject.NewResources(session)

query := &query{
  sobject:  "Account",
  id:       "AccExID234",
  external: "MyUID__c",
  fields: []string{
    "Name",
    "Phone",
    "MyCustomText__c",
  },
}

record, err := sobjResources.ExternalQuery(query)
if err != nil {
  fmt.Printf("Query Error %s", err.Error())
  fmt.Println()
  return
}

fmt.Println("Account Query")
fmt.Println("-------------------")
fmt.Printf("%+v", record)
fmt.Println()
```
### List of Deleted Records
```go
sobjResources := sobject.NewResources(session)

deletedRecords, err := sobjResources.DeletedRecords("Account", time.Now().Add(time.Hour*-12), time.Now())
if err != nil {
  fmt.Printf("Deleted Records Error %s", err.Error())
  fmt.Println()
  return
}

fmt.Println("Deleted Account Records")
fmt.Println("-------------------")
fmt.Printf("%+v", deletedRecords)
fmt.Println()
```
### List of Updated Records
```go
sobjResources := sobject.NewResources(session)

updatedRecords, err := sobjResources.UpdatedRecords("Account", time.Now().Add(time.Hour*-12), time.Now())
if err != nil {
  fmt.Printf("Deleted Records Error %s", err.Error())
  fmt.Println()
  return
}

fmt.Println("Updated Account Records")
fmt.Println("-------------------")
fmt.Printf("%+v", updatedRecords)
fmt.Println()

```
### Get Attachment and Document Content
```go
sobjResources := sobject.NewResources(session)

attachment, err := sobjResources.GetContent("Attachment ID", sobject.AttachmentType)
if err != nil {
  fmt.Printf("Error %s", err.Error())
  fmt.Println()
  return
}

document, err := sobjResources.GetContent("Document ID", sobject.DocumentType)
if err != nil {
  fmt.Printf("Error %s", err.Error())
  fmt.Println()
  return
}
```
