package main

import (
	"context"
	"database/sql"
	"fmt"
	"mtg-report/config/cjobcfg"
	"mtg-report/internal/adapters/gateway/cardgateway"
	"mtg-report/internal/adapters/gateway/exchangegateway"
	"mtg-report/internal/adapters/handlers/conciliatehandler"
	"mtg-report/internal/adapters/repositories/conciliaterepo"
	"mtg-report/internal/core/services/conciliateservice"
	"mtg-report/internal/sources/databases/mysql"
	"mtg-report/internal/sources/logger/logrus"
	"mtg-report/internal/sources/web"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg, err := cjobcfg.New()
	if err != nil {
		panic(err)
	}

	log := logrus.New(cfg.LogLevel)

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database))
	if err != nil {
		log.WithError(err).Fatal("failed in db connection")
	}
	defer db.Close()

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.Job.Timeout)
	defer cancelCtx()

	mysql := mysql.New(db)
	http := web.New()

	cardRepo := conciliaterepo.New(mysql)
	cardGateway := cardgateway.New(http, log)
	exchangegateway := exchangegateway.New(http, cfg.ExchangeGateway.Url, log)
	cardSrv := conciliateservice.New(cardRepo, cardGateway, exchangegateway, cfg.Database.CommitSize, log)
	cardHand := conciliatehandler.New(cardSrv, log)

	err = cardHand.Conciliate(ctx)
	if err != nil {
		log.WithError(err).Fatal("failed to conciliate")
	}

}
