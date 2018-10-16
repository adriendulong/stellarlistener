package router

import (
	"context"
	"fmt"
	"net/http"
	"time"

	c "github.com/adriendulong/go/stellar/controller"
	"github.com/adriendulong/go/stellar/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
)

//OperationRoutes list all the routes for the Operations
func OperationRoutes(r chi.Router) {
	r.Get("/", GeneralOperationHandler)
	r.Route("/{date}", OperationDateRouter)
	// router.Get("/{date}", GeneralOperationHandler)
	// router.Get("/count/{date}", GetOperationCountDay)
	// router.Get("/count/{date}/{type}", GetOperationCountDayType)
	// router.Get("/count/offers/count/{date}/", GetOperationCountDayType)
	// router.Get("/offers/buyingassets/{date}", GetBuyingAssetsOfDay)
}

//OperationDateRouter is a router that handles all the route that contains a date
func OperationDateRouter(r chi.Router) {
	r.Use(DateCtx)
	r.Get("/", GeneralOperationHandler)
	r.Get("/offers", GetGeneralInfosOffers)
}

//DateCtx is a middleware that parse a date and verify its format
func DateCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		date := chi.URLParam(r, "date")
		dateTime, err := time.Parse("02012006", date)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "date", dateTime)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GeneralOperationHandler is the general operation handler that give multuple infos
func GeneralOperationHandler(w http.ResponseWriter, r *http.Request) {

	type GeneralOperationsResp struct {
		TotalCount int `json:"total_count"`
		Types      struct {
			Payment     int `json:"payment"`
			ManageOffer int `json:"manage_offer"`
			Inflation   int `json:"inflation"`
		} `json:"types"`
	}

	resp := GeneralOperationsResp{}

	log.Info("WORKING")

	requestDate := time.Now()
	date, ok := r.Context().Value("date").(time.Time)
	if ok {
		requestDate = date
	}

	redis := GetRedis(r)
	//Get total operations of the day
	nbOperations, err := c.TotalOperationsOfDay(requestDate, redis)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
		return
	}
	resp.TotalCount = nbOperations

	types := make(map[string](chan int))
	types["payment"] = make(chan int)
	types["manage_offer"] = make(chan int)
	types["inflation"] = make(chan int)

	nbOpeTypes := make(map[string]int)
	nbOpeTypes["payment"] = 0
	nbOpeTypes["manage_offer"] = 0
	nbOpeTypes["inflation"] = 0

	for k, v := range types {
		go c.TotalOpearationsOfDayType(requestDate, k, redis, v)
	}

	for i := 0; i < len(nbOpeTypes); i++ {
		select {
		case nbOpeTypes["payment"] = <-types["payment"]:
			resp.Types.Payment = nbOpeTypes["payment"]
		case nbOpeTypes["manage_offer"] = <-types["manage_offer"]:
			resp.Types.ManageOffer = nbOpeTypes["manage_offer"]
		case nbOpeTypes["inflation"] = <-types["inflation"]:
			resp.Types.Inflation = nbOpeTypes["inflation"]
		}
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp)
}

//GetOperationCountDay return the count of operations on a given date
func GetOperationCountDay(w http.ResponseWriter, r *http.Request) {
	date := chi.URLParam(r, "date")
	dateTime, err := time.Parse("02012006", date)
	if err != nil {
		panic(err)
	}

	log.Info(dateTime)
	redis := GetRedis(r)
	nbOperations, err := c.TotalOperationsOfDay(dateTime, redis)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		return
	}

	type CountResponse struct {
		Count int `json:"count"`
	}

	render.JSON(w, r, CountResponse{Count: nbOperations})
}

//GetOperationCountDayType return the count of operations on a given date for a given type
func GetOperationCountDayType(w http.ResponseWriter, r *http.Request) {
	date := chi.URLParam(r, "date")
	dateTime, err := time.Parse("02012006", date)
	if err != nil {
		panic(err)
	}

	typeOpe := chi.URLParam(r, "type")

	redis := GetRedis(r)
	countChannel := make(chan int)
	go c.TotalOpearationsOfDayType(dateTime, typeOpe, redis, countChannel)

	count := <-countChannel

	type CountResponse struct {
		Count int `json:"count"`
	}

	render.JSON(w, r, CountResponse{Count: count})
}

//GetOperationOffersCount return the count of operations on a given date for a given type
func GetOperationOffersCount(w http.ResponseWriter, r *http.Request) {

	render.PlainText(w, r, "Later")
}

//GetBuyingAssetsOfDay return a list of the assets that received an offer today
func GetBuyingAssetsOfDay(w http.ResponseWriter, r *http.Request) {
	date := chi.URLParam(r, "date")
	dateTime, err := time.Parse("02012006", date)
	if err != nil {
		panic(err)
	}

	redis := GetRedis(r)
	buyingAssetsChannel := make(chan []string)
	go c.GetAllBuyinAssetsOfDay(dateTime, redis, buyingAssetsChannel)

	assets := <-buyingAssetsChannel

	// type AssetsResponse struct {
	// 	Asset int `json:"asset"`
	// }

	// assetsResponse := []AssetsResponse{}
	// for asset := range assets {
	// 	assetsResponse.
	// }

	render.JSON(w, r, assets)
}

//GetGeneralInfosOffers is a route to get all the infos about the
//manage offers of the day
func GetGeneralInfosOffers(w http.ResponseWriter, r *http.Request) {
	type BuyingAssetResp struct {
		AssetCode  string `json:"asset_code"`
		TotalCount int    `json:"total_count"`
	}

	type GeneralOffersResp struct {
		TotalCount   int               `json:"total_count"`
		BuyingAssets []BuyingAssetResp `json:"buying_assets"`
	}

	resp := GeneralOffersResp{BuyingAssets: []BuyingAssetResp{}}

	requestDate := time.Now()
	date, ok := r.Context().Value("date").(time.Time)
	if ok {
		requestDate = date
	}

	redis := GetRedis(r)

	c1 := make(chan int)
	c2 := make(chan utils.CountPairList)
	countOffers := 0
	topAssets := utils.CountPairList{}
	go c.TotalOpearationsOfDayType(requestDate, "manage_offer", redis, c1)
	go c.GetTopBuyingAssetsOfDay(requestDate, redis, c2)

	for i := 0; i < 2; i++ {
		select {
		case countOffers = <-c1:
			fmt.Println("Count Offers", countOffers)
		case topAssets = <-c2:
			fmt.Println("Top assets", topAssets)
		}
	}

	resp.TotalCount = countOffers
	for _, a := range topAssets {
		b := BuyingAssetResp{AssetCode: a.Key, TotalCount: a.Count}
		resp.BuyingAssets = append(resp.BuyingAssets, b)
	}

	render.JSON(w, r, resp)

}
