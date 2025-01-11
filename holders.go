package holders

import (
	"context"
	"sync"
)

type Asset struct {
	Name     string
	UnitName string
	AssetID  uint64
}

type AssetHolding struct {
	Name     string
	UnitName string
	Address  string
	Amount   uint64
	AssetID  uint64
}

type Collection struct {
	Name                    string
	Addresses               []string
	UnitNamePrefixes        []string
	ExcludedAssets          []uint64
	AssetIndexGreaterThan   uint64
	ExcludedHolderAddresses []string
	IncludeNameContains     []string
	ExcludeNameContains     []string
}

type CollectionClient interface {
	GetAssetHoldingsByCollection(ctx context.Context, collection Collection) ([]AssetHolding, error)
	GetAssetsByCollection(ctx context.Context, collection Collection) ([]Asset, error)
	IsAssetOwned(ctx context.Context, asset Asset) (bool, error)
}

func GetAssetHoldingsByCollection(ctx context.Context, client CollectionClient, collections []Collection, concurrency int) (map[string][]AssetHolding, error) {
	result := make(map[string][]AssetHolding)
	resultMutex := &sync.Mutex{}
	semaphore := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	errChan := make(chan error, len(collections))

	for _, collection := range collections {
		semaphore <- struct{}{}
		wg.Add(1)

		go func(collection Collection) {
			defer wg.Done()
			defer func() { <-semaphore }()

			assetHoldings, err := client.GetAssetHoldingsByCollection(ctx, collection)
			if err != nil {
				errChan <- err
				return
			}

			resultMutex.Lock()
			result[collection.Name] = assetHoldings
			resultMutex.Unlock()
		}(collection)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	return result, nil
}
