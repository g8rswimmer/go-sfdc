package collections

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/g8rswimmer/goforce"
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
	response, err := c.send(cd.session)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var insertErrs []goforce.Error
		err = decoder.Decode(&insertErrs)
		var errMsg error
		if err == nil {
			for _, insertErr := range insertErrs {
				errMsg = fmt.Errorf("insert response err: %s: %s", insertErr.ErrorCode, insertErr.Message)
			}
		} else {
			errMsg = fmt.Errorf("insert response err: %d %s", response.StatusCode, response.Status)
		}
		return nil, errMsg
	}
	var values []DeleteValue
	err = decoder.Decode(&values)
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
