package main

import (
	"betprophet1.com/wagers/internal/domains"
	"betprophet1.com/wagers/internal/handlers"
	"betprophet1.com/wagers/internal/repositories"
	"betprophet1.com/wagers/internal/services"
	"betprophet1.com/wagers/pkg"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:              time.Second,
			LogLevel:                   logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
	conf := pkg.Get()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.MysqlUser, conf.MysqlPassword, conf.MysqlHost, conf.MysqlPort,conf.MysqlDatabase)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   dbLogger,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	err = db.AutoMigrate(&domains.Wager{}, &domains.Purchase{})
	if err != nil {
		fmt.Println(err.Error())
	}

	wagerRepository := repositories.NewWagerRepository(db)
	purchaseRepository := repositories.NewPurchaseRepository(db)

	wagerService    := services.NewWagerService(wagerRepository)
	purchaseService := services.NewPurchaseService(purchaseRepository, wagerRepository)

	wagerHandler    := handlers.NewWagerHandler(wagerService, purchaseService)

	r := mux.NewRouter()
	r.HandleFunc("/wagers", wagerHandler.PlaceWager).Methods(http.MethodPost)
	r.HandleFunc("/buy/{wager_id}", wagerHandler.BuyWager).Methods(http.MethodPost)
	r.HandleFunc("/wagers", wagerHandler.ListWager).
		Queries("page", "{page:[0-9,]+}", "limit", "{limit:[0-9,]+}").
		Methods(http.MethodGet)
	srv := &http.Server{
		Addr:              "0.0.0.0:8080",
		Handler:           r,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
	}
	log.Fatalln(srv.ListenAndServe())
}