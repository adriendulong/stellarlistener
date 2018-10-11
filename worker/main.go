package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/adriendulong/go/stellar/api"
	"github.com/adriendulong/go/stellar/controller"
	"github.com/adriendulong/go/stellar/database"
	"github.com/adriendulong/go/stellar/model"
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

	r := database.New()
	redis = &r

	client := horizon.DefaultPublicNetClient
	cursor := horizon.Cursor("now")

	ctx := context.Background()

	err := client.StreamLedgers(ctx, &cursor, func(l horizon.Ledger) {
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
