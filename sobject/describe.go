package sobject

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TuSKan/go-sfdc"
	"github.com/TuSKan/go-sfdc/session"
)

// DescribeValue is a structure that is returned from the from the Salesforce
// API SObject describe.
type DescribeValue struct {
	Activateable         bool                `json:"activateable"`
	CompactLayoutable    bool                `json:"compactLayoutable"`
	Createable           bool                `json:"createable"`
	Custom               bool                `json:"custom"`
	CustomSetting        bool                `json:"customSetting"`
	Deletable            bool                `json:"deletable"`
	DeprecatedAndHidden  bool                `json:"deprecatedAndHidden"`
	FeedEnabled          bool                `json:"feedEnabled"`
	HasSubtypes          bool                `json:"hasSubtypes"`
	IsSubType            bool                `json:"isSubtype"`
	Layoutable           bool                `json:"layoutable"`
	Mergeable            bool                `json:"mergeable"`
	MRUEnabled           bool                `json:"mruEnabled"`
	Queryable            bool                `json:"queryable"`
	Replicateable        bool                `json:"replicateable"`
	Retrieveable         bool                `json:"retrieveable"`
	SearchLayoutable     bool                `json:"searchLayoutable"`
	Searchable           bool                `json:"searchable"`
	Triggerable          bool                `json:"triggerable"`
	Undeletable          bool                `json:"undeletable"`
	Updateable           bool                `json:"updateable"`
	KeyPrefix            string              `json:"keyPrefix"`
	Label                string              `json:"label"`
	LabelPural           string              `json:"labelPlural"`
	Name                 string              `json:"name"`
	NetworkScopeFielName string              `json:"networkScopeFieldName"`
	Listviewable         interface{}         `json:"listviewable"`
	LookupLayoutable     interface{}         `json:"lookupLayoutable"`
	URLs                 ObjectURLs          `json:"urls"`
	ActionOverrides      []ActionOverride    `json:"actionOverrides"`
	ChildRelationships   []ChildRelationship `json:"childRelationships"`
	Fields               []Field             `json:"fields"`
	RecordTypeInfos      []RecordTypeInfo    `json:"recordTypeInfos"`
	SupportedScopes      []SupportedScope    `json:"supportedScopes"`
	NamedLayoutInfos     []interface{}       `json:"namedLayoutInfos"`
}

// ActionOverride describes the objects overrides.
type ActionOverride struct {
	IsAvailableInTouch bool   `json:"isAvailableInTouch"`
	FormFactor         string `json:"formFactor"`
	Name               string `json:"name"`
	PageID             string `json:"pageId"`
	URL                string `json:"url"`
}

// ChildRelationship describes the child relationship of the SObject.
type ChildRelationship struct {
	CascadeDelete       bool     `json:"cascadeDelete"`
	RestrictedDelete    bool     `json:"restrictedDelete"`
	DeprecatedAndHidden bool     `json:"deprecatedAndHidden"`
	ChildSObject        string   `json:"childSObject"`
	Field               string   `json:"field"`
	RelationshipName    string   `json:"relationshipName"`
	JunctionIDListNames []string `json:"junctionIdListNames"`
	JunctionReferenceTo []string `json:"junctionReferenceTo"`
}

// PickListValue describes the SObject's field picklist values.
type PickListValue struct {
	Active       bool   `json:"active"`
	DefaultValue bool   `json:"defaultValue"`
	Label        string `json:"label"`
	ValidFor     string `json:"validFor"`
	Value        string `json:"value"`
}

// Field describes the SOBject's fields.
type Field struct {
	Aggregatable                 bool            `json:"aggregatable"`
	AIPredictionField            bool            `json:"aiPredictionField"`
	AutoNumber                   bool            `json:"autoNumber"`
	Calculated                   bool            `json:"calculated"`
	CascadeDelete                bool            `json:"cascadeDelete"`
	CaseSensitive                bool            `json:"caseSensitive"`
	Createable                   bool            `json:"createable"`
	Custom                       bool            `json:"custom"`
	DefaultedOnCreate            bool            `json:"defaultedOnCreate"`
	DependentPicklist            bool            `json:"dependentPicklist"`
	DeprecatedAndHidden          bool            `json:"deprecatedAndHidden"`
	DisplayLocationInDecimal     bool            `json:"displayLocationInDecimal"`
	Encrypted                    bool            `json:"encrypted"`
	ExternalID                   bool            `json:"externalId"`
	Filterable                   bool            `json:"filterable"`
	FormulaTreatNullNumberAsZero bool            `json:"formulaTreatNullNumberAsZero"`
	Groupable                    bool            `json:"groupable"`
	HighScaleNumber              bool            `json:"highScaleNumber"`
	HTMLFormatted                bool            `json:"htmlFormatted"`
	IDLookup                     bool            `json:"idLookup"`
	NameField                    bool            `json:"nameField"`
	NamePointing                 bool            `json:"namePointing"`
	Nillable                     bool            `json:"nillable"`
	Permissionable               bool            `json:"permissionable"`
	PolymorphicForeignKey        bool            `json:"polymorphicForeignKey"`
	QueryByDistance              bool            `json:"queryByDistance"`
	RestrictedDelete             bool            `json:"restrictedDelete"`
	RestrictedPicklist           bool            `json:"restrictedPicklist"`
	SearchPrefilterable          bool            `json:"searchPrefilterable"`
	Sortable                     bool            `json:"sortable"`
	Unique                       bool            `json:"unique"`
	Updateable                   bool            `json:"updateable"`
	WriteRequiredMasterRead      bool            `json:"writeRequiresMasterRead"`
	Digits                       int             `json:"digits"`
	Length                       int             `json:"length"`
	Precision                    int             `json:"precision"`
	ByteLength                   int             `json:"byteLength"`
	Scale                        int             `json:"scale"`
	InlineHelpText               string          `json:"inlineHelpText"`
	Label                        string          `json:"label"`
	Name                         string          `json:"name"`
	RelationshipName             string          `json:"relationshipName"`
	Type                         string          `json:"type"`
	SoapType                     string          `json:"soapType"`
	CompoundFieldName            string          `json:"compoundFieldName"`
	ControllerName               string          `json:"controllerName"`
	ReferenceTargetField         string          `json:"referenceTargetField"`
	ReferenceTo                  []string        `json:"referenceTo"`
	CalculatedFormula            interface{}     `json:"calculatedFormula"`
	DefaultValue                 interface{}     `json:"defaultValue"`
	DefaultValueFormula          interface{}     `json:"defaultValueFormula"`
	ExtraTypeInfo                interface{}     `json:"extraTypeInfo"`
	FilteredLookupInfo           interface{}     `json:"filteredLookupInfo"`
	Mask                         interface{}     `json:"mask"`
	MaskType                     interface{}     `json:"maskType"`
	RelationshipOrder            interface{}     `json:"relationshipOrder"`
	PicklistValues               []PickListValue `json:"picklistValues"`
}

// RecordTypeInfo describes the SObjects record types assocaited with it.
type RecordTypeInfo struct {
	Active                   bool          `json:"active"`
	Available                bool          `json:"available"`
	DefaultRecordTypeMapping bool          `json:"defaultRecordTypeMapping"`
	Master                   bool          `json:"master"`
	Name                     string        `json:"name"`
	RecordTypeID             string        `json:"recordTypeId"`
	DeveloperName            string        `json:"developerName"`
	URLs                     RecordTypeURL `json:"urls"`
}

// RecordTypeURL contains the record type's URLs.
type RecordTypeURL struct {
	Layout string `json:"layout"`
}

// SupportedScope describes the supported scope.
type SupportedScope struct {
	Label string `json:"label"`
	Name  string `json:"name"`
}

const describeEndpoint = "/describe"

type describe struct {
	session session.ServiceFormatter
}

func (d *describe) callout(sobject string) (DescribeValue, error) {

	request, err := d.request(sobject)

	if err != nil {
		return DescribeValue{}, err
	}

	value, err := d.response(request)

	if err != nil {
		return DescribeValue{}, err
	}

	return value, nil
}

func (d *describe) request(sobject string) (*http.Request, error) {
	url := d.session.ServiceURL() + objectEndpoint + sobject + describeEndpoint

	request, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/json")
	d.session.AuthorizationHeader(request)
	return request, nil

}

func (d *describe) response(request *http.Request) (DescribeValue, error) {
	response, err := d.session.Client().Do(request)

	if err != nil {
		return DescribeValue{}, err
	}

	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var respErrs []sfdc.Error
		err = decoder.Decode(&respErrs)
		var errMsg error
		if err == nil {
			for _, respErr := range respErrs {
				errMsg = fmt.Errorf("metadata response err: %s: %s", respErr.ErrorCode, respErr.Message)
			}
		} else {
			errMsg = fmt.Errorf("metadata response err: %d %s", response.StatusCode, response.Status)
		}

		return DescribeValue{}, errMsg
	}

	var value DescribeValue
	err = decoder.Decode(&value)
	if err != nil {
		return DescribeValue{}, err
	}

	return value, nil
}
