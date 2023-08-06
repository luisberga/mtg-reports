package ports

type Email interface {
	SendEmail(cardsTable, cardsPriceTable string) error
}
