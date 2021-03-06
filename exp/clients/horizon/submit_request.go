package horizonclient

import (
	"fmt"
	"net/url"

	"github.com/stellar/go/support/errors"
)

// BuildUrl returns the url for submitting transactions to a running horizon instance
func (sr submitRequest) BuildUrl() (endpoint string, err error) {
	if sr.endpoint == "" || sr.transactionXdr == "" {
		return endpoint, errors.New("Invalid request. Too few parameters")
	}

	query := url.Values{}
	query.Set("tx", sr.transactionXdr)

	endpoint = fmt.Sprintf("%s?%s", sr.endpoint, query.Encode())
	return endpoint, err
}
