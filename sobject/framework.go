package sobject

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/g8rswimmer/goforce"
	"github.com/g8rswimmer/goforce/session"
)

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
	query    *query
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
		query: &query{
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

// Query returns a SObject record using the Salesforce ID.
func (a *SalesforceAPI) Query(querier Querier) (*goforce.Record, error) {
	if a == nil {
		panic("salesforce api metadata has nil values")
	}

	if a.query == nil {
		return nil, errors.New("salesforce api is not initialized properly")
	}

	if querier == nil {
		return nil, errors.New("querier can not be nil")
	}

	return a.query.Query(querier)
}

// ExternalQuery returns a SObject record using an external ID field.
func (a *SalesforceAPI) ExternalQuery(querier ExternalQuerier) (*goforce.Record, error) {
	if a == nil {
		panic("salesforce api metadata has nil values")
	}

	if a.query == nil {
		return nil, errors.New("salesforce api is not initialized properly")
	}

	if querier == nil {
		return nil, errors.New("querier can not be nil")
	}

	return a.query.ExternalQuery(querier)
}

// DeletedRecords returns a list of records that have been deleted from a date range.
func (a *SalesforceAPI) DeletedRecords(sobject string, startDate, endDate time.Time) (DeletedRecords, error) {
	if a == nil {
		panic("salesforce api metadata has nil values")
	}

	if a.query == nil {
		return DeletedRecords{}, errors.New("salesforce api is not initialized properly")
	}

	matching, err := regexp.MatchString(`\w`, sobject)
	if err != nil {
		return DeletedRecords{}, err
	}

	if matching == false {
		return DeletedRecords{}, fmt.Errorf("sobject salesforce api: %s is not a valid sobject", sobject)
	}

	return a.query.DeletedRecords(sobject, startDate, endDate)
}

// UpdatedRecords returns a list of records that have been updated from a date range.
func (a *SalesforceAPI) UpdatedRecords(sobject string, startDate, endDate time.Time) (UpdatedRecords, error) {
	if a == nil {
		panic("salesforce api metadata has nil values")
	}

	if a.query == nil {
		return UpdatedRecords{}, errors.New("salesforce api is not initialized properly")
	}

	matching, err := regexp.MatchString(`\w`, sobject)
	if err != nil {
		return UpdatedRecords{}, err
	}

	if matching == false {
		return UpdatedRecords{}, fmt.Errorf("sobject salesforce api: %s is not a valid sobject", sobject)
	}

	return a.query.UpdatedRecords(sobject, startDate, endDate)
}

// GetContent returns the blob from a content SObject.
func (a *SalesforceAPI) GetContent(id string, content ContentType) ([]byte, error) {
	if a == nil {
		panic("salesforce api metadata has nil values")
	}

	if a.query == nil {
		return nil, errors.New("salesforce api is not initialized properly")
	}

	if id == "" {
		return nil, fmt.Errorf("sobject salesforce api: %s can not be empty", id)
	}

	switch content {
	case AttachmentType:
	case DocumentType:
	default:
		return nil, fmt.Errorf("sobject salesforce: content type (%s) is not supported", string(content))
	}

	return a.query.GetContent(id, content)
}
