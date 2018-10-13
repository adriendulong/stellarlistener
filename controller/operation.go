package controller

import (
	"time"

	"github.com/adriendulong/go/stellar/database"
	m "github.com/adriendulong/go/stellar/model"
	log "github.com/sirupsen/logrus"
)

// ListenNewOperations receives new operation got by the API
func ListenNewOperations(channelOperation <-chan m.Operation, r *database.Redis) {
	for op := range channelOperation {
		go op.Save(r)
		log.Info(op.DescribeOperation())
	}
}

// TotalOperationsOfDay returns the total number of operations that has been done
// on the stellar network this particular day
func TotalOperationsOfDay(t time.Time, r *database.Redis) (int, error) {
	resp := r.Client.Get(database.GetCountDayOperationsKey(t))
	if resp.Err() != nil {
		return 0, resp.Err()
	}

	n, err := resp.Int()
	if err != nil {
		return 0, err
	}

	return n, nil
}

// TotalOpearationsOfDayType returns the total number of operations that has been done
// on the stellar network this particular day for this type of operation
func TotalOpearationsOfDayType(t time.Time, typeOpe string, r *database.Redis, opeChannel chan<- int) {
	resp := r.Client.Get(database.GetCountDayOperationsKeyType(t, typeOpe))
	n := 0
	if resp.Err() != nil {
		log.Warn(resp.Err())
	}

	n, err := resp.Int()
	if err != nil {
		log.Warn(resp.Err())
	}

	opeChannel <- n
}

//GetAllBuyinAssetsOfDay returns the assets that have been bought during this day
func GetAllBuyinAssetsOfDay(t time.Time, r *database.Redis, assetsChannel chan<- []string) {
	resp := r.Client.SMembers(database.GetSetBuyingAssetsManageOffer(t))
	if resp.Err() != nil {
		log.Warn(resp.Err())
	}

	buyingAssets, err := resp.Result()
	if err != nil {
		log.Warn(err)
	}

	assetsChannel <- buyingAssets
}
