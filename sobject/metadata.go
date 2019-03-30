package sobject

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/g8rswimmer/goforce/session"
)

// MetadataValue is the response from the SObject metadata API.
type MetadataValue struct {
	ObjectDescribe ObjectDescribe           `json:"objectDescribe"`
	RecentItems    []map[string]interface{} `json:"recentItems"`
}

// ObjectDescribe is the SObject metadata describe.
type ObjectDescribe struct {
	Activatable         bool       `json:"activateable"`
	Creatable           bool       `json:"createable"`
	Custom              bool       `json:"custom"`
	CustomSetting       bool       `json:"customSetting"`
	Deletable           bool       `json:"deletable"`
	DeprecatedAndHidden bool       `json:"deprecatedAndHidden"`
	FeedEnabled         bool       `json:"feedEnabled"`
	HasSubtype          bool       `json:"hasSubtypes"`
	IsSubtype           bool       `json:"isSubtype"`
	KeyPrefix           string     `json:"keyPrefix"`
	Label               string     `json:"label"`
	LabelPlural         string     `json:"labelPlural"`
	Layoutable          bool       `json:"layoutable"`
	Mergeable           bool       `json:"mergeable"`
	MruEnabled          bool       `json:"mruEnabled"`
	Name                string     `json:"name"`
	Queryable           bool       `json:"queryable"`
	Replicateable       bool       `json:"replicateable"`
	Retrieveable        bool       `json:"retrieveable"`
	Searchable          bool       `json:"searchable"`
	Triggerable         bool       `json:"triggerable"`
	Undeletable         bool       `json:"undeletable"`
	Updateable          bool       `json:"updateable"`
	URLs                ObjectURLs `json:"urls"`
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

type metadata struct {
	session session.Formatter
}

func (md *metadata) Metadata(sobject string) (MetadataValue, error) {

	request, err := md.request(sobject)

	if err != nil {
		return MetadataValue{}, err
	}

	value, err := md.response(request)

	if err != nil {
		return MetadataValue{}, err
	}

	return value, nil
}

func (md *metadata) request(sobject string) (*http.Request, error) {
	url := md.session.ServiceURL() + objectEndpoint + sobject

	request, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/json")
	md.session.AuthorizationHeader(request)
	return request, nil

}

func (md *metadata) response(request *http.Request) (MetadataValue, error) {
	response, err := md.session.Client().Do(request)

	if err != nil {
		return MetadataValue{}, err
	}

	if response.StatusCode != http.StatusOK {
		return MetadataValue{}, fmt.Errorf("metadata response err: %d %s", response.StatusCode, response.Status)
	}

	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	var value MetadataValue
	err = decoder.Decode(&value)
	if err != nil {
		return MetadataValue{}, err
	}

	return value, nil
}
