package collections

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/g8rswimmer/goforce/session"
	"github.com/g8rswimmer/goforce/sobject"
)

type DeleteValue struct {
	sobject.InsertValue
}

type Delete struct {
	session session.ServiceFormatter
	records []string
}

func (d *Delete) Callout(allOrNone bool) ([]DeleteValue, error) {
	if d == nil {
		panic("collections: Collection Delete can not be nil")
	}
	c := &collection{
		method:   http.MethodDelete,
		endpoint: endpoint,
		values:   d.values(allOrNone),
	}
	var values []DeleteValue
	err := c.send(d.session, &values)
	if err != nil {
		return nil, err
	}
	return values, nil
}
func (d *Delete) Records(records ...string) {
	if d == nil {
		panic("collections: Collection Delete can not be nil")
	}
	d.records = append(d.records, records...)
}
func (d *Delete) values(allOrNone bool) *url.Values {
	values := &url.Values{}
	values.Add("ids", strings.Join(d.records, ","))
	values.Add("allOrNone", fmt.Sprintf("%t", allOrNone))
	return values
}
