package tree

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/g8rswimmer/goforce/sobject"
)

// Builder is the SObject Tree builder for the
// composite SObject Tree API.
type Builder interface {
	sobject.Inserter
	ReferenceID() string
}

// RecordBuilder is record builder for the
// composite SObject Tree API.
type RecordBuilder struct {
	record Record
}

// NewRecordBuilder will create a new builder.  If the SObject is
// not value or the reference ID is empyt, an error will be returned.
func NewRecordBuilder(builder Builder) (*RecordBuilder, error) {
	sobject := builder.SObject()
	matching, err := regexp.MatchString(`\w`, sobject)
	if err != nil {
		return nil, err
	}
	if matching == false {
		return nil, fmt.Errorf("tree builder: %s is not a valid sobject", sobject)
	}
	if builder.ReferenceID() == "" {
		return nil, errors.New("tree builder: reference id must be present")
	}
	return &RecordBuilder{
		record: Record{
			Attributes: Attributes{
				Type:        sobject,
				ReferenceID: builder.ReferenceID(),
			},
			Fields:  builder.Fields(),
			Records: make(map[string][]*Record),
		},
	}, nil
}

// SubRecords will add subrecords to the object.
func (rb *RecordBuilder) SubRecords(sobjects string, records ...*Record) {
	if rb == nil {
		panic("record builder can not be nil")
	}
	var subRecords []*Record
	if subRec, ok := rb.record.Records[sobjects]; ok {
		subRecords = subRec
	}
	subRecords = append(subRecords, records...)
	rb.record.Records[sobjects] = subRecords
}

// Build will create the composite tree record.
func (rb *RecordBuilder) Build() *Record {
	if rb == nil {
		panic("record builder can not be nil")
	}
	return &rb.record
}
