package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/adriendulong/go/stellar/database"
	log "github.com/sirupsen/logrus"
)

// Operation types
const (
	createAccount      = iota
	payment            = iota
	pathPayment        = iota
	manageOffer        = iota
	createPassiveOffer = iota
	setOptions         = iota
	changeTrust        = iota
	allowTrust         = iota
	accountMerge       = iota
	inflation          = iota
	manageData         = iota
	bumpSequence       = iota
)

// Operation struct contains fields of a stellar operation
type Operation struct {
	ID                 string `json:"id,omitempty"`
	Type               string `json:"type,omitempty"`
	TypeI              int    `json:"type_i,omitempty"`
	StartingBalance    string `json:"starting_balance,omitempty"`
	Funder             string `json:"funder,omitempty"`
	AssetType          string `json:"asset_type,omitempty"`
	AssetCode          string `json:"asset_code,omitempty"`
	AssetIssuer        string `json:"asset_issuer,omitempty"`
	From               string `json:"from,omitempty"`
	To                 string `json:"to,omitempty"`
	Amount             string `json:"amount,omitempty"`
	SellingAssetCode   string `json:"selling_asset_code,omitempty"`
	SellingAssetIssuer string `json:"selling_asset_issuer,omitempty"`
	SellingAssetType   string `json:"selling_asset_type,omitempty"`
	BuyingAssetCode    string `json:"buying_asset_code,omitempty"`
	BuyingAssetIssuer  string `json:"buying_asset_issuer,omitempty"`
	BuyingAssetType    string `json:"buying_asset_type,omitempty"`
	Price              string `json:"price,omitempty"`
	PriceR             PriceR `json:"price_r,omitempty"`
}

// PriceR is the Numerator (Byuying) and Denominator (Selling)
type PriceR struct {
	BuyingPrice  int `json:"n,omitempty"`
	SellingPrice int `json:"d,omitempty"`
}

// DescribeOperation returns a description of the operations
func (o *Operation) DescribeOperation() (s string) {
	resp := "Don't know this Operation"
	switch o.TypeI {
	case createAccount:
		resp = "Account created founded with " + o.StartingBalance + " by account " + o.Funder
	case payment:
		resp = "Payment from " + o.From + " to " + o.To + " of " + o.Amount + "" + o.AssetCode
	case pathPayment:
		resp = "Path Payment: "
	case manageOffer:
		resp = fmt.Sprintf("Manage offer: propose to sell %s %s for some %s at a ratio of %d %s for %d %s (%s)", o.Amount, o.SellingAssetCode, o.BuyingAssetCode, o.PriceR.BuyingPrice, o.BuyingAssetCode, o.PriceR.SellingPrice, o.SellingAssetCode, o.Price)
	case createPassiveOffer:
		resp = "Passive offer: "
	case setOptions:
		resp = "Set Options: "
	case changeTrust:
		resp = "Change Trust: "
	case allowTrust:
		resp = "Allow Trust: "
	case accountMerge:
		resp = "Account Merge: "
	case inflation:
		resp = "Inflation: "
	case manageData:
		resp = "Manage Data: "
	case bumpSequence:
		resp = "Bump Sequence: "
	}

	return resp
}

// Save insert a new operation in the database
func (o *Operation) Save(r *database.Redis) {
	if r.Client == nil {
		panic(errors.New("No Client opened on redis"))
	}

	// Get the date of today
	// It will allow us to build the key
	now := time.Now()

	// Increment the  number of operations of today
	if r.Client.Incr(database.GetCountDayOperationsKey(now)).Err() != nil {
		log.WithFields(log.Fields{
			"operation_id":   o.ID,
			"operation_type": o.Type,
			"date":           now.String(),
		}).Error("Problem incrementig the general number of operation of the day")
	}

	// Increment the number of operations of today depending on the type of the operation
	if r.Client.Incr(database.GetCountDayOperationsKeyType(now, o.Type)).Err() != nil {
		log.WithFields(log.Fields{
			"operation_id":   o.ID,
			"operation_type": o.Type,
			"date":           now.String(),
		}).Error("Problem incrementig the number of operation of a specific type of the day")
	}

	if o.TypeI == manageOffer {
		// Increment the number of operations of today depending on the type of the operation
		buyingAssetCode := o.BuyingAssetCode
		if buyingAssetCode == "" {
			buyingAssetCode = "XLM"
		}

		sellingAssetCode := o.SellingAssetCode
		if sellingAssetCode == "" {
			sellingAssetCode = "XLM"
		}

		//Incr the number total of times an offer has been done for an asset (buying asset)
		if r.Client.Incr(database.GetCountDayManageOfferPerAssetCount(now, buyingAssetCode)).Err() != nil {
			log.WithFields(log.Fields{
				"operation_id":      o.ID,
				"operation_type":    o.Type,
				"buying_asset_code": buyingAssetCode,
				"date":              now.String(),
			}).Error("Problem incrementig the number of manage offer for a certain buying asset")
		}

		//Add this buying asset in a set to list all the buying asset of the day
		if r.Client.SAdd(database.GetSetBuyingAssetsManageOffer(now), buyingAssetCode).Err() != nil {
			log.WithFields(log.Fields{
				"operation_id":      o.ID,
				"operation_type":    o.Type,
				"buying_asset_code": buyingAssetCode,
				"date":              now.String(),
			}).Error("Problem adding a buying asset to the set of buying asset of the day")
		}

		//Incr the number of offer that has been done between this buygin asset and this selling asset
		if r.Client.Incr(database.GetKeyCountBuyingAssetForAnAsset(now, buyingAssetCode, sellingAssetCode)).Err() != nil {
			log.WithFields(log.Fields{
				"operation_id":       o.ID,
				"operation_type":     o.Type,
				"buying_asset_code":  buyingAssetCode,
				"selling_asset_code": sellingAssetCode,
				"date":               now.String(),
			}).Error("Problem incrementig the number of offer that has been done for this asset per this selling asset")
		}

		//Add this selling asset in a set for this buying asset
		if r.Client.SAdd(database.GetKeySetSellingAssetsForBuyingAsset(now, buyingAssetCode), sellingAssetCode).Err() != nil {
			log.WithFields(log.Fields{
				"operation_id":       o.ID,
				"operation_type":     o.Type,
				"buying_asset_code":  buyingAssetCode,
				"selling_asset_code": sellingAssetCode,
				"date":               now.String(),
			}).Error("Problem adding a selling asset code to the set of this buying asset")
		}

		// Add the price proposed in this manage offer between this selling asset and this buying asset
		if r.Client.LPush(database.GetPricesProposedBtwBuyingAssetAndSellingAssetListKey(now, buyingAssetCode, sellingAssetCode), o.Price).Err() != nil {
			log.WithFields(log.Fields{
				"operation_id":       o.ID,
				"operation_type":     o.Type,
				"buying_asset_code":  buyingAssetCode,
				"selling_asset_code": sellingAssetCode,
				"price":              o.Price,
				"date":               now.String(),
			}).Error("Problem adding the price proposed in the list of price")
		}

	}
}
