package main

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/rs/zerolog/log"
	"github.com/yellowbackground/holders"
	"github.com/yellowbackground/holders/algorand"
)

var collections = []holders.WeightedCollection{
	{
		Weight: 6,
		Collection: holders.Collection{
			Name: "Yieldlings Flambos",
			Addresses: []string{
				"5DYIZMX7N4SAB44HLVRUGLYBPSN4UMPDZVTX7V73AIRMJQA3LKTENTLFZ4",
			},
			IncludeNameContains: []string{"Flamborghini"},
		},
	},
	{
		Weight: 4,
		Collection: holders.Collection{
			Name: "Yieldlings",
			Addresses: []string{
				"5DYIZMX7N4SAB44HLVRUGLYBPSN4UMPDZVTX7V73AIRMJQA3LKTENTLFZ4",
			},
			ExcludeNameContains: []string{"Flamborghini"},
			UnitNamePrefixes:    []string{"TLDG", "YLD"},
		},
	},
	{
		Weight: 2,
		Collection: holders.Collection{
			Name: "M.N.G.O",
			Addresses: []string{
				"MNGOLDXO723TDRM6527G7OZ2N7JLNGCIH6U2R4MOCPPLONE3ZATOBN7OQM",
				"MNGORTG4A3SLQXVRICQXOSGQ7CPXUPMHZT3FJZBIZHRYAQCYMEW6VORBIA",
				"MNGOZ3JAS3C4QTGDQ5NVABUEZIIF4GAZY52L3EZE7BQIBFTZCNLQPXHRHE",
				"MNGO4JTLBN64PJLWTQZYHDMF2UBHGJGW5L7TXDVTJV7JGVD5AE4Y3HTEZM",
			},
			UnitNamePrefixes: []string{"MNGO"},
		},
	},
	{
		Weight: 1,
		Collection: holders.Collection{
			Name: "Mostly Frens",
			Addresses: []string{
				"MOSTLYSNUJP7PG6Q3FNJCGGENQXMOH3PXXMIJRFLODLG2DNDBHI7QHJSOE",
			},
			UnitNamePrefixes: []string{"MFER"},
		},
	},
	{
		Weight: 6,
		Collection: holders.Collection{
			Name: "Best Frens",
			Addresses: []string{
				"MOSTLYSNUJP7PG6Q3FNJCGGENQXMOH3PXXMIJRFLODLG2DNDBHI7QHJSOE",
			},
			UnitNamePrefixes: []string{"BFER"},
		},
	},
}

const (
	concurrency     = 1
	numberOfWinners = 3
)

var algonodeRefererHeader = []*common.Header{
	{
		Key:   "Referer",
		Value: "https://www.mostlyfrens.xyz/",
	},
}

func main() {
	var excludedWallets = []string{}
	for _, collection := range collections {
		excludedWallets = append(excludedWallets, collection.Collection.Addresses...)
	}

	algoD, _ := algod.MakeClientWithHeaders("https://mainnet-api.algonode.cloud", "", algonodeRefererHeader)
	idx, _ := indexer.MakeClientWithHeaders("https://mainnet-idx.algonode.cloud", "", algonodeRefererHeader)
	collectionClient := algorand.NewCollectionClient(algoD, idx)

	log.Info().Msg("Running raffle...")

	winningAssets, err := holders.RunWeightedCollectionRaffle(context.Background(), collectionClient, collections, numberOfWinners, concurrency, excludedWallets)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get holders")
	}

	log.Info().Msgf("Winning asset: %v", winningAssets)
}
