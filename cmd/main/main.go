package main

import (
	"betprophet1.com/wagers/internal/domains"
	"betprophet1.com/wagers/internal/handlers"
	"betprophet1.com/wagers/internal/repositories"
	"betprophet1.com/wagers/internal/services"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

func main() {
	dsn := "wager:123456@tcp(127.0.0.1:3306)/wager?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
	}
	err = db.AutoMigrate(&domains.Wager{})
	if err != nil {
		fmt.Println(err.Error())
	}

	wagerRepository := repositories.NewWagerRepositoryImpl(*db)
	wagerService    := services.NewWagerServiceImpl(wagerRepository)
	wagerHandler    := handlers.NewWagerHandlerImpl(wagerService)
	r := mux.NewRouter()
	r.HandleFunc("/wagers", wagerHandler.PlaceWager).Methods(http.MethodPost)
	r.HandleFunc("/buy", wagerHandler.BuyWager).Methods(http.MethodPost)
	r.HandleFunc("/wagers", wagerHandler.ListWager).Methods(http.MethodGet)
	srv := &http.Server{
		Addr:              "0.0.0.0:8080",
		Handler:           r,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
	}
	log.Fatalln(srv.ListenAndServe())
}