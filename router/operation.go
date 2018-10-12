package router

import (
	"fmt"
	"net/http"
	"time"

	c "github.com/adriendulong/go/stellar/controller"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
)

//OperationRoutes list all the routes for the Operations
func OperationRoutes() http.Handler {
	router := chi.NewRouter()
	router.Get("/", GeneralOperationHandler)
	router.Get("/count/{date}", GetOperationCountDay)
	router.Get("/count/{date}/{type}", GetOperationCountDayType)

	return router
}

// GeneralOperationHandler is the general operation handler that give multuple infos
func GeneralOperationHandler(w http.ResponseWriter, r *http.Request) {

	redis := GetRedis(r)
	//Get total operations of the day
	nbOperations, err := c.TotalOperationsOfDay(time.Now(), redis)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
		return
	}

	types := make(map[string](chan int))
	types["payment"] = make(chan int)
	types["manage_offer"] = make(chan int)
	types["inflation"] = make(chan int)

	nbOpeTypes := make(map[string]int)
	nbOpeTypes["payment"] = 0
	nbOpeTypes["manage_offer"] = 0
	nbOpeTypes["inflation"] = 0

	for k, v := range types {
		go c.TotalOpearationsOfDayType(time.Now(), k, redis, v)
	}

	for i := 0; i < len(nbOpeTypes); i++ {
		select {
		case nbOpeTypes["payment"] = <-types["payment"]:
			log.Info("Got payment nb")
		case nbOpeTypes["manage_offer"] = <-types["manage_offer"]:
			log.Info("Got manage_offer nb")
		case nbOpeTypes["inflation"] = <-types["inflation"]:
			log.Info("Got inflation nb")
		}
	}

	w.WriteHeader(http.StatusOK)
	sResp := fmt.Sprintf("Number total of operation on %s %dth is %d\nNumber of Payment is %d\nNumber of Manage Offer is %d\nNumber of Inlfation is %d", time.Now().Month().String(), time.Now().Day(), nbOperations, nbOpeTypes["payment"], nbOpeTypes["manage_offer"], nbOpeTypes["inflation"])
	render.PlainText(w, r, sResp)
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
