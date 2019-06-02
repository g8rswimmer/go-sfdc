package bulk

import (
	"errors"
	"fmt"
	"strings"
)

type Record interface {
	Fields() map[string]interface{}
}

type Formatter struct {
	job        *Job
	fields     []string
	insertNull bool
	sb         strings.Builder
}

func NewFormatter(job *Job, fields []string, insertNull bool) (*Formatter, error) {
	if job == nil {
		return nil, errors.New("bulk formatter: job is required for the formatter")
	}
	if len(fields) == 0 {
		return nil, errors.New("bulk formatter: fields are required")
	}
	return &Formatter{
		job:        job,
		fields:     fields,
		insertNull: insertNull,
		sb:         strings.Builder{},
	}, nil
}

func (f *Formatter) Add(records ...Record) error {
	if records == nil {
		return errors.New("bulk formatter: record interface can not be nil")
	}

	for _, record := range records {
		recFields := record.Fields()
		values := make([]string, len(f.fields))
		for idx, field := range f.fields {
			if f.insertNull {
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
		f.sb.WriteString(strings.Join(values, f.job.delimiter()))
		f.sb.WriteString(f.job.newline())

	}

	return nil
}

func (f *Formatter) Reader() *strings.Reader {
	return strings.NewReader(f.sb.String())
}
