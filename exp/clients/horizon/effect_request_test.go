package horizonclient

import (
	"context"
	"testing"
	"time"

	"github.com/stellar/go/protocols/horizon/effects"
	"github.com/stellar/go/support/http/httptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEffectRequestBuildUrl(t *testing.T) {
	er := EffectRequest{}
	endpoint, err := er.BuildUrl()

	// It should return valid all effects endpoint and no errors
	require.NoError(t, err)
	assert.Equal(t, "effects", endpoint)

	er = EffectRequest{ForAccount: "GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU"}
	endpoint, err = er.BuildUrl()

	// It should return valid account effects endpoint and no errors
	require.NoError(t, err)
	assert.Equal(t, "accounts/GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU/effects", endpoint)

	er = EffectRequest{ForLedger: "123"}
	endpoint, err = er.BuildUrl()

	// It should return valid ledger effects endpoint and no errors
	require.NoError(t, err)
	assert.Equal(t, "ledgers/123/effects", endpoint)

	er = EffectRequest{ForOperation: "123"}
	endpoint, err = er.BuildUrl()

	// It should return valid operation effects endpoint and no errors
	require.NoError(t, err)
	assert.Equal(t, "operations/123/effects", endpoint)

	er = EffectRequest{ForTransaction: "123"}
	endpoint, err = er.BuildUrl()

	// It should return valid transaction effects endpoint and no errors
	require.NoError(t, err)
	assert.Equal(t, "transactions/123/effects", endpoint)

	er = EffectRequest{ForLedger: "123", ForOperation: "789"}
	endpoint, err = er.BuildUrl()

	// error case: too many parameters for building any effect endpoint
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "Invalid request. Too many parameters")
	}

	er = EffectRequest{Cursor: "123456", Limit: 30, Order: OrderAsc}
	endpoint, err = er.BuildUrl()
	// It should return valid all effects endpoint with query params and no errors
	require.NoError(t, err)
	assert.Equal(t, "effects?cursor=123456&limit=30&order=asc", endpoint)

}

func TestEffectStream(t *testing.T) {
	hmock := httptest.NewClient()
	client := &Client{
		HorizonURL: "https://localhost/",
		HTTP:       hmock,
	}

	// All effects
	effectRequest := EffectRequest{}
	ctx, cancel := context.WithCancel(context.Background())

	hmock.On(
		"GET",
		"https://localhost/effects?cursor=now",
	).ReturnString(200, effectStreamResponse)

	go func() {
		// Stop streaming after 1 second.
		time.Sleep(1 * time.Second)
		cancel()
	}()

	var effectStream []effects.Base
	err := client.Stream(ctx, effectRequest, func(effect interface{}) {
		resp, ok := effect.(effects.Base)
		if ok {
			effectStream = append(effectStream, resp)
		}
	})

	if assert.NoError(t, err) {
		assert.Equal(t, effectStream[0].Type, "account_credited")
	}

	// Account effects
	effectRequest = EffectRequest{ForAccount: "GBNZN27NAOHRJRCMHQF2ZN2F6TAPVEWKJIGZIRNKIADWIS2HDENIS6CI"}
	ctx, cancel = context.WithCancel(context.Background())

	hmock.On(
		"GET",
		"https://localhost/accounts/GBNZN27NAOHRJRCMHQF2ZN2F6TAPVEWKJIGZIRNKIADWIS2HDENIS6CI/effects?cursor=now",
	).ReturnString(200, effectStreamResponse)

	go func() {
		// Stop streaming after 1 second.
		time.Sleep(1 * time.Second)
		cancel()
	}()

	err = client.Stream(ctx, effectRequest, func(effect interface{}) {
		resp, ok := effect.(effects.Base)
		if ok {
			effectStream = append(effectStream, resp)
		}
	})

	if assert.NoError(t, err) {
		assert.Equal(t, effectStream[0].Account, "GBNZN27NAOHRJRCMHQF2ZN2F6TAPVEWKJIGZIRNKIADWIS2HDENIS6CI")
	}

	// test error
	effectRequest = EffectRequest{}
	ctx, cancel = context.WithCancel(context.Background())

	hmock.On(
		"GET",
		"https://localhost/effects?cursor=now",
	).ReturnString(500, effectStreamResponse)

	err = client.Stream(ctx, effectRequest, func(effect interface{}) {
		resp, ok := effect.(effects.Base)
		if ok {
			effectStream = append(effectStream, resp)
		}
		cancel()
	})

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "Got bad HTTP status code 500")

	}
}

var effectStreamResponse = `data: {"_links":{"operation":{"href":"https://horizon-testnet.stellar.org/operations/2531135896703017"},"succeeds":{"href":"https://horizon-testnet.stellar.org/effects?order=desc\u0026cursor=2531135896703017-1"},"precedes":{"href":"https://horizon-testnet.stellar.org/effects?order=asc\u0026cursor=2531135896703017-1"}},"id":"0002531135896703017-0000000001","paging_token":"2531135896703017-1","account":"GBNZN27NAOHRJRCMHQF2ZN2F6TAPVEWKJIGZIRNKIADWIS2HDENIS6CI","type":"account_credited","type_i":2,"created_at":"2019-04-03T10:14:17Z","asset_type":"credit_alphanum4","asset_code":"qwop","asset_issuer":"GBM4HXXNDBWWQBXOL4QCTZIUQAP6XFUI3FPINUGUPBMULMTEHJPIKX6T","amount":"0.0460000"}
`
