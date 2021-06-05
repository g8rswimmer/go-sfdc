package sobject

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/nettolicious/go-sfdc/session"
)

func Test_describe_Describe(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
	}
	type args struct {
		sobject string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    DescribeValue
		wantErr bool
	}{
		{
			name: "Request Error",
			fields: fields{
				session: &mockSessionFormatter{
					url: "123://wrong",
				},
			},
			args: args{
				sobject: "someobject",
			},
			want:    DescribeValue{},
			wantErr: true,
		},
		{
			name: "Response HTTP Error",
			fields: fields{
				session: &mockSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {

						resp := `[
							{
								"message" : "The requested resource does not exist",
								"errorCode" : "NOT_FOUND"
							}							
						]
						`
						return &http.Response{
							StatusCode: 500,
							Status:     "Some Status",
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				sobject: "someobject",
			},
			want:    DescribeValue{},
			wantErr: true,
		},
		{
			name: "Response JSON Error",
			fields: fields{
				session: &mockSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {
						resp := `
						{`

						return &http.Response{
							StatusCode: 200,
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				sobject: "someobject",
			},
			want:    DescribeValue{},
			wantErr: true,
		},
		{
			name: "Response Passing",
			fields: fields{
				session: &mockSessionFormatter{
					url: "https://test.salesforce.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {
						resp := `
						{
							"actionOverrides": [
								{
									"formFactor": "LARGE",
									"isAvailableInTouch": false,
									"name": "New",
									"pageId": "0Abm00000009E2TCAU",
									"url": "some.url"
								}
							],
							"activateable": false,
							"childRelationships": [
								{
									"cascadeDelete": false,
									"childSObject": "AcceptedEventRelation",
									"deprecatedAndHidden": false,
									"field": "RelationId",
									"junctionIdListNames": [],
									"junctionReferenceTo": [],
									"relationshipName": "PersonAcceptedEventRelations",
									"restrictedDelete": false
								}
							],
							"compactLayoutable": true,
							"createable": true,
							"custom": false,
							"customSetting": false,
							"deletable": true,
							"deprecatedAndHidden": false,
							"feedEnabled": true,
							"fields": [
								{
									"aggregatable": true,
									"aiPredictionField": false,
									"autoNumber": false,
									"byteLength": 18,
									"calculated": false,
									"calculatedFormula": null,
									"cascadeDelete": false,
									"caseSensitive": false,
									"compoundFieldName": null,
									"controllerName": null,
									"createable": false,
									"custom": false,
									"defaultValue": null,
									"defaultValueFormula": null,
									"defaultedOnCreate": true,
									"dependentPicklist": false,
									"deprecatedAndHidden": false,
									"digits": 0,
									"displayLocationInDecimal": false,
									"encrypted": false,
									"externalId": false,
									"extraTypeInfo": null,
									"filterable": true,
									"filteredLookupInfo": null,
									"formulaTreatNullNumberAsZero": false,
									"groupable": true,
									"highScaleNumber": false,
									"htmlFormatted": false,
									"idLookup": true,
									"inlineHelpText": null,
									"label": "Account ID",
									"length": 18,
									"mask": null,
									"maskType": null,
									"name": "Id",
									"nameField": false,
									"namePointing": false,
									"nillable": false,
									"permissionable": false,
									"picklistValues": [
										{
											"active": true,
											"defaultValue": false,
											"label": "Hand surgery - plastic",
											"validFor": null,
											"value": "Hand surgery - plastic"
										}
									],
									"polymorphicForeignKey": false,
									"precision": 0,
									"queryByDistance": false,
									"referenceTargetField": null,
									"referenceTo": [],
									"relationshipName": null,
									"relationshipOrder": null,
									"restrictedDelete": false,
									"restrictedPicklist": false,
									"scale": 0,
									"searchPrefilterable": false,
									"soapType": "tns:ID",
									"sortable": true,
									"type": "id",
									"unique": false,
									"updateable": false,
									"writeRequiresMasterRead": false
								}
							],
							"hasSubtypes": true,
							"isSubtype": false,
							"keyPrefix": "001",
							"label": "Account",
							"labelPlural": "Accounts",
							"layoutable": true,
							"listviewable": null,
							"lookupLayoutable": null,
							"mergeable": true,
							"mruEnabled": true,
							"name": "Account",
							"namedLayoutInfos": [],
							"networkScopeFieldName": null,
							"queryable": true,
							"recordTypeInfos": [
								{
									"active": true,
									"available": true,
									"defaultRecordTypeMapping": false,
									"developerName": "DeveloperName",
									"master": false,
									"name": "Some Record Name",
									"recordTypeId": "xxx1234",
									"urls": {
										"layout": "/services/data/v44.0/sobjects/Account/describe/layouts/xxx1234"
									}
								}
							],
							"replicateable": true,
							"retrieveable": true,
							"searchLayoutable": true,
							"searchable": true,
							"supportedScopes": [
								{
									"label": "All accounts",
									"name": "everything"
								}
							],
							"triggerable": true,
							"undeletable": true,
							"updateable": true,
							"urls": {
								"compactLayouts": "/services/data/v44.0/sobjects/Account/describe/compactLayouts",
								"rowTemplate": "/services/data/v44.0/sobjects/Account/{ID}",
								"approvalLayouts": "/services/data/v44.0/sobjects/Account/describe/approvalLayouts",
								"uiDetailTemplate": "https:/my.salesforce.com/{ID}",
								"uiEditTemplate": "https://my.salesforce.com/{ID}/e",
								"defaultValues": "/services/data/v44.0/sobjects/Account/defaultValues?recordTypeId&fields",
								"listviews": "/services/data/v44.0/sobjects/Account/listviews",
								"describe": "/services/data/v44.0/sobjects/Account/describe",
								"uiNewRecord": "https://my.salesforce.com/001/e",
								"quickActions": "/services/data/v44.0/sobjects/Account/quickActions",
								"layouts": "/services/data/v44.0/sobjects/Account/describe/layouts",
								"sobject": "/services/data/v44.0/sobjects/Account"
							}
						}`

						return &http.Response{
							StatusCode: 200,
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
			},
			args: args{
				sobject: "someobject",
			},
			want: DescribeValue{
				ActionOverrides: []ActionOverride{
					{
						FormFactor: "LARGE",
						Name:       "New",
						PageID:     "0Abm00000009E2TCAU",
						URL:        "some.url",
					},
				},
				ChildRelationships: []ChildRelationship{
					{
						ChildSObject:        "AcceptedEventRelation",
						Field:               "RelationId",
						JunctionIDListNames: make([]string, 0),
						JunctionReferenceTo: make([]string, 0),
						RelationshipName:    "PersonAcceptedEventRelations",
					},
				},
				CompactLayoutable: true,
				Createable:        true,
				Deletable:         true,
				FeedEnabled:       true,
				Fields: []Field{
					{
						Aggregatable:        true,
						ByteLength:          18,
						CalculatedFormula:   nil,
						CompoundFieldName:   "",
						ControllerName:      "",
						DefaultValue:        nil,
						DefaultValueFormula: nil,
						DefaultedOnCreate:   true,
						Digits:              0,
						ExtraTypeInfo:       nil,
						Filterable:          true,
						FilteredLookupInfo:  nil,
						Groupable:           true,
						IDLookup:            true,
						InlineHelpText:      "",
						Label:               "Account ID",
						Length:              18,
						Mask:                nil,
						MaskType:            nil,
						Name:                "Id",
						PicklistValues: []PickListValue{
							{
								Active:   true,
								Label:    "Hand surgery - plastic",
								ValidFor: "",
								Value:    "Hand surgery - plastic",
							},
						},
						Precision:            0,
						ReferenceTargetField: "",
						ReferenceTo:          make([]string, 0),
						RelationshipName:     "",
						RelationshipOrder:    nil,
						Scale:                0,
						SoapType:             "tns:ID",
						Sortable:             true,
						Type:                 "id",
					},
				},
				HasSubtypes:          true,
				KeyPrefix:            "001",
				Label:                "Account",
				LabelPural:           "Accounts",
				Layoutable:           true,
				Listviewable:         nil,
				LookupLayoutable:     nil,
				Mergeable:            true,
				MRUEnabled:           true,
				Name:                 "Account",
				NamedLayoutInfos:     make([]interface{}, 0),
				NetworkScopeFielName: "",
				Queryable:            true,
				RecordTypeInfos: []RecordTypeInfo{
					{
						Active:        true,
						Available:     true,
						DeveloperName: "DeveloperName",
						Name:          "Some Record Name",
						RecordTypeID:  "xxx1234",
						URLs: RecordTypeURL{
							Layout: "/services/data/v44.0/sobjects/Account/describe/layouts/xxx1234",
						},
					},
				},
				Replicateable:    true,
				Retrieveable:     true,
				SearchLayoutable: true,
				Searchable:       true,
				SupportedScopes: []SupportedScope{
					{
						Label: "All accounts",
						Name:  "everything",
					},
				},
				Triggerable: true,
				Undeletable: true,
				Updateable:  true,
				URLs: ObjectURLs{
					CompactLayouts:   "/services/data/v44.0/sobjects/Account/describe/compactLayouts",
					RowTemplate:      "/services/data/v44.0/sobjects/Account/{ID}",
					ApprovalLayouts:  "/services/data/v44.0/sobjects/Account/describe/approvalLayouts",
					DefaultValues:    "/services/data/v44.0/sobjects/Account/defaultValues?recordTypeId&fields",
					ListViews:        "/services/data/v44.0/sobjects/Account/listviews",
					Describe:         "/services/data/v44.0/sobjects/Account/describe",
					QuickActions:     "/services/data/v44.0/sobjects/Account/quickActions",
					Layouts:          "/services/data/v44.0/sobjects/Account/describe/layouts",
					SObject:          "/services/data/v44.0/sobjects/Account",
					UIDetailTemplate: "https:/my.salesforce.com/{ID}",
					UIEditTemplate:   "https://my.salesforce.com/{ID}/e",
					UINewRecord:      "https://my.salesforce.com/001/e",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &describe{
				session: tt.fields.session,
			}
			got, err := d.callout(tt.args.sobject)
			if (err != nil) != tt.wantErr {
				t.Errorf("describe.Describe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("describe.Describe() = %v, want %v", got, tt.want)
			}
		})
	}
}
