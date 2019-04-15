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

type CollectionDelete struct {
	session session.ServiceFormatter
	records []string
}

func (cd *CollectionDelete) Delete(allOrNone bool) ([]DeleteValue, error) {
	c := &collection{
		method:   http.MethodDelete,
		endpoint: cd.session.ServiceURL() + endpoint,
		values:   cd.values(allOrNone),
	}
	var values []DeleteValue
	err := c.send(cd.session, &values)
	if err != nil {
		return nil, err
	}
	return values, nil
}
func (cd *CollectionDelete) Records(records ...string) {
	if cd == nil {
		panic("collections: Collection Delete can not be nil")
	}
	cd.records = append(cd.records, records...)
}
func (cd *CollectionDelete) values(allOrNone bool) *url.Values {
	values := &url.Values{}
	values.Add("ids", strings.Join(cd.records, ","))
	values.Add("allOrNone", fmt.Sprintf("%t", allOrNone))
	return values
}
