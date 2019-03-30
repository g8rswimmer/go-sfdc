package sobject

import (
	"time"

	"github.com/g8rswimmer/goforce"
)

// Framework is the Salesforce SObject API framework.
//
// Metadata will return a SObject's metadata.
//
// Describe will return a SObject's descibe
//
// Insert will insert a record into the Salesforce org.
//
// Update will update a record into the Salesforce org.
//
// Upsert will upsert (update or insert) a record inot the Salesforce org.
// The external id, defined in the record, can be used as an ID.
//
// Delete will remove the record into the Salesforce org.
//
// Get will retrieve an objects record from the id or external id.
//
// AttachmentBody will retrieve the attachment's blob.
//
// DocumentBody will retrieve the document's blob.
//
// DeletedRecords will return the delete records from a date range.
//
// UpdatedRecords will return the updated records from a date range.
type Framework interface {
	Metadata(string) (MetadataValue, error)
	Describe(string) (DesribeValue, error)
	Insert(*goforce.Record) (InsertValue, error)
	Update(*goforce.Record) (UpdateValue, error)
	Upsert(*goforce.Record) (UpdateValue, error)
	Delete(*goforce.Record) (DeleteValue, error)
	Get(*goforce.Record) error
	AttachmentBody(string) (AttachmentBody, error)
	DocumentBody(string) (DocumentBody, error)
	DeletedRecords(string, time.Time, time.Time) ([]DeleteValue, error)
	UpdatedRecords(string, time.Time, time.Time) ([]UpdateValue, error)
}

type MetadataValue struct {
}

type DesribeValue struct {
}

type InsertValue struct {
}

type UpdateValue struct {
}

type DeleteValue struct {
}

type AttachmentBody struct {
}
type DocumentBody struct {
}
