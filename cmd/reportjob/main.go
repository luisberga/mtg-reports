package main

import (
	"context"
	"database/sql"
	"fmt"
	"mtg-report/config/rjobcfg"
	"mtg-report/internal/adapters/email/simplemailtp"
	"mtg-report/internal/adapters/handlers/reporthandler"
	"mtg-report/internal/adapters/repositories/reportrepo"
	"mtg-report/internal/core/services/reportservice"
	"mtg-report/internal/sources/databases/mysql"
	"mtg-report/internal/sources/logger/logrus"
	"mtg-report/internal/sources/timer"
	"net/smtp"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg, err := rjobcfg.New()
	if err != nil {
		panic(err)
	}

	log := logrus.New(cfg.LogLevel)

	auth := smtp.PlainAuth("", cfg.Email.Username, cfg.Email.Password, cfg.Email.Host)

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database))
	if err != nil {
		log.WithError(err).Fatal("failed in db connection")
	}
	defer db.Close()

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.Job.Timeout)
	defer cancelCtx()

	mysql := mysql.New(db)
	timer := timer.New()

	add := cfg.Email.Host + ":" + cfg.Email.Port
	smtp := simplemailtp.New(auth, timer, cfg.Email.Username, cfg.Email.To, add)
	reportRepo := reportrepo.New(mysql)
	reportSrv := reportservice.New(reportRepo, smtp, log)
	reportHand := reporthandler.New(reportSrv, log)

	err = reportHand.ProcessAndSend(ctx)
	if err != nil {
		log.WithError(err).Fatal("failed to process and send")
	}

}
