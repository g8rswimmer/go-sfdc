package sobject

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/g8rswimmer/goforce/session"
)

type DesribeValue struct {
	ActionOverrides      []ActionOverride    `json:"actionOverrides"`
	Activateable         bool                `json:"activateable"`
	ChildRelationships   []ChildRelationship `json:"childRelationships"`
	CompactLayoutable    bool                `json:"compactLayoutable"`
	Createable           bool                `json:"createable"`
	Custom               bool                `json:"custom"`
	CustomSetting        bool                `json:"customSetting"`
	Deletable            bool                `json:"deletable"`
	DeprecatedAndHidden  bool                `json:"deprecatedAndHidden"`
	FeedEnabled          bool                `json:"feedEnabled"`
	Fields               []Field             `json:"fields"`
	HasSubtypes          bool                `json:"hasSubtypes"`
	IsSubType            bool                `json:"isSubtype"`
	KeyPrefix            string              `json:"keyPrefix"`
	Label                string              `json:"label"`
	LabelPural           string              `json:"labelPlural"`
	Layoutable           bool                `json:"layoutable"`
	Listviewable         interface{}         `json:"listviewable"`
	LookupLayoutable     interface{}         `json:"lookupLayoutable"`
	Mergeable            bool                `json:"mergeable"`
	MRUEnabled           bool                `json:"mruEnabled"`
	Name                 string              `json:"name"`
	NamedLayoutInfos     []interface{}       `json:"namedLayoutInfos"`
	NetworkScopeFielName string              `json:"networkScopeFieldName"`
	Queryable            bool                `json:"queryable"`
	RecordTypeInfos      []RecordTypeInfo    `json:"recordTypeInfos"`
	Replicateable        bool                `json:"replicateable"`
	Retrieveable         bool                `json:"retrieveable"`
	SearchLayoutable     bool                `json:"searchLayoutable"`
	Searchable           bool                `json:"searchable"`
	SupportedScopes      []SupportedScope    `json:"supportedScopes"`
	Triggerable          bool                `json:"triggerable"`
	Undeletable          bool                `json:"undeletable"`
	Updateable           bool                `json:"updateable"`
	URLs                 []ObjectURLs        `json:"urls"`
}

type ActionOverride struct {
	FormFactor         string `json:"formFactor"`
	IsAvailableInTouch bool   `json:"isAvailableInTouch"`
	Name               string `json:"name"`
	PageID             string `json:"pageId"`
	URL                string `json:"url"`
}

type ChildRelationship struct {
	CascadeDelete       bool     `json:"cascadeDelete"`
	ChildSObject        string   `json:"childSObject"`
	DeprecatedAndHidden bool     `json:"deprecatedAndHidden"`
	Field               string   `json:"field"`
	JunctionIDListNames []string `json:"junctionIdListNames"`
	JunctionReferenceTo []string `json:"junctionReferenceTo"`
	RelationshipName    string   `json:"relationshipName"`
	RestrictedDelete    bool     `json:"restrictedDelete"`
}

type PickListValue struct {
	Active       bool   `json:"active"`
	DefaultValue bool   `json:"defaultValue"`
	Label        string `json:"label"`
	ValidFor     string `json:"validFor"`
	Value        string `json:"value"`
}

type Field struct {
	Aggregatable                 bool            `json:"aggregatable"`
	AIPredictionField            bool            `json:"aiPredictionField"`
	AutoNumber                   bool            `json:"autoNumber"`
	ByteLength                   int             `json:"byteLength"`
	Calculated                   bool            `json:"calculated"`
	CalculatedFormula            interface{}     `json:"calculatedFormula"`
	CascadeDelete                bool            `json:"cascadeDelete"`
	CaseSensitive                bool            `json:"caseSensitive"`
	CompoundFieldName            string          `json:"compoundFieldName"`
	ControllerName               string          `json:"controllerName"`
	Createable                   bool            `json:"createable"`
	Custom                       bool            `json:"custom"`
	DefaultValue                 interface{}     `json:"defaultValue"`
	DefaultValueFormula          interface{}     `json:"defaultValueFormula"`
	DefaultedOnCreate            bool            `json:"defaultedOnCreate"`
	DependentPicklist            bool            `json:"dependentPicklist"`
	DeprecatedAndHidden          bool            `json:"deprecatedAndHidden"`
	Digits                       int             `json:"digits"`
	DisplayLocationInDecimal     bool            `json:"displayLocationInDecimal"`
	Encrypted                    bool            `json:"encrypted"`
	ExternalID                   bool            `json:"externalId"`
	ExtraTypeInfo                interface{}     `json:"extraTypeInfo"`
	Filterable                   bool            `json:"filterable"`
	FilteredLookupInfo           interface{}     `json:"filteredLookupInfo"`
	FormulaTreatNullNumberAsZero bool            `json:"formulaTreatNullNumberAsZero"`
	Groupable                    bool            `json:"groupable"`
	HighScaleNumber              bool            `json:"highScaleNumber"`
	HTMLFormatted                bool            `json:"htmlFormatted"`
	IDLookup                     bool            `json:"idLookup"`
	InlineHelpText               string          `json:"inlineHelpText"`
	Label                        string          `json:"label"`
	Length                       int             `json:"length"`
	Mask                         interface{}     `json:"mask"`
	MakeType                     interface{}     `json:"maskType"`
	Name                         string          `json:"name"`
	NameField                    bool            `json:"nameField"`
	NamePointing                 bool            `json:"namePointing"`
	Nillable                     bool            `json:"nillable"`
	Permissionable               bool            `json:"permissionable"`
	PicklistValues               []PickListValue `json:"picklistValues"`
	PolymorphicForeignKey        bool            `json:"polymorphicForeignKey"`
	Precision                    int             `json:"precision"`
	QueryByDistance              bool            `json:"queryByDistance"`
	ReferenceTargetField         bool            `json:"referenceTargetField"`
	ReferenceTo                  []string        `json:"referenceTo"`
	RelationshipName             string          `json:"relationshipName"`
	RelationshipOrder            interface{}     `json:"relationshipOrder"`
	RestrictedDelete             bool            `json:"restrictedDelete"`
	RestrictedPicklist           bool            `json:"restrictedPicklist"`
	Scale                        int             `json:"scale"`
	SearchPrefilterable          bool            `json:"searchPrefilterable"`
	SoapType                     string          `json:"soapType"`
	Sortable                     bool            `json:"sortable"`
	Type                         string          `json:"type"`
	Unique                       bool            `json:"unique"`
	Updateable                   bool            `json:"updateable"`
	WriteRequiredMasterRead      bool            `json:"writeRequiresMasterRead"`
}

type RecordTypeInfo struct {
	Active                   bool          `json:"active"`
	Available                bool          `json:"available"`
	DefaultRecordTypeMapping bool          `json:"defaultRecordTypeMapping"`
	DeveloperName            string        `json:"developerName"`
	Master                   bool          `json:"master"`
	Name                     string        `json:"name"`
	RecordTypeID             string        `json:"recordTypeId"`
	URLs                     RecordTypeURL `json:"urls"`
}

type RecordTypeURL struct {
	Layout string `json:"layout"`
}

type SupportedScope struct {
	Label string `json:"label"`
	Name  string `json:"name"`
}

const describeEndpoint = "/describe"

type describe struct {
	session session.Formatter
}

func (d *describe) Describe(sobject string) (DesribeValue, error) {

	request, err := d.request(sobject)

	if err != nil {
		return DesribeValue{}, err
	}

	value, err := d.response(request)

	if err != nil {
		return DesribeValue{}, err
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

func (d *describe) response(request *http.Request) (DesribeValue, error) {
	response, err := d.session.Client().Do(request)

	if err != nil {
		return DesribeValue{}, err
	}

	if response.StatusCode != http.StatusOK {
		return DesribeValue{}, fmt.Errorf("metadata response err: %d %s", response.StatusCode, response.Status)
	}

	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	var value DesribeValue
	err = decoder.Decode(&value)
	if err != nil {
		return DesribeValue{}, err
	}

	return value, nil
}
