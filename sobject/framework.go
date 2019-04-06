package sobject

import (
	"errors"
	"fmt"
	"regexp"

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
	// Metadata(string) (MetadataValue, error)
	// Describe(string) (DescribeValue, error)
	// Insert(Inserter) (InsertValue, error)
	// Update(*goforce.Record) error
	// Upsert(*goforce.Record) (UpdateValue, error)
	// Delete(*goforce.Record) (DeleteValue, error)
	// Get(*goforce.Record) error
	// AttachmentBody(string) (AttachmentBody, error)
	// DocumentBody(string) (DocumentBody, error)
	// DeletedRecords(string, time.Time, time.Time) ([]DeleteValue, error)
	// UpdatedRecords(string, time.Time, time.Time) ([]UpdateValue, error)
}

// ObjectURLs is the URL for the SObject metadata.
type ObjectURLs struct {
	CompactLayouts   string `json:"compactLayouts"`
	RowTemplate      string `json:"rowTemplate"`
	ApprovalLayouts  string `json:"approvalLayouts"`
	DefaultValues    string `json:"defaultValues"`
	ListViews        string `json:"listviews"`
	Describe         string `json:"describe"`
	QuickActions     string `json:"quickActions"`
	Layouts          string `json:"layouts"`
	SObject          string `json:"sobject"`
	UIDetailTemplate string `json:"uiDetailTemplate"`
	UIEditTemplate   string `json:"uiEditTemplate"`
	UINewRecord      string `json:"uiNewRecord"`
}

// SalesforceAPI is the structure for the Salesforce APIs for SObjects.
type SalesforceAPI struct {
	metadata *metadata
	describe *describe
	dml      *dml
}

const objectEndpoint = "/sobjects/"

// NewSalesforceAPI forms the Salesforce SObject API structure.  The
// session formatter is required to form the proper URLs and authorization
// header.
func NewSalesforceAPI(session session.Formatter) *SalesforceAPI {
	return &SalesforceAPI{
		metadata: &metadata{
			session: session,
		},
		describe: &describe{
			session: session,
		},
		dml: &dml{
			session: session,
		},
	}
}

// Metadata retrieves the SObject's metadata.
func (a *SalesforceAPI) Metadata(sobject string) (MetadataValue, error) {
	if a == nil {
		panic("salesforce api metadata has nil values")
	}

	if a.metadata == nil {
		return MetadataValue{}, errors.New("salesforce api is not initialized properly")
	}

	matching, err := regexp.MatchString(`\w`, sobject)
	if err != nil {
		return MetadataValue{}, err
	}

	if matching == false {
		return MetadataValue{}, fmt.Errorf("sobject salesforce api: %s is not a valid sobject", sobject)
	}

	return a.metadata.Metadata(sobject)
}

// Describe retrieves the SObject's describe.
func (a *SalesforceAPI) Describe(sobject string) (DescribeValue, error) {
	if a == nil {
		panic("salesforce api metadata has nil values")
	}

	if a.describe == nil {
		return DescribeValue{}, errors.New("salesforce api is not initialized properly")
	}

	matching, err := regexp.MatchString(`\w`, sobject)
	if err != nil {
		return DescribeValue{}, err
	}

	if matching == false {
		return DescribeValue{}, fmt.Errorf("sobject salesforce api: %s is not a valid sobject", sobject)
	}

	return a.describe.Describe(sobject)
}

// Insert will create a new Salesforce record.
func (a *SalesforceAPI) Insert(inserter Inserter) (InsertValue, error) {
	if a == nil {
		panic("salesforce api metadata has nil values")
	}

	if a.dml == nil {
		return InsertValue{}, errors.New("salesforce api is not initialized properly")
	}

	if inserter == nil {
		return InsertValue{}, errors.New("inserter can not be nil")
	}

	return a.dml.Insert(inserter)

}

// Update will update an existing Salesforce record.
func (a *SalesforceAPI) Update(updater Updater) error {
	if a == nil {
		panic("salesforce api metadata has nil values")
	}

	if a.dml == nil {
		return errors.New("salesforce api is not initialized properly")
	}

	if updater == nil {
		return errors.New("updater can not be nil")
	}

	return a.dml.Update(updater)

}

// Upsert will upsert an existing or new Salesforce record.
func (a *SalesforceAPI) Upsert(upserter Upserter) (UpsertValue, error) {
	if a == nil {
		panic("salesforce api metadata has nil values")
	}

	if a.dml == nil {
		return UpsertValue{}, errors.New("salesforce api is not initialized properly")
	}

	if upserter == nil {
		return UpsertValue{}, errors.New("upserter can not be nil")
	}

	return a.dml.Upsert(upserter)

}

// Delete will delete an existing Salesforce record.
func (a *SalesforceAPI) Delete(deleter Deleter) error {
	if a == nil {
		panic("salesforce api metadata has nil values")
	}

	if a.dml == nil {
		return errors.New("salesforce api is not initialized properly")
	}

	if deleter == nil {
		return errors.New("deleter can not be nil")
	}

	return a.dml.Delete(deleter)
}
