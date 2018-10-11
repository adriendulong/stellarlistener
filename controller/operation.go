package controller

import (
	"fmt"
	"time"

	"github.com/adriendulong/go/stellar/database"
	m "github.com/adriendulong/go/stellar/model"
	"github.com/adriendulong/go/stellar/utils"
)

// ListenNewOperations receives new operation got by the API
func ListenNewOperations(channelOperation <-chan m.Operation, r *database.Redis) {
	for op := range channelOperation {
		go op.Save(r)
		fmt.Println(op.DescribeOperation())
	}
}

// TotalOperationsOfDay returns the total number of operations that has been done
// on the stellar network this particular day
func TotalOperationsOfDay(t time.Time, r *database.Redis) (int, error) {
	resp := r.Pool.Cmd("get", utils.GetCountDayOperationsKey(t))
	if resp.Err != nil {
		return 0, resp.Err
	}

	n, err := resp.Int()
	if err != nil {
		return 0, err
	}

	return n, nil
}
