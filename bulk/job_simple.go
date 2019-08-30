package bulk

import (
	"io/ioutil"

	b64 "encoding/base64"

	"github.com/g8rswimmer/go-sfdc/session"
)

const (
	// A request can provide CSV data that does not in total exceed 150 MB of base64 encoded content.
	// https://developer.salesforce.com/docs/atlas.en-us.api_bulk_v2.meta/api_bulk_v2/upload_job_data.htm
	maxUploadSizeBytes = 150000000 // I think Salesforce is using this number and not the actual 150MB byte value
)

// ProcessDataAsBulkJobs recieves all the necessary information to upload large amounts of data
// breaking it into individual jobs as needed if they cross the size threshold, returns the created jobs
func ProcessDataAsBulkJobs(session session.ServiceFormatter, options Options, fields []string, data []Record) ([]*Job, error) {
	// Get an initial job
	jobs := []*Job{}
	j, err := openBulkJob(session, options)
	if err != nil {
		return jobs, err
	}
	jobs = append(jobs, j)

	f, err := NewFormatter(j, fields)
	if err != nil {
		return jobs, err
	}
	// Get the by counts of the headers and the line encoding
	headers, err := ioutil.ReadAll(f.Reader())
	if err != nil {
		return jobs, err
	}
	sizeOfHeaders := len(b64.StdEncoding.EncodeToString(headers))
	sizeOfLineEnding := len(b64.StdEncoding.EncodeToString([]byte(options.LineEnding)))

	// loop through output from records, analyze the impact of adding another record to the total CSV size
	rollingSize := sizeOfHeaders + sizeOfLineEnding
	rollingRowCount := 0
	for _, v := range data {

		// get size of encoded data for this row
		sizeOfRow := len(b64.StdEncoding.EncodeToString([]byte(f.buildRecordString(v))))
		rollingSize += sizeOfRow + sizeOfLineEnding
		// if the new record would push it over the threshold send the CSV to the open job and start a new one
		if rollingSize > maxUploadSizeBytes {
			// Send the data, close the job, start a new job, make a new formatter, do it again
			err := sendDataToJob(j, f)
			if err != nil {
				return jobs, err
			}

			// Get a new job to load against
			j, err = openBulkJob(session, options)
			if err != nil {
				return jobs, err
			}
			jobs = append(jobs, j)

			// Build a new formatter and reset counters
			f, err = NewFormatter(j, fields)
			if err != nil {
				return jobs, err
			}
			rollingSize = sizeOfHeaders + sizeOfLineEnding + sizeOfRow + sizeOfLineEnding
			rollingRowCount = 0
		}
		f.Add(v)
		rollingRowCount++
	}

	// send the last batch
	if rollingRowCount > 0 {
		err = sendDataToJob(j, f)
		if err != nil {
			return jobs, err
		}
	}
	return jobs, nil
}

func sendDataToJob(j *Job, f *Formatter) error {
	err := j.Upload(f.Reader())
	if err != nil {
		return err
	}

	err = closeBulkJob(j)
	if err != nil {
		return err
	}
	return nil
}

func openBulkJob(session session.ServiceFormatter, jobOpts Options) (*Job, error) {
	resource, err := NewResource(session)
	if err != nil {
		return &Job{}, err
	}

	job, err := resource.CreateJob(jobOpts)
	if err != nil {
		return &Job{}, err
	}
	return job, nil
}

func closeBulkJob(job *Job) error {
	_, err := job.Close()
	if err != nil {
		return err
	}
	return nil
}
