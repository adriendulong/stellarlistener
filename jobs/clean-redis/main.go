package main

import (
	"os"
	"sort"
	"time"

	"github.com/adriendulong/go/stellar/database"
	r "github.com/go-redis/redis"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

const (
	limit = 10
)

func main() {
	// Get the date
	dateFormat := os.Args[1]

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	redis := database.New()
	defer redis.Client.Close()

	t, _ := time.Parse("02012006", dateFormat)
	log.Info(t)
	//t := time.Now().Add(time.Hour * -24)

	keyErased := 0

	// Get the list of all the buying assets of the day
	setKey := database.GetSetBuyingAssetsManageOffer(t)
	stringSlice := redis.Client.SMembers(setKey)
	if stringSlice.Err() != nil {
		log.Error(stringSlice.Err())
		panic(stringSlice.Err())
	}

	buyingAssetsList, err := stringSlice.Result()
	if err != nil {
		panic(err)
	}

	if len(buyingAssetsList) < 1 {
		log.Info("Nothing to do")
		return
	}

	// Get thhe count of all these buying assets
	pipe := redis.Client.Pipeline()
	countsResult := make(map[string]*r.StringCmd)
	counts := make(map[string]int)
	for _, buyingAsset := range buyingAssetsList {
		key := database.GetCountDayManageOfferPerAssetCount(t, buyingAsset)
		result := pipe.Get(key)
		countsResult[buyingAsset] = result
	}
	_, err = pipe.Exec()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Problem getting the count of the buying assets")
	}
	for k, result := range countsResult {
		counts[k], err = result.Int()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Problem getting the count of this buying asset")
		}
	}

	log.Info(counts)

	// Sort al these assets count
	assetListSorted := rankByAssetsCount(counts)
	log.Info(assetListSorted)

	//Erase all the keys that are linked to these assets
	i := 1
	nbToErase := len(assetListSorted) - limit
	doneChannel := make(chan int, nbToErase)
	log.Infof("Nb to erase : %d", nbToErase)
	logCount := 0
	for _, countPairAsset := range assetListSorted {
		if i > limit {
			logCount++
			log.Infof("Erase everything for asset %s", countPairAsset.Key)
			go eraseAllForThisAsset(t, countPairAsset.Key, &redis, doneChannel)
		}
		i++
	}
	log.Infof("Go routines launched: %d", logCount)

	for j := 1; j <= nbToErase; j++ {
		c := <-doneChannel
		keyErased += c
		log.Info("Finish one erase")
	}
	close(doneChannel)
	log.Infof("Finished ALL erase. Nb of erased keys: %d", keyErased)

}

// eraseAllForThisAsset erase all the keys that concern this buying asset
func eraseAllForThisAsset(t time.Time, buyingAsset string, r *database.Redis, doneChannel chan<- int) {

	//Erase the count of this asset
	// key := database.GetCountDayManageOfferPerAssetCount(t, buyingAsset)
	// result := r.Client.Del(key)
	// if result.Err() != nil {
	// 	log.WithFields(log.Fields{
	// 		"buying_asset": buyingAsset,
	// 	}).Error("Problem trying to del the count of this asset")
	// }

	count := 0

	// Erase the list of selling asset for this buying asset
	sellingAssetsListKey := database.GetKeySetSellingAssetsForBuyingAsset(t, buyingAsset)
	stringSlice := r.Client.SMembers(sellingAssetsListKey)
	if stringSlice.Err() != nil {
		log.WithFields(log.Fields{
			"buying_asset": buyingAsset,
		}).Error("Problem trying to get the list of selling assets for this buying asset")
	}

	sellingAssetsList, err := stringSlice.Result()
	if err != nil {
		log.WithFields(log.Fields{
			"buying_asset": buyingAsset,
		}).Error("Problem trying to get the list of selling assets for this buying asset")
	}

	// Iterate over selling assets to remove the keys
	for _, asset := range sellingAssetsList {
		countKey := database.GetKeyCountBuyingAssetForAnAsset(t, buyingAsset, asset)
		priceKey := database.GetPricesProposedBtwBuyingAssetAndSellingAssetListKey(t, buyingAsset, asset)
		pipe := r.Client.Pipeline()
		pipe.Del(countKey)
		pipe.Del(priceKey)

		_, err := pipe.Exec()

		if err != nil {
			log.WithFields(log.Fields{
				"buying_asset":  buyingAsset,
				"selling_asset": asset,
			}).Error("Problem trying to remove the count and price of this buying asset and selling asset")
		} else {
			count += 2
		}
	}

	doneChannel <- count

}

type countPair struct {
	Key   string
	Count int
}
type countPairList []countPair

func (c countPairList) Len() int           { return len(c) }
func (c countPairList) Less(i, j int) bool { return c[i].Count < c[j].Count }
func (c countPairList) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

func rankByAssetsCount(assetsCount map[string]int) countPairList {
	cl := make(countPairList, len(assetsCount))
	i := 0
	for k, v := range assetsCount {
		cl[i] = countPair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(cl))
	return cl
}
