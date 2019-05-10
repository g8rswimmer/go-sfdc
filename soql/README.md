# SOQL APIs
[back](../README.md)

The `soql` package is an implementation of the `Salesforce APIs` centered on `SOQL` operations.  These operations include:
* `SOQL` query formatter
* `SOQL` query
* `SOQL` query all

 As a reference, see `Salesforce API` [documentation](https://developer.salesforce.com/docs/atlas.en-us.api_rest.meta/api_rest/intro_what_is_rest_api.htm)

## Examples
The following are examples to access the `APIs`.  It is assumed that a `go-sfdc` [session](../session/README.md) has been created.
### SOQL Builder
The following examplas cenrter around `SOQL` builder.  Although using the builder is not required to use the `API`, it is recommended as it generates the proper query statement.
#### SELECT Name, Id FROM Account WHERE Name = 'Golang'
```go
	where, err := soql.WhereEquals("Name", "Golang")
	if err != nil {
		fmt.Printf("SOQL Query Where Statement Error %s\n", err.Error())
		return
	}
	input := soql.QueryInput{
		ObjectType: "Account",
		FieldList: []string{
			"Name",
			"Id",
		},
		Where: where,
	}
	queryStmt, err := soql.NewQuery(input)
	if err != nil {
		fmt.Printf("SOQL Query Statement Error %s\n", err.Error())
		return
	}
	stmt, err := queryStmt.Format()
	if err != nil {
		fmt.Printf("SOQL Query Statement Error %s\n", err.Error())
		return
	}

	fmt.Println("SOQL Query Statement")
	fmt.Println("-------------------")
	fmt.Println(stmt)
```
#### SELECT Name, Id, (SELECT LastName FROM Contacts) FROM Account 
```go
	subInput := soql.QueryInput{
		ObjectType: "Contacts",
		FieldList: []string{
			"LastName",
		},
	}
	subQuery, err := soql.NewQuery(subInput)
	if err != nil {
		fmt.Printf("SOQL Sub Query Error %s\n", err.Error())
		return
	}

	input := soql.QueryInput{
		ObjectType: "Account",
		FieldList: []string{
			"Name",
			"Id",
		},
		SubQuery: []soql.Querier{
			subQuery,
		},
	}
	queryStmt, err := soql.NewQuery(input)
	if err != nil {
		fmt.Printf("SOQL Query Statement Error %s\n", err.Error())
		return
	}

	stmt, err := queryStmt.Format()
	if err != nil {
		fmt.Printf("SOQL Query Statement Error %s\n", err.Error())
		return
	}

	fmt.Println("SOQL Query Statement")
	fmt.Println("-------------------")
	fmt.Println(stmt)
```
### SOQL Query
The following example demostrates how to `SOQL` query.  It is assumed that a session has need created and a `SOQL` statement has been built.
The `SOQL` statement is as follows:
```
SELECT
  Name,
  Id,
  (
      SELECT
        LastName
      FROM
        Contacts  
  )
FROM 
  Account
```  
```go
	resource := soql.NewResource(session)
	result, err := resource.Query(queryStmt, false)
	if err != nil {
		fmt.Printf("SOQL Query Error %s\n", err.Error())
		return
	}
	fmt.Println("SOQL Query")
	fmt.Println("-------------------")
	fmt.Printf("Done: %t\n", result.Done())
	fmt.Printf("Total Size: %d\n", result.TotalSize())
	fmt.Printf("Next Records URL: %t\n\n", result.MoreRecords())

	for _, rec := range result.Records() {
		r := rec.Record()
		fmt.Printf("SObject: %s\n", r.SObject())
		fmt.Printf("Fields: %v\n", r.Fields())
		for obj, subResult := range rec.Subresults() {
			fmt.Printf("Sub Result: %s\n", obj)
			fmt.Printf("Done: %t\n", subResult.Done())
			fmt.Printf("Total Size: %d\n", subResult.TotalSize())
			fmt.Printf("Next Records URL: %t\n", subResult.MoreRecords())
			fmt.Println()
			for _, subRec := range subResult.Records() {
				sr := subRec.Record()
				fmt.Printf("SObject: %s\n", sr.SObject())
				fmt.Printf("Fields: %v\n", sr.Fields())
			}
		}
	}
```