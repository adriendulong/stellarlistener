package database

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

//Redis is a struct that hold the pool connection
type Redis struct {
	Client *redis.Client
}

//RedisKey represents all the redis keys
type RedisKey uint32

// Convert the Level to a string. E.g. PanicLevel becomes "panic".
func (key RedisKey) String() string {
	switch key {
	case OperationCountDate:
		return "%d%d%d:operations:count"
	case OperationTypeCountDate:
		return "%d%d%d:operations:%s:count"
	case OperationManageOfferPerAssetCount:
		return "%d%d%d:operations:manage_offer:buying:%s:count"
	case OperationManageOfferBuyingAssets:
		return "%d%d%d:operations:manage_offer:buyingassets"
	case OperationManageOfferBuyingAssetsCount:
		return "%d%d%d:operations:manage_offer:%s:%s:count"
	case OperationManageOfferSellingAssetForBuyginAssetSetKey:
		return "%d%d%d:operations:manage_offer:%s"
	case PricesProposedBtwBuyingAssetAndSellingAssetListKey:
		return "%d%d%d:operations:manage_offer:%s:%s:prices"
	}

	return "unknown"
}

// These are the different redis keys
const (
	// OperationCountDate represents the count of total operations per date
	OperationCountDate RedisKey = iota

	// OperationTypeCountDate represents the count of the operation of a certain type at this date
	OperationTypeCountDate

	// OperationManageOfferPerAssetCount is the key for an incr of the number of manage offer that
	// has been done for the particular asset (as a buying asset)
	OperationManageOfferPerAssetCount

	// OperationManageOfferBuyingAssets represents a set that list all buying assets
	// that were at leat in one manage offer on this date
	OperationManageOfferBuyingAssets

	// OperationManageOfferBuyingAssetsCount represents a incr that count the number of manage offer
	// that has been done between this buying asset and this selling asset on this date
	OperationManageOfferBuyingAssetsCount

	// OperationManageOfferSellingAssetForBuyginAssetSetKey is a set that list all the selling
	// asset that were proposed for a particular buying asset on this day
	OperationManageOfferSellingAssetForBuyginAssetSetKey

	// PricesProposedBtwBuyingAssetAndSellingAssetListKey is a list with all the prices that
	// has been found in a manage offer between a particuler selling asset and a buying asset
	PricesProposedBtwBuyingAssetAndSellingAssetListKey
)

// New open the pool connection and return
func New() (r Redis) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}
	opt.PoolSize = 5
	client := redis.NewClient(opt)
	r = Redis{Client: client}
	return
}

// SetKey simple function that set a key with a value
func (r *Redis) SetKey(k string, v interface{}) {
	if r.Client == nil {
		panic(errors.New("No CLient found"))
	}

	resp := r.Client.Set(k, v, 0)
	if resp.Err() != nil {
		panic(resp.Err())
	}
}

// GetCountDayOperationsKey return a string that is the key
// that must be used in order to get the total number of operations
// on that day
func GetCountDayOperationsKey(t time.Time) (s string) {
	s = fmt.Sprintf(OperationCountDate.String(), t.Day(), t.Month(), t.Year())
	return
}

// GetCountDayOperationsKeyType return a string that is the key
// that must be used in order to get number of operations for
// a specific type
func GetCountDayOperationsKeyType(t time.Time, typeOpe string) (s string) {
	s = fmt.Sprintf(OperationTypeCountDate.String(), t.Day(), t.Month(), t.Year(), typeOpe)
	return
}

// GetCountDayManageOfferPerAssetCount return a string that is the key
// that must be used in order to get number of manage offer for a specific asset
func GetCountDayManageOfferPerAssetCount(t time.Time, assetCode string) (s string) {
	s = fmt.Sprintf(OperationManageOfferPerAssetCount.String(), t.Day(), t.Month(), t.Year(), assetCode)
	return
}

// GetSetBuyingAssetsManageOffer return a list of buying asset for the manage offers of the day
func GetSetBuyingAssetsManageOffer(t time.Time) (s string) {
	s = fmt.Sprintf(OperationManageOfferBuyingAssets.String(), t.Day(), t.Month(), t.Year())
	return
}

// GetKeyCountBuyingAssetForAnAsset return a list of buying asset for the manage offers of the day
func GetKeyCountBuyingAssetForAnAsset(t time.Time, buyingAsset string, sellingAsset string) (s string) {
	s = fmt.Sprintf(OperationManageOfferBuyingAssetsCount.String(), t.Day(), t.Month(), t.Year(), buyingAsset, sellingAsset)
	return
}

// GetKeySetSellingAssetsForBuyingAsset returns the key of a set that lists all the selling assets of offers
// that has been made for this buying asset on a specific day
func GetKeySetSellingAssetsForBuyingAsset(t time.Time, buyingAsset string) (s string) {
	s = fmt.Sprintf(OperationManageOfferSellingAssetForBuyginAssetSetKey.String(), t.Day(), t.Month(), t.Year(), buyingAsset)
	return
}

// GetPricesProposedBtwBuyingAssetAndSellingAssetListKey returns the PricesProposedBtwBuyingAssetAndSellingAssetListKey key
// formattes with the date, the buying asset code and the selling asset code
func GetPricesProposedBtwBuyingAssetAndSellingAssetListKey(t time.Time, buyingAsset string, sellingAsset string) (s string) {
	s = fmt.Sprintf(PricesProposedBtwBuyingAssetAndSellingAssetListKey.String(), t.Day(), t.Month(), t.Year(), buyingAsset, sellingAsset)
	return
}
