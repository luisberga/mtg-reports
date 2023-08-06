package simplemailtp

import (
	"fmt"
	"mtg-report/internal/sources/timer"
	"net/smtp"
)

type email struct {
	auth   smtp.Auth
	timer  timer.Timer
	from   string
	to     string
	adress string
}

func New(auth smtp.Auth, timer timer.Timer, from string, to string, adress string) *email {
	return &email{
		auth:   auth,
		timer:  timer,
		from:   from,
		to:     to,
		adress: adress,
	}
}

func (e *email) SendEmail(cardsTable, cardsPrice string) error {
	to := []string{e.to}

	timestamp := e.timer.Now()
	subject := "Subject: Daily MTG Investment Report\r\n"
	mime := "MIME-Version: 1.0\r\n"
	contentType := "Content-Type: text/html; charset=UTF-8\r\n"
	htmlOpening := "<html><head><style>body { font-family: Arial, sans-serif; }</style></head><body>"
	title := "<h1>Daily MTG Investment Report</h1>"
	date := "<p><strong>Report Date: </strong>" + timestamp + "</p>"
	cardPrice := "<p>" + cardsPrice + "</p>"
	htmlClosing := "</body></html>\r\n"

	msg := []byte(subject + mime + contentType + "\r\n" + htmlOpening + title + date + cardPrice + cardsTable + htmlClosing)

	err := smtp.SendMail(e.adress, e.auth, e.from, to, msg)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
