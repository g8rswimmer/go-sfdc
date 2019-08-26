package common

import (
	"github.com/g8rswimmer/go-sfdc/session"
	"github.com/g8rswimmer/go-sfdc/sobject"
)

// DML is the structure for all of the DML REST API operations.
type DML struct {
	Metadata *Metadata `json:"metadata,omitempty"`
	Describe *Describe `json:"describe,omitempty"`
	Insert   *Insert   `json:"insert,omitempty"`
	Update   *Update   `json:"update,omitempty"`
	Upsert   *Upsert   `json:"upsert,omitempty"`
	Delete   *Delete   `json:"delete,omitempty"`
}

// Run will, if present, run the different operations.
func (d *DML) Run(session *session.Session) {
	sobject, err := sobject.NewResources(session)
	Error(err)

	if d.Metadata != nil {
		d.Metadata.Run(sobject)
	}

	if d.Describe != nil {
		d.Describe.Run(sobject)
	}

	if d.Insert != nil {
		d.Insert.Run(sobject)
	}

	if d.Update != nil {
		d.Update.Run(sobject)
	}

	if d.Upsert != nil {
		d.Upsert.Run(sobject)
	}

	if d.Delete != nil {
		d.Delete.Run(sobject)
	}
}

// Metadata is the structure to retrieve the metadata of a sobject.
type Metadata struct {
	SObject string `json:"sobject"`
}

// Run will retreive the metadata.
func (m *Metadata) Run(sobject *sobject.Resources) {
	metadata, err := sobject.Metadata(m.SObject)
	Error(err)

	Info("%s Metadata", m.SObject)
	Info("------------------------")
	Info("%+v", metadata)
	Info("------------------------")
}

// Describe is the structure to retrieve the description of the sobject.
type Describe struct {
	SObject string `json:"sobject"`
}

// Run will retrieve the description.
func (d *Describe) Run(sobject *sobject.Resources) {
	describe, err := sobject.Describe(d.SObject)
	Error(err)

	Info("%s Describe", d.SObject)
	Info("------------------------")
	Info("%+v", describe)
	Info("------------------------")
}

// Insert is the structure for insert fields of a sobject.
type Insert struct {
	SObj       string                 `json:"sobject"`
	SObjFields map[string]interface{} `json:"fields"`
}

// SObject is the sobject to insert.
func (i *Insert) SObject() string {
	return i.SObj
}

// Fields is the set of fields and values to insert.
func (i *Insert) Fields() map[string]interface{} {
	return i.SObjFields
}

// Run will insert a sobject using the REST API.
func (i *Insert) Run(sobject *sobject.Resources) {
	value, err := sobject.Insert(i)
	Error(err)

	Info("%s Insert", i.SObj)
	Info("------------------------")
	Info("%+v", value)
	Info("------------------------")
}

// Update is the structure to update a sobject.
type Update struct {
	Insert
	SObjID string `json:"id"`
}

// SObject is the sobject to update.
func (u *Update) SObject() string {
	return u.SObj
}

// Fields is the set of fields an values to update.
func (u *Update) Fields() map[string]interface{} {
	return u.SObjFields
}

// ID is the sobject id.
func (u *Update) ID() string {
	return u.SObjID
}

// Run will update a sobject using the REST API.
func (u *Update) Run(sobject *sobject.Resources) {
	err := sobject.Update(u)
	Error(err)

	Info("%s %s Updated", u.SObj, u.SObjID)
	Info("------------------------")
}

// Upsert is the structure to upsert a sobject.
type Upsert struct {
	Update
	SObjExternalField string `json:"external_id"`
}

// SObject is the sobject to upsert.
func (u *Upsert) SObject() string {
	return u.SObj
}

// Fields is the set of fields an values to upsert.
func (u *Upsert) Fields() map[string]interface{} {
	return u.SObjFields
}

// ID is the sobject's external id
func (u *Upsert) ID() string {
	return u.SObjID
}

// ExternalField is the sobject's external field to use in the upsert.
func (u *Upsert) ExternalField() string {
	return u.SObjExternalField
}

// Run will upsert a sobject using REST API
func (u *Upsert) Run(sobject *sobject.Resources) {
	value, err := sobject.Upsert(u)
	Error(err)

	Info("%s Upsert", u.SObj)
	Info("------------------------")
	Info("%+v", value)
	Info("------------------------")
}

// Delete is the structure to delete an sobject
type Delete struct {
	SObjID string `json:"id"`
	SObj   string `json:"sobject"`
}

// SObject is the sobject to delete.
func (d *Delete) SObject() string {
	return d.SObj
}

// ID is the sobject id.
func (d *Delete) ID() string {
	return d.SObjID
}

// Run will delete a sobject using REST API
func (d *Delete) Run(sobject *sobject.Resources) {
	err := sobject.Delete(d)
	Error(err)

	Info("%s %s Deleted", d.SObj, d.SObjID)
	Info("------------------------")
}
