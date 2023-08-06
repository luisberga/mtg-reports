package main

import (
	"context"
	"database/sql"
	"fmt"
	"mtg-report/config/apicfg"
	"mtg-report/internal/adapters/handlers/apihandler"
	"mtg-report/internal/adapters/repositories/cardrepo"
	"mtg-report/internal/core/services/cardservice"
	"mtg-report/internal/core/validate"
	"mtg-report/internal/sources/databases/mysql"
	"mtg-report/internal/sources/logger/logrus"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg, err := apicfg.New()
	if err != nil {
		panic(err)
	}

	log := logrus.New(cfg.LogLevel)

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database))
	if err != nil {
		log.WithError(err).Fatal("failed in db connection")
	}
	defer db.Close()

	mysql := mysql.New(db)
	requestVal := validate.New()

	cardRepo := cardrepo.New(mysql, log)
	cardSrv := cardservice.New(cardRepo, cfg.Database.CommitSize, log)
	cardHand := apihandler.New(requestVal, cardSrv, log)

	router := apihandler.SetupRouter(cardHand)

	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		err := http.ListenAndServe(cfg.Api.Port, router)
		if err != nil {
			log.WithError(err).Fatal("server initialization error")
			cancelCtx()
		}
	}()

	log.WithFields(logrus.Fields{
		"port": cfg.Api.Port,
	}).Info("server initialized")

	<-ctx.Done()
}
