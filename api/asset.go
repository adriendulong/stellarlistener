package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/stellar/go/protocols/horizon"
)

// Assets is the api response for asset endpoint
type Assets struct {
	Links    Links `json:"_links"`
	Embedded struct {
		Records []horizon.AssetStat `json:"records"`
	} `json:"_embedded"`
}

// LoadAllAssets Will load all the assets of stellar
func (c *Client) LoadAllAssets() ([]horizon.AssetStat, error) {
	assetURLString := c.BaseURL.String() + "/assets"
	assetURL, err := url.Parse(assetURLString)
	if err != nil {
		panic(err)
	}

	//Add how many asset we want at once
	v := url.Values{}
	v.Set("limit", "100")
	assetURL.RawQuery = v.Encode()

	fmt.Println(assetURL.String())
	resp, err := c.client.Get(assetURL.String())
	if err != nil {
		return nil, err
	}

	ops := new(Assets)
	if err := json.NewDecoder(resp.Body).Decode(ops); err != nil {
		return nil, err
	}

	return ops.Embedded.Records, nil
}
