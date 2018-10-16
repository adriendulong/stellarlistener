package utils

import "sort"

//CountPair is an helper struct that allows to sort a map of assets counter
type CountPair struct {
	Key   string
	Count int
}

//CountPairList is an helper struct that allows to sort a map of assets counter
type CountPairList []CountPair

func (c CountPairList) Len() int           { return len(c) }
func (c CountPairList) Less(i, j int) bool { return c[i].Count < c[j].Count }
func (c CountPairList) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

//RankAssetsByCount sort a map of assets counter
func RankAssetsByCount(assetsCount map[string]int) CountPairList {
	cl := make(CountPairList, len(assetsCount))
	i := 0
	for k, v := range assetsCount {
		cl[i] = CountPair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(cl))
	return cl
}
