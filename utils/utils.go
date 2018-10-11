package utils

import (
	"fmt"
	"time"
)

// GetCountDayOperationsKey return a string that is the key
// that must be used in order to get the total number of operations
// on that day
func GetCountDayOperationsKey(t time.Time) (s string) {
	s = fmt.Sprintf("operations:count:%d%d%d", t.Day(), t.Month(), t.Year())
	return
}
