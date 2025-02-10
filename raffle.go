package holders

import (
	"context"
	"fmt"
	"github.com/mroth/weightedrand/v2"
	"github.com/rs/zerolog/log"
)

type RaffleConfig struct {
	RandSeed string
}

type WeightedCollection struct {
	Collection Collection
	Weight     uint64
}

func RunWeightedCollectionRaffle(ctx context.Context, client CollectionClient, weightedCollections []WeightedCollection, numberOfWinners int, concurrency int, excludedWinnerWallets []string) ([]AssetHolding, error) {
	collections := extractCollections(weightedCollections)

	assetsHoldingsByCollection, err := GetAssetHoldingsByCollection(ctx, client, collections, concurrency)
	if err != nil {
		return nil, err
	}

	weightedTickets := createWeightedLotteryTickets(assetsHoldingsByCollection, weightedCollections)
	chooser, err := weightedrand.NewChooser(weightedTickets...)
	if err != nil {
		return nil, err
	}

	return pickUniqueWinners(weightedTickets, chooser, numberOfWinners, excludedWinnerWallets)
}

func extractCollections(weightedCollections []WeightedCollection) []Collection {
	collections := make([]Collection, len(weightedCollections))
	for i, wc := range weightedCollections {
		collections[i] = wc.Collection
	}
	return collections
}

func pickUniqueWinners(weightedTickets []weightedrand.Choice[AssetHolding, uint64], chooser *weightedrand.Chooser[AssetHolding, uint64], numberOfWinners int, excludedWallets []string) ([]AssetHolding, error) {
	selectedWinners := make(map[AssetHolding]bool)
	var winners []AssetHolding

	isExcludedWallet := func(wallet string) bool {
		for _, excluded := range excludedWallets {
			if wallet == excluded {
				return true
			}
		}
		return false
	}

	for len(winners) < numberOfWinners {
		if len(weightedTickets) == 0 {
			return nil, fmt.Errorf("not enough unique assets to select %d winners", numberOfWinners)
		}

		winner := chooser.Pick()

		if !isExcludedWallet(winner.Address) && !selectedWinners[winner] {
			selectedWinners[winner] = true
			winners = append(winners, winner)
		} else {
			// Remove the picked asset and recreate chooser
			weightedTickets = removePickedAsset(weightedTickets, winner)
			newChooser, err := weightedrand.NewChooser(weightedTickets...)
			if err != nil {
				return nil, err
			}
			chooser = newChooser
		}
	}

	return winners, nil
}

func removePickedAsset(tickets []weightedrand.Choice[AssetHolding, uint64], assetToRemove AssetHolding) []weightedrand.Choice[AssetHolding, uint64] {
	var updatedTickets []weightedrand.Choice[AssetHolding, uint64]
	for _, ticket := range tickets {
		if ticket.Item != assetToRemove {
			updatedTickets = append(updatedTickets, ticket)
		}
	}
	return updatedTickets
}

func createWeightedLotteryTickets(assetsByCollection map[string][]AssetHolding, collections []WeightedCollection) []weightedrand.Choice[AssetHolding, uint64] {
	var choices []weightedrand.Choice[AssetHolding, uint64]

	for collectionName, holdings := range assetsByCollection {
		collection, found := findWeightedCollection(collections, collectionName)
		if !found {
			log.Fatal().Msgf("collection %s not found", collectionName)
		}

		for _, holding := range holdings {
			weightedHolding := weightedrand.Choice[AssetHolding, uint64]{
				Item:   holding,
				Weight: collection.Weight * holding.Amount,
			}
			choices = append(choices, weightedHolding)
		}
	}

	return choices
}

func findWeightedCollection(weightedCollections []WeightedCollection, collectionName string) (WeightedCollection, bool) {
	for _, weightedCollection := range weightedCollections {
		if weightedCollection.Collection.Name == collectionName {
			return weightedCollection, true
		}
	}
	return WeightedCollection{}, false
}
