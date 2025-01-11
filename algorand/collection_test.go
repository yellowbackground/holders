package algorand_test

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
	raffle "github.com/yellowbackground/holders"
	"github.com/yellowbackground/holders/algorand"
	"github.com/yellowbackground/holders/testdata"
	"net/http"
	"testing"
)

func TestGetCollection(t *testing.T) {
	tests := map[string]struct {
		GotCollection           raffle.Collection
		GotCreatedAssetResponse string
		GotBalancesResponse     string
		WantHoldings            []raffle.AssetHolding
	}{
		"account holds created asset": {
			GotCollection: raffle.Collection{
				Addresses: []string{testdata.TestAccount1Address},
			},
			GotCreatedAssetResponse: `{
			  "created-assets": [
				{
				  "index": 1,
				  "params": {
					"unit-name": "BRO#1"
				  }
				}
			  ]
			}`,
			GotBalancesResponse: fmt.Sprintf(`{
			  "balances": [
				{
				  "address": "%s",
				  "amount": 1,
				  "deleted": false
				}
			  ]
			}`, testdata.TestAccount2Address),
			WantHoldings: []raffle.AssetHolding{
				{
					Address:  testdata.TestAccount2Address,
					Amount:   1,
					AssetID:  1,
					UnitName: "BRO#1",
				},
			},
		},
		"no matching unit name prefix": {
			GotCollection: raffle.Collection{
				Addresses:        []string{testdata.TestAccount1Address},
				UnitNamePrefixes: []string{"A"},
			},
			GotCreatedAssetResponse: `{
			  "created-assets": [
				{
				  "index": 1,
				  "params": {
					"unit-name": "B"
				  }
				}
			  ]
			}`,
			GotBalancesResponse: fmt.Sprintf(`{
			  "balances": [
				{
				  "address": "%s",
				  "amount": 1,
				  "deleted": false
				}
			  ]
			}`, testdata.TestAccount2Address),
			WantHoldings: []raffle.AssetHolding{},
		},
		"excluded asset id": {
			GotCollection: raffle.Collection{
				Addresses:      []string{testdata.TestAccount1Address},
				ExcludedAssets: []uint64{1},
			},
			GotCreatedAssetResponse: `{
			  "created-assets": [
				{
				  "index": 1,
				  "params": {
					"unit-name": "B"
				  }
				}
			  ]
			}`,
			GotBalancesResponse: fmt.Sprintf(`{
			  "balances": [
				{
				  "address": "%s",
				  "amount": 1,
				  "deleted": false
				}
			  ]
			}`, testdata.TestAccount2Address),
			WantHoldings: []raffle.AssetHolding{},
		},
		"excluded holder": {
			GotCollection: raffle.Collection{
				Addresses:               []string{testdata.TestAccount1Address},
				ExcludedHolderAddresses: []string{testdata.TestAccount2Address},
			},
			GotCreatedAssetResponse: `{
			  "created-assets": [
				{
				  "index": 1,
				  "params": {
					"unit-name": "B"
				  }
				}
			  ]
			}`,
			GotBalancesResponse: fmt.Sprintf(`{
			  "balances": [
				{
				  "address": "%s",
				  "amount": 1,
				  "deleted": false
				}
			  ]
			}`, testdata.TestAccount2Address),
			WantHoldings: []raffle.AssetHolding{},
		},
		"not excluded holder": {
			GotCollection: raffle.Collection{
				Addresses:               []string{testdata.TestAccount1Address},
				ExcludedHolderAddresses: []string{testdata.TestAccount1Address},
			},
			GotCreatedAssetResponse: `{
			  "created-assets": [
				{
				  "index": 1,
				  "params": {
					"unit-name": "B"
				  }
				}
			  ]
			}`,
			GotBalancesResponse: fmt.Sprintf(`{
			  "balances": [
				{
				  "address": "%s",
				  "amount": 1,
				  "deleted": false
				}
			  ]
			}`, testdata.TestAccount2Address),
			WantHoldings: []raffle.AssetHolding{
				{
					Address:  testdata.TestAccount2Address,
					Amount:   1,
					UnitName: "B",
					AssetID:  1,
				},
			},
		},
		"asset id not greater than": {
			GotCollection: raffle.Collection{
				Addresses:             []string{testdata.TestAccount1Address},
				AssetIndexGreaterThan: 1,
			},
			GotCreatedAssetResponse: `{
			  "created-assets": [
				{
				  "index": 1,
				  "params": {
					"unit-name": "B"
				  }
				}
			  ]
			}`,
			GotBalancesResponse: fmt.Sprintf(`{
			  "balances": [
				{
				  "address": "%s",
				  "amount": 1,
				  "deleted": false
				}
			  ]
			}`, testdata.TestAccount2Address),
			WantHoldings: []raffle.AssetHolding{},
		},
		"asset id greater than": {
			GotCollection: raffle.Collection{
				Addresses:             []string{testdata.TestAccount1Address},
				AssetIndexGreaterThan: 0,
			},
			GotCreatedAssetResponse: `{
			  "created-assets": [
				{
				  "index": 1,
				  "params": {
					"unit-name": "B"
				  }
				}
			  ]
			}`,
			GotBalancesResponse: fmt.Sprintf(`{
			  "balances": [
				{
				  "address": "%s",
				  "amount": 1,
				  "deleted": false
				}
			  ]
			}`, testdata.TestAccount2Address),
			WantHoldings: []raffle.AssetHolding{
				{
					Address:  testdata.TestAccount2Address,
					Amount:   1,
					AssetID:  1,
					UnitName: "B",
				},
			},
		},
		"excluded asset id - not matching": {
			GotCollection: raffle.Collection{
				Addresses:      []string{testdata.TestAccount1Address},
				ExcludedAssets: []uint64{2},
			},
			GotCreatedAssetResponse: `{
			  "created-assets": [
				{
				  "index": 1,
				  "params": {
					"unit-name": "B"
				  }
				}
			  ]
			}`,
			GotBalancesResponse: fmt.Sprintf(`{
			  "balances": [
				{
				  "address": "%s",
				  "amount": 1,
				  "deleted": false
				}
			  ]
			}`, testdata.TestAccount2Address),
			WantHoldings: []raffle.AssetHolding{
				{
					Address:  testdata.TestAccount2Address,
					Amount:   1,
					AssetID:  1,
					UnitName: "B",
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			nodeCli, _ := algod.MakeClient("http://localhost:8000", "")
			idxCli, _ := indexer.MakeClient("http://localhost:9000", "")
			underTest := algorand.NewCollectionClient(nodeCli, idxCli)

			resetTransport := setupMocks(test.GotCreatedAssetResponse, test.GotBalancesResponse).End()
			defer resetTransport()

			holdings, err := underTest.GetAssetHoldingsByCollection(context.Background(), test.GotCollection)

			assert.NoError(t, err)
			assert.Equal(t, test.WantHoldings, holdings)
		})
	}
}

func setupMocks(gotCreatedAssets string, gotBalances string) *apitest.StandaloneMocks {
	getCreatedAssetsMock := apitest.NewMock().
		Get("/v2/accounts/" + testdata.TestAccount1Address).
		RespondWith().
		Status(http.StatusOK).
		JSON(gotCreatedAssets).
		End()
	getBalancesMock := apitest.NewMock().
		Getf("/v2/assets/1/balances").
		RespondWith().
		Status(http.StatusOK).
		JSON(gotBalances).
		End()
	return apitest.NewStandaloneMocks(getCreatedAssetsMock, getBalancesMock)
}
