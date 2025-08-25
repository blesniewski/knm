package main

import (
	"os"

	"github.com/blesniewski/kryptonim/api"
	"github.com/blesniewski/kryptonim/clients/cryptoexchange"
	"github.com/blesniewski/kryptonim/clients/oxr"
)

const (
	oxrBaseUrl = "https://openexchangerates.org/api"
	listenAddr = ":8080"
)

func main() {
	// Should probably use something like godotenv for a real world scenario
	// also for the base url, listening address
	oxrAppId := os.Getenv("OPENEXCHANGERATES_APP_ID")
	if oxrAppId == "" {
		panic("OPENEXCHANGERATES_APP_ID env variable must be set")
	}

	oxrClient := oxr.NewClient(oxrBaseUrl, oxrAppId)
	cryptoClient := cryptoexchange.NewClient()
	httpServer := api.NewServer(oxrClient, cryptoClient)
	httpServer.Run(listenAddr)
}
