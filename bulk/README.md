# Bulk 2.0 API
[back](../README.md)

The `bulk` package is an implementation of `Salesforce APIs` centered on `Bulk 2.0` operations.  These operations include:
* Creating a job
* Upload job data
* Close or Abort a job
* Delete a job
* Get all jobs
* Get job info
* Get job successful records
* Get job failed records
* Get job unprocessed records

As a reference, see `Salesforce API` [documentation](https://developer.salesforce.com/docs/atlas.en-us.api_rest.meta/api_rest/intro_what_is_rest_api.htm)

## Examples
The following are examples to access the `APIs`.  It is assumed that a `sfdc` [session](../session/README.md) has been created.
### Creating a Job
```go
	resource, err := bulk.NewResource(session)
	if err != nil {
		fmt.Printf("Bulk Resource Error %s\n", err.Error())
		return
	}

	jobOpts := bulk.Options{
		ColumnDelimiter: bulk.Pipe,
		Operation:       bulk.Insert,
		Object:          "Account",
	}
	job, err := resource.CreateJob(jobOpts)
	if err != nil {
		fmt.Printf("Job Create Error %s\n", err.Error())
		return
	}
```
### Uploading Job Data
```go
	fields := []string{
		"Name",
		"FirstName__c",
		"LastName__c",
		"Site",
	}
	formatter, err := bulk.NewFormatter(job, fields)
	if err != nil {
		fmt.Printf("Formatter Error %s\n", err.Error())
		return
	}

	failedRecord := &bulkRecord{
		fields: map[string]interface{}{
			"FirstName__c": "TallFailing",
			"LastName__c":  "PersonFailing",
			"Site":         "MySite",
		},
	}
	successRecord := &bulkRecord{
		fields: map[string]interface{}{
			"Name":         "Tall Person Success",
			"FirstName__c": "TallSuccess",
			"LastName__c":  "PersonSuccess",
			"Site":         "MySite",
		},
	}
	err = formatter.Add(failedRecord, successRecord)
	if err != nil {
		fmt.Printf("Formatter Record Error %s\n", err.Error())
		return
	}

	err = job.Upload(formatter.Reader())
	if err != nil {
		fmt.Printf("Job Upload Error %s\n", err.Error())
		return
	}
```
### Close or Abort Job
```go
	response, err := job.Close()
	if err != nil {
		fmt.Printf("Job Info Error %s\n", err.Error())
		return
	}
	fmt.Println("Bulk Job Closed")
	fmt.Println("-------------------")
	fmt.Printf("%+v\n", response)
```
### Delete a Job
```go
	err := job.Delete()
	if err != nil {
		fmt.Printf("Job Delete Error %s\n", err.Error())
		return
	}
```
### Get All Jobs
```go
	parameters := bulk.Parameters{
		IsPkChunkingEnabled: false,
		JobType:             bulk.V2Ingest,
	}
	jobs, err := resource.AllJobs(parameters)
	if err != nil {
		fmt.Printf("All Jobs Error %s\n", err.Error())
		return
	}
	fmt.Println("All Jobs")
	fmt.Println("-------------------")
	fmt.Printf("%+v\n\n", jobs)
```
### Get Job Info
```go
	info, err := job.Info()
	if err != nil {
		fmt.Printf("Job Info Error %s\n", err.Error())
		return
	}
	fmt.Println("Bulk Job Information")
	fmt.Println("-------------------")
	fmt.Printf("%+v\n", info)
```
### Get Job Successful Records
```go
	info, err = job.Info()
	if err != nil {
		fmt.Printf("Job Info Error %s\n", err.Error())
		return
	}

	if (info.NumberRecordsProcessed - info.NumberRecordsFailed) > 0 {
		successRecords, err := job.SuccessfulRecords()
		if err != nil {
			fmt.Printf("Job Success Records Error %s\n", err.Error())
			return
		}
		fmt.Println("Successful Record(s)")
		fmt.Println("-------------------")
		for _, successRecord := range successRecords {
			fmt.Printf("%+v\n\n", successRecord)
		}

	}
```
### Get Job Failed Records
```go
	info, err = job.Info()
	if err != nil {
		fmt.Printf("Job Info Error %s\n", err.Error())
		return
	}

	if info.NumberRecordsFailed > 0 {
		failedRecords, err := job.FailedRecords()
		if err != nil {
			fmt.Printf("Job Failed Records Error %s\n", err.Error())
			return
		}
		fmt.Println("Failed Record(s)")
		fmt.Println("-------------------")
		for _, failedRecord := range failedRecords {
			fmt.Printf("%+v\n\n", failedRecord)
		}
	}
```
### Get Job Unprocessed Records
```go
	info, err = job.Info()
	if err != nil {
		fmt.Printf("Job Info Error %s\n", err.Error())
		return
	}

	unprocessedRecords, err := job.UnprocessedRecords()
	if err != nil {
		fmt.Printf("Job Unprocessed Records Error %s\n", err.Error())
		return
	}
	fmt.Println("Unprocessed Record(s)")
	fmt.Println("-------------------")
	for _, unprocessedRecord := range unprocessedRecords {
		fmt.Printf("%+v\n\n", unprocessedRecord)
	}
```

