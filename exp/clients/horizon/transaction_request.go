package horizonclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	hProtocol "github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/support/errors"
)

// BuildUrl creates the endpoint to be queried based on the data in the TransactionRequest struct.
// If no data is set, it defaults to the build the URL for all transactions
func (tr TransactionRequest) BuildUrl() (endpoint string, err error) {
	nParams := countParams(tr.ForAccount, tr.ForLedger, tr.forTransactionHash)

	if nParams > 1 {
		return endpoint, errors.New("Invalid request. Too many parameters")
	}

	endpoint = "transactions"
	if tr.ForAccount != "" {
		endpoint = fmt.Sprintf("accounts/%s/transactions", tr.ForAccount)
	}
	if tr.ForLedger > 0 {
		endpoint = fmt.Sprintf("ledgers/%d/transactions", tr.ForLedger)
	}
	if tr.forTransactionHash != "" {
		endpoint = fmt.Sprintf("transactions/%s", tr.forTransactionHash)
	}

	queryParams := addQueryParams(cursor(tr.Cursor), limit(tr.Limit), tr.Order,
		includeFailed(tr.IncludeFailed))
	if queryParams != "" {
		endpoint = fmt.Sprintf("%s?%s", endpoint, queryParams)
	}

	_, err = url.Parse(endpoint)
	if err != nil {
		err = errors.Wrap(err, "failed to parse endpoint")
	}

	return endpoint, err
}

// TransactionHandler is a function that is called when a new transaction is received
type TransactionHandler func(hProtocol.Transaction)

// StreamTransactions streams executed transactions. It can be used to stream all transactions and  transactions for an account. Use context.WithCancel to stop streaming or context.Background() if you want
// to stream indefinitely. TransactionHandler is a user-supplied function that is executed for each streamed transaction received.
func (tr TransactionRequest) StreamTransactions(ctx context.Context, client *Client,
	handler TransactionHandler) (err error) {
	endpoint, err := tr.BuildUrl()
	if err != nil {
		return errors.Wrap(err, "Unable to build endpoint")
	}

	url := fmt.Sprintf("%s%s", client.getHorizonURL(), endpoint)

	return client.stream(ctx, url, func(data []byte) error {
		var transaction hProtocol.Transaction
		err = json.Unmarshal(data, &transaction)
		if err != nil {
			return errors.Wrap(err, "Error unmarshaling data")
		}
		handler(transaction)
		return nil
	})
}
