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
### Creating a Query Job
```go
	resource, err := bulk.NewResource(session)
	if err != nil {
		fmt.Printf("Bulk Resource Error %s\n", err.Error())
		return
	}

	query := "SELECT Id, Name, BillingCity FROM Account"

	jobOpts := bulk.Options{
		ColumnDelimiter: bulk.Pipe,
		Operation:       bulk.Query,
		Query:          query,
	}
	job, err := resource.CreateJob(jobOpts)
	if err != nil {
		fmt.Printf("Job Create Error %s\n", err.Error())
		return
	}
```

### Get Results from a Query Job
```go
	info, err = job.Info()
	if err != nil {
		fmt.Printf("Job Info Error %s\n", err.Error())
		return
	}

	if info.State == bulk.JobComplete {
		f, err := os.Create("data/data.csv")
		defer f.Close()
		if err != nil {
			fmt.Printf("File Create Error %s\n", err.Error())
			return
		}
		if err := job.QueryResults(f, -1, ""); err != nil {
			fmt.Printf("QueryResults Error %s\n", err.Error())
			return
		}
	}
```

### Get Results from a Query Job using Wait
```go
	if err := job.Wait(5 * time.Minute); err != nil {
		fmt.Printf("[Wait]: %s\n", err.Error())
		return
	}

	f, err := os.Create("data/data.csv")
	if err != nil {
		fmt.Printf("[File Create]: %s\n", err.Error())
		return
	}

	if err := job.QueryResults(f, -1, ""); err != nil {
		fmt.Printf("[QueryResults]: %s\n", err.Error())
		return
	}
```

### Info from many Jobs
```go
	jobsInfo, err := resource.JobsInfo(bulk.Parameters{JobType: bulk.V2Query})
	if err != nil {
		fmt.Printf("[JobsInfo]: %s\n", err.Error())
		return
	}
```

### QueryJobResults from many Query Jobs
```go
	mapErrs, err := resource.QueryJobsResults([]*Jobs{job1, job2}, []io.Write{f1, f2}, bulk.Parameters{JobType: bulk.V2Query}, 5*time.Minute, -1)
	if err != nil {
		fmt.Printf("[QueryJobsResults]: %s\n", err.Error())
		return
	}
	if len(mapErrs) > 0 {
		for k, e := range mapErrs {
			if e != nil {
				fmt.Printf("[QueryJobsResults]: JobID %s - %s\n", k, err.Error())
			}
		}
		return
	}
```