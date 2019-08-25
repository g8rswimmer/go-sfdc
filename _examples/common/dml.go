package common

import (
	"github.com/g8rswimmer/go-sfdc/session"
	"github.com/g8rswimmer/go-sfdc/sobject"
)

type DML struct {
	Metadata *Metadata `json:"metadata,omitempty"`
	Describe *Describe `json:"describe,omitempty"`
	Insert   *Insert   `json:"insert,omitempty"`
	Update   *Update   `json:"update,omitempty"`
	Upsert   *Upsert   `json:"upsert,omitempty"`
	Delete   *Delete   `json:"delete,omitempty"`
}

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

type Metadata struct {
	SObject string `json:"sobject"`
}

func (m *Metadata) Run(sobject *sobject.Resources) {
	metadata, err := sobject.Metadata(m.SObject)
	Error(err)

	Info("%s Metadata", m.SObject)
	Info("------------------------")
	Info("%+v", metadata)
	Info("------------------------")
}

type Describe struct {
	SObject string `json:"sobject"`
}

func (d *Describe) Run(sobject *sobject.Resources) {
	describe, err := sobject.Describe(d.SObject)
	Error(err)

	Info("%s Describe", d.SObject)
	Info("------------------------")
	Info("%+v", describe)
	Info("------------------------")
}

type Insert struct {
	SObj       string                 `json:"sobject"`
	SObjFields map[string]interface{} `json:"fields"`
}

func (i *Insert) SObject() string {
	return i.SObj
}

func (i *Insert) Fields() map[string]interface{} {
	return i.SObjFields
}

func (i *Insert) Run(sobject *sobject.Resources) {
	value, err := sobject.Insert(i)
	Error(err)

	Info("%s Insert", i.SObj)
	Info("------------------------")
	Info("%+v", value)
	Info("------------------------")
}

type Update struct {
	Insert
	SObjID string `json:"id"`
}

func (u *Update) SObject() string {
	return u.SObj
}

func (u *Update) Fields() map[string]interface{} {
	return u.SObjFields
}

func (u *Update) ID() string {
	return u.SObjID
}

func (u *Update) Run(sobject *sobject.Resources) {
	err := sobject.Update(u)
	Error(err)

	Info("%s %s Updated", u.SObj, u.SObjID)
	Info("------------------------")
}

type Upsert struct {
	Update
	SObjExternalField string `json:"external_id"`
}

func (u *Upsert) SObject() string {
	return u.SObj
}

func (u *Upsert) Fields() map[string]interface{} {
	return u.SObjFields
}

func (u *Upsert) ID() string {
	return u.SObjID
}

func (u *Upsert) ExternalField() string {
	return u.SObjExternalField
}

func (u *Upsert) Run(sobject *sobject.Resources) {
	value, err := sobject.Upsert(u)
	Error(err)

	Info("%s Upsert", u.SObj)
	Info("------------------------")
	Info("%+v", value)
	Info("------------------------")
}

type Delete struct {
	SObjID string `json:"id"`
	SObj   string `json:"sobject"`
}

func (d *Delete) SObject() string {
	return d.SObj
}

func (d *Delete) ID() string {
	return d.SObjID
}

func (d *Delete) Run(sobject *sobject.Resources) {
	err := sobject.Delete(d)
	Error(err)

	Info("%s %s Deleted", d.SObj, d.SObjID)
	Info("------------------------")
}
