package api

import (
	"net/http"
	"net/url"

	"github.com/stellar/go/support/render/hal"
)

const (
	testNetClientURL   = "https://horizon-testnet.stellar.org"
	publicNetClientURL = "https://horizon.stellar.org"
)

// Client struct that contains data necessary to communicate with horizon server
type Client struct {
	client *http.Client

	BaseURL *url.URL
}

// Links is the element we get back from each request on Horizon server
type Links struct {
	Self hal.Link `json:"self"`
	Next hal.Link `json:"next"`
	Prev hal.Link `json:"prev"`
}

// NewTestClient return a Horizon API Client pointing to the test core stellar network
func NewTestClient() *Client {
	baseURL, _ := url.Parse(testNetClientURL)
	c := &Client{client: http.DefaultClient, BaseURL: baseURL}
	return c
}

// NewClient return a Horizon API Client pointing to the real core stellar network
func NewClient() *Client {
	baseURL, _ := url.Parse(publicNetClientURL)
	c := &Client{client: http.DefaultClient, BaseURL: baseURL}
	return c
}
