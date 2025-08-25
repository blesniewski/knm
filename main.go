package main

import (
	"os"

	"github.com/blesniewski/kryptonim/api"
	"github.com/blesniewski/kryptonim/clients/cryptoexchange"
	"github.com/blesniewski/kryptonim/clients/oxr"
)

const (
	oxrBaseUrl = "https://openexchangerates.org/api"
)

func main() {
	oxrAppId := os.Getenv("OPENEXCHANGERATES_APP_ID")
	if oxrAppId == "" {
		panic("OPENEXCHANGERATES_APP_ID env variable must be set")
	}

	oxrClient := oxr.NewClient(oxrBaseUrl, oxrAppId)
	cryptoClient := cryptoexchange.NewClient()
	httpServer := api.NewServer(oxrClient, cryptoClient)
	httpServer.Run(":8080")
}
