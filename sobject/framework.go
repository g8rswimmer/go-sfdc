package sobject

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"go-sfdc"
	"go-sfdc/session"
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

// Resources is the structure for the Salesforce APIs for SObjects.
type Resources struct {
	metadata *metadata
	describe *describe
	dml      *dml
	query    *query
}

const objectEndpoint = "/sobjects/"

// NewResources forms the Salesforce SObject resource structure.  The
// session formatter is required to form the proper URLs and authorization
// header.
func NewResources(session session.ServiceFormatter) (*Resources, error) {
	if session == nil {
		return nil, errors.New("sobject resource: session can not be nil")
	}
	return &Resources{
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
	}, nil
}

// Metadata retrieves the SObject's metadata.
func (r *Resources) Metadata(sobject string) (MetadataValue, error) {
	if r.metadata == nil {
		return MetadataValue{}, errors.New("salesforce api is not initialized properly")
	}

	matching, err := regexp.MatchString(`\w`, sobject)
	if err != nil {
		return MetadataValue{}, err
	}

	if matching == false {
		return MetadataValue{}, fmt.Errorf("sobject salesforce api: %s is not a valid sobject", sobject)
	}

	return r.metadata.callout(sobject)
}

// Describe retrieves the SObject's describe.
func (r *Resources) Describe(sobject string) (DescribeValue, error) {
	if r.describe == nil {
		return DescribeValue{}, errors.New("salesforce api is not initialized properly")
	}

	matching, err := regexp.MatchString(`\w`, sobject)
	if err != nil {
		return DescribeValue{}, err
	}

	if matching == false {
		return DescribeValue{}, fmt.Errorf("sobject salesforce api: %s is not a valid sobject", sobject)
	}

	return r.describe.callout(sobject)
}

// Insert will create a new Salesforce record.
func (r *Resources) Insert(inserter Inserter) (InsertValue, error) {
	if r.dml == nil {
		return InsertValue{}, errors.New("salesforce api is not initialized properly")
	}

	if inserter == nil {
		return InsertValue{}, errors.New("inserter can not be nil")
	}

	return r.dml.insertCallout(inserter)

}

// Update will update an existing Salesforce record.
func (r *Resources) Update(updater Updater) error {
	if r.dml == nil {
		return errors.New("salesforce api is not initialized properly")
	}

	if updater == nil {
		return errors.New("updater can not be nil")
	}

	return r.dml.updateCallout(updater)

}

// Upsert will upsert an existing or new Salesforce record.
func (r *Resources) Upsert(upserter Upserter) (UpsertValue, error) {
	if r.dml == nil {
		return UpsertValue{}, errors.New("salesforce api is not initialized properly")
	}

	if upserter == nil {
		return UpsertValue{}, errors.New("upserter can not be nil")
	}

	return r.dml.upsertCallout(upserter)

}

// Delete will delete an existing Salesforce record.
func (r *Resources) Delete(deleter Deleter) error {
	if r.dml == nil {
		return errors.New("salesforce api is not initialized properly")
	}

	if deleter == nil {
		return errors.New("deleter can not be nil")
	}

	return r.dml.deleteCallout(deleter)
}

// Query returns a SObject record using the Salesforce ID.
func (r *Resources) Query(querier Querier) (*sfdc.Record, error) {
	if r.query == nil {
		return nil, errors.New("salesforce api is not initialized properly")
	}

	if querier == nil {
		return nil, errors.New("querier can not be nil")
	}

	return r.query.callout(querier)
}

// ExternalQuery returns a SObject record using an external ID field.
func (r *Resources) ExternalQuery(querier ExternalQuerier) (*sfdc.Record, error) {
	if r.query == nil {
		return nil, errors.New("salesforce api is not initialized properly")
	}

	if querier == nil {
		return nil, errors.New("querier can not be nil")
	}

	return r.query.externalCallout(querier)
}

// DeletedRecords returns a list of records that have been deleted from a date range.
func (r *Resources) DeletedRecords(sobject string, startDate, endDate time.Time) (DeletedRecords, error) {
	if r.query == nil {
		return DeletedRecords{}, errors.New("salesforce api is not initialized properly")
	}

	matching, err := regexp.MatchString(`\w`, sobject)
	if err != nil {
		return DeletedRecords{}, err
	}

	if matching == false {
		return DeletedRecords{}, fmt.Errorf("sobject salesforce api: %s is not a valid sobject", sobject)
	}

	return r.query.deletedRecordsCallout(sobject, startDate, endDate)
}

// UpdatedRecords returns a list of records that have been updated from a date range.
func (r *Resources) UpdatedRecords(sobject string, startDate, endDate time.Time) (UpdatedRecords, error) {
	if r.query == nil {
		return UpdatedRecords{}, errors.New("salesforce api is not initialized properly")
	}

	matching, err := regexp.MatchString(`\w`, sobject)
	if err != nil {
		return UpdatedRecords{}, err
	}

	if matching == false {
		return UpdatedRecords{}, fmt.Errorf("sobject salesforce api: %s is not a valid sobject", sobject)
	}

	return r.query.updatedRecordsCallout(sobject, startDate, endDate)
}

// GetContent returns the blob from a content SObject.
func (r *Resources) GetContent(id string, content ContentType) ([]byte, error) {
	if r.query == nil {
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

	return r.query.contentCallout(id, content)
}
