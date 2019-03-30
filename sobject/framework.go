package sobject

import (
	"time"

	"github.com/g8rswimmer/goforce"
	"github.com/g8rswimmer/goforce/session"
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

// ObjectURLs is the URL for the SObject metadata.
type ObjectURLs struct {
	CompactLayouts  string `json:"compactLayouts"`
	RowTemplate     string `json:"rowTemplate"`
	ApprovalLayouts string `json:"approvalLayouts"`
	DefaultValues   string `json:"defaultValues"`
	ListViews       string `json:"listviews"`
	Describe        string `json:"describe"`
	QuickActions    string `json:"quickActions"`
	Layouts         string `json:"layouts"`
	SObject         string `json:"sobject"`
}

type SObjectAPI struct {
	metadata *metadata
	describe *describe
}

func NewSObjectAPI(session session.Formatter) *SObjectAPI {
	return &SObjectAPI{
		metadata: &metadata{
			session: session,
		},
		describe: &describe{
			session: session,
		},
	}
}

func (api *SObjectAPI) Metadata(sobject string) (MetadataValue, error) {
	return api.metadata.Metadata(sobject)
}

func (api *SObjectAPI) Describe(sobject string) (DesribeValue, error) {
	return api.describe.Describe(sobject)
}

const objectEndpoint = "/sobjects/"

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
