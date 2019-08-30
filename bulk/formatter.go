package bulk

import (
	"errors"
	"fmt"
	"strings"
)

// Record is the interface to the fields of the bulk uploader record.
type Record interface {
	Fields() map[string]interface{}
	InsertNull() bool
}

// Formatter is the object that will add records for the bulk uploader.
type Formatter struct {
	job    *Job
	fields []string
	sb     strings.Builder
}

// NewFormatter creates a new formatter using the job and the list of fields.
func NewFormatter(job *Job, fields []string) (*Formatter, error) {
	if job == nil {
		return nil, errors.New("bulk formatter: job is required for the formatter")
	}
	if len(fields) == 0 {
		return nil, errors.New("bulk formatter: fields are required")
	}

	f := &Formatter{
		job:    job,
		fields: fields,
		sb:     strings.Builder{},
	}

	if _, err := f.sb.WriteString(strings.Join(fields, job.delimiter())); err != nil {
		return nil, err
	}
	if _, err := f.sb.WriteString(job.newline()); err != nil {
		return nil, err
	}

	return f, nil
}

// Add will place a record in the bulk uploader.
func (f *Formatter) Add(records ...Record) error {
	if records == nil {
		return errors.New("bulk formatter: record interface can not be nil")
	}

	for _, record := range records {
		_, err := f.sb.WriteString(f.buildRecordString(record))
		if err != nil {
			return err
		}
		_, err = f.sb.WriteString(f.job.newline())
		if err != nil {
			return err
		}

	}

	return nil
}

func (f *Formatter) buildRecordString(record Record) string {
	recFields := record.Fields()
	values := make([]string, len(f.fields))
	insertNull := record.InsertNull()
	for idx, field := range f.fields {
		if insertNull {
			values[idx] = "#N/A"
		} else {
			values[idx] = ""
		}
		if value, ok := recFields[field]; ok {
			if value != nil {
				values[idx] = fmt.Sprintf("%v", value)
			}
		}
	}
	return strings.Join(values, f.job.delimiter())
}

// Reader will return a reader of the bulk uploader field record body.
func (f *Formatter) Reader() *strings.Reader {
	return strings.NewReader(f.sb.String())
}
