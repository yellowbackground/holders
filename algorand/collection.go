package algorand

import (
	"context"
	"errors"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/yellowbackground/holders"
	"strings"
)

func NewCollectionClient(algoD *algod.Client, idxClient *indexer.Client) holders.CollectionClient {
	return &collectionClient{
		algodClient:   algoD,
		indexerClient: idxClient,
	}
}

type collectionClient struct {
	indexerClient *indexer.Client
	algodClient   *algod.Client
}

// IsAssetOwned determines if an asset is owned by a regular holder - not an escrow or creator
func (c collectionClient) IsAssetOwned(ctx context.Context, asset holders.Asset) (bool, error) {
	assetDetails, err := c.algodClient.GetAssetByID(asset.AssetID).Do(ctx)
	if err != nil {
		return false, err
	}

	holder, err := c.getAssetHolder(ctx, asset)
	if err != nil {
		return false, err
	}

	if holder == assetDetails.Params.Creator {
		return false, nil
	}

	accountInfo, err := c.algodClient.AccountInformation(holder).Do(ctx)
	if err != nil {
		return false, err
	}

	// assume an account that is not opted into more than 1 asset is an escrow and is not eligible
	var optedInAssetCount int
	for _, asset := range accountInfo.Assets {
		if asset.Amount > 0 {
			optedInAssetCount++
		}
		if optedInAssetCount > 1 {
			return true, nil
		}
	}

	return false, nil
}

func (c collectionClient) getAssetHolder(ctx context.Context, asset holders.Asset) (string, error) {
	nextToken := ""
	for {
		res, err := c.indexerClient.LookupAssetBalances(asset.AssetID).
			Limit(1000).
			NextToken(nextToken).
			Do(ctx)
		if err != nil {
			return "", err
		}

		for _, a := range res.Balances {
			if a.Amount < 1 || a.Deleted || a.IsFrozen {
				continue
			}
			return a.Address, nil
		}

		if res.NextToken == "" {
			break
		}
		nextToken = res.NextToken
	}

	return "", errors.New("no holders found")
}

func (c collectionClient) GetAssetsByCollection(ctx context.Context, collection holders.Collection) ([]holders.Asset, error) {
	var createdAssets []holders.Asset

	for _, address := range collection.Addresses {
		accountInfo, err := c.algodClient.AccountInformation(address).Do(ctx)
		if err != nil {
			return nil, err
		}
		for _, asset := range accountInfo.CreatedAssets {
			if !matchesUnitNamePrefix(collection.UnitNamePrefixes, asset.Params.UnitName) {
				continue
			}
			if isExcludedAsset(collection.ExcludedAssets, asset.Index) {
				continue
			}
			if !isAssetIndexGreaterThan(collection.AssetIndexGreaterThan, asset.Index) {
				continue
			}
			if !nameContains(collection.IncludeNameContains, asset.Params.Name) {
				continue
			}
			if len(collection.ExcludeNameContains) > 0 && nameContains(collection.ExcludeNameContains, asset.Params.Name) {
				continue
			}
			createdAssets = append(createdAssets, holders.Asset{
				Name:     asset.Params.Name,
				UnitName: asset.Params.UnitName,
				AssetID:  asset.Index,
			})
		}
	}
	return createdAssets, nil
}

func (c collectionClient) GetAssetHoldingsByCollection(ctx context.Context, collection holders.Collection) ([]holders.AssetHolding, error) {
	createdAssets, err := c.GetAssetsByCollection(ctx, collection)
	if err != nil {
		return nil, err
	}

	holdings := []holders.AssetHolding{}
	for _, asset := range createdAssets {
		balancesResponse, err := c.indexerClient.LookupAssetBalances(asset.AssetID).Do(ctx)
		if err != nil {
			return nil, err
		}
		for _, balance := range balancesResponse.Balances {
			if balance.Amount > 0 &&
				balance.Deleted == false &&
				!isExcludedHolderAddress(collection.ExcludedHolderAddresses, balance.Address) {

				holdings = append(holdings, holders.AssetHolding{
					Address:  balance.Address,
					Amount:   balance.Amount,
					AssetID:  asset.AssetID,
					Name:     asset.Name,
					UnitName: asset.UnitName,
				})
			}
		}
	}

	return holdings, nil
}

func matchesUnitNamePrefix(unitNamePrefixes []string, unitName string) bool {
	if len(unitNamePrefixes) == 0 {
		return true
	}
	for _, prefix := range unitNamePrefixes {
		if strings.HasPrefix(strings.TrimSpace(unitName), strings.TrimSpace(prefix)) {
			return true
		}
	}
	return false
}

func isAssetIndexGreaterThan(assetIndexGreaterThan uint64, assetID uint64) bool {
	return assetID > assetIndexGreaterThan
}

func isExcludedAsset(excludedAssets []uint64, assetID uint64) bool {
	if len(excludedAssets) == 0 {
		return false
	}
	for _, excludedAsset := range excludedAssets {
		if excludedAsset == assetID {
			return true
		}
	}
	return false
}

func isExcludedHolderAddress(excludedHolderAddresses []string, address string) bool {
	if len(excludedHolderAddresses) == 0 {
		return false
	}
	for _, excludedHolderAddress := range excludedHolderAddresses {
		if excludedHolderAddress == address {
			return true
		}
	}
	return false
}

func nameContains(nameContains []string, name string) bool {
	if len(nameContains) == 0 {
		return true
	}
	for _, nameContainsString := range nameContains {
		if strings.Contains(strings.ToLower(name), strings.ToLower(nameContainsString)) {
			return true
		}
	}
	return false
}
