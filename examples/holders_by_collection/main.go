package main

import (
	"context"
	"encoding/json"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/rs/zerolog/log"
	"github.com/yellowbackground/holders"
	"github.com/yellowbackground/holders/algorand"
	"github.com/yellowbackground/holders/examples"
	"os"
	"time"
)

const concurrency = 1

var algonodeRefererHeader = []*common.Header{
	{
		Key:   "Referer",
		Value: "https://www.mostlyfrens.xyz/",
	},
}

func main() {
	startTime := time.Now()
	algoD, _ := algod.MakeClientWithHeaders("https://mainnet-api.algonode.cloud", "", algonodeRefererHeader)
	idx, _ := indexer.MakeClientWithHeaders("https://mainnet-idx.algonode.cloud", "", algonodeRefererHeader)
	collectionClient := algorand.NewCollectionClient(algoD, idx)

	log.Info().Msg("Getting holders...")

	holdings, err := holders.GetAssetHoldingsByCollection(context.Background(), collectionClient, examples.Collections, concurrency)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get holders")
	}

	holdingsJSON, err := json.MarshalIndent(holdings, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal holdings to JSON")
		return
	}

	fileName := "holdings.json"
	err = os.WriteFile(fileName, holdingsJSON, 0644)
	if err != nil {
		log.Error().Err(err).Msg("failed to write holdings to file")
		return
	}

	log.Info().Msgf("holdings successfully written to %s", fileName)
	elapsedTime := time.Since(startTime)
	log.Info().Msgf("Elapsed time: %s", elapsedTime)
}
