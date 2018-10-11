package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/adriendulong/go/stellar/api"
	"github.com/stellar/go/clients/horizon"
)

var totalOperations = 0
var channelTotalOperations = 0
var types = map[string]int{
	"manage_offer":   0,
	"payment":        0,
	"create_account": 0,
	"change_trust":   0,
}

// Fetch all the operations of a operations link
// Display the type of the operation
func getOperationsType(operationsURL string) {
	// client := api.NewClient()

	// oChan := make(chan api.Operation)
	// go api.ReadOperation(oChan, &channelTotalOperations)
	// go client.GetOperationsFromLedger(operationsURL, oChan)
	// ops, err := client.GetOperationsFromLedger(legderID)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// for _, o := range ops {
	// 	totalOperations++
	// 	_, ok := types[o.Type]
	// 	if ok {
	// 		types[o.Type]++
	// 	}
	// 	fmt.Println(o.DescribeOperation())
	// }

	// if operations.Links.Next.Href != "" && len(operations.Embedded.Records) != 0 {
	// 	go getOperationsType(operations.Links.Next)
	// }
}

// StreamOperations log all the operations that happens on stellar during one minute
func StreamOperations() {
	client := horizon.DefaultPublicNetClient
	cursor := horizon.Cursor("now")

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		// Stop streaming after 60 seconds.
		time.Sleep(30 * time.Second)
		log.Printf("TOTAL OPE : %d", totalOperations)
		log.Printf("TOTAL CHANNEL OPE: %d", channelTotalOperations)
		cancel()
	}()

	err := client.StreamLedgers(ctx, &cursor, func(l horizon.Ledger) {
		fmt.Printf("Number of Transaction: %d\n", l.TransactionCount)
		fmt.Printf("Number of OPerations: %d\n", l.OperationCount)
		totalOperations += int(l.OperationCount)
		fmt.Printf("Closed at: %s\n", l.ClosedAt.String())
		go getOperationsType(strings.Split(l.Links.Operations.Href, "{")[0])
	})

	if err != nil {
		fmt.Println(err)
	}
}

// LoadAllAssets take car of loading assets from the stellar network
func LoadAllAssets() {
	client := api.NewClient()
	assets, err := client.LoadAllAssets()
	if err != nil {
		log.Fatalln(err)
	}

	counter := 0
	for _, a := range assets {
		counter++
		fmt.Printf("Asset code : %s, Num accounts: %d, Amount: %s", a.Asset.Code, a.NumAccounts, a.Amount)
		fmt.Println(a.Asset.Code)
		//fmt.Println(o.DescribeOperation())
	}

	fmt.Println(counter)
}

func main() {

	fmt.Println("What do you want to do?\n1) Load all assets 2) Stream operations during one minute")
	var a int
	fmt.Scanln(&a)

	switch a {
	case 1:
		LoadAllAssets()
	case 2:
		StreamOperations()
	default:
		fmt.Println("Don't know this command")
	}
}
