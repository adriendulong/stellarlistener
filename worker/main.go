package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/adriendulong/go/stellar/api"
	"github.com/adriendulong/go/stellar/controller"
	"github.com/adriendulong/go/stellar/database"
	"github.com/adriendulong/go/stellar/model"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/stellar/go/clients/horizon"
)

var totalOperations int32
var channelTotalOperations int32
var types = map[string]int{
	"manage_offer":   0,
	"payment":        0,
	"create_account": 0,
	"change_trust":   0,
}

var redis *database.Redis

func getOperations(operationsURL string) {
	client := api.NewClient()

	oChan := make(chan model.Operation)
	go controller.ListenNewOperations(oChan, redis)
	go client.GetOperationsFromLedger(operationsURL, oChan)
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Info("ENVIRONMENT IS")
	log.Info(os.Getenv("ENVIRONMENT"))
	if os.Getenv("ENVIRONMENT") == "DEV" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	r := database.New()
	redis = &r

	client := horizon.DefaultPublicNetClient
	cursor := horizon.Cursor("now")

	ctx := context.Background()

	err = client.StreamLedgers(ctx, &cursor, func(l horizon.Ledger) {
		log.WithFields(log.Fields{
			"transactions": l.TransactionCount,
			"operations":   l.OperationCount,
			"closed_time":  l.ClosedAt.String(),
		}).Info("A new ledger")

		totalOperations += l.OperationCount
		go getOperations(strings.Split(l.Links.Operations.Href, "{")[0])
	})

	if err != nil {
		fmt.Println(err)
	}
}
