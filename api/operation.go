package api

import (
	"encoding/json"
	"fmt"
	"log"
	"sync/atomic"

	m "github.com/adriendulong/go/stellar/model"
)

// Operations is a list of operation
type Operations struct {
	Links    Links `json:"_links"`
	Embedded struct {
		Records []m.Operation `json:"records"`
	} `json:"_embedded"`
}

// GetOperationsFromLedger get all the operations that has been save in a ledger
func (c *Client) GetOperationsFromLedger(operationsURL string, channelOperation chan<- m.Operation) {
	resp, err := c.client.Get(operationsURL)
	if err != nil {
		log.Fatalln(err)
	}

	ops := new(Operations)
	if err := json.NewDecoder(resp.Body).Decode(ops); err != nil {
		log.Fatalln(err)
	}

	for _, o := range ops.Embedded.Records {
		channelOperation <- o
	}

	if ops.Links.Next.Href != "" && len(ops.Embedded.Records) > 0 {
		c.GetOperationsFromLedger(ops.Links.Next.Href, channelOperation)
	} else {
		close(channelOperation)
	}
}

// ReadOperation list a channel of operation and displays it
func ReadOperation(channelOperation <-chan m.Operation, totalOperation *int32) {
	for op := range channelOperation {
		fmt.Println(op.DescribeOperation())
		atomic.AddInt32(totalOperation, 1)
	}
}
