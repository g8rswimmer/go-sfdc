package collections

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/nettolicious/go-sfdc/session"
	"github.com/nettolicious/go-sfdc/sobject"
)

// DeleteValue is the return value from the
// Salesforce API.
type DeleteValue struct {
	sobject.InsertValue
}

type remove struct {
	session session.ServiceFormatter
}

func (r *remove) callout(allOrNone bool, records []string) ([]DeleteValue, error) {
	if r == nil {
		panic("collections: Collection Delete can not be nil")
	}
	c := &collection{
		method:   http.MethodDelete,
		endpoint: endpoint,
		values:   r.values(allOrNone, records),
	}
	var values []DeleteValue
	err := c.send(r.session, &values)
	if err != nil {
		return nil, err
	}
	return values, nil
}
func (r *remove) values(allOrNone bool, records []string) *url.Values {
	values := &url.Values{}
	values.Add("ids", strings.Join(records, ","))
	values.Add("allOrNone", fmt.Sprintf("%t", allOrNone))
	return values
}
