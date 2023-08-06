package exchangegateway

type ErrFailedToGetExchangeRequest struct{}

func (e ErrFailedToGetExchangeRequest) Error() string {
	return "exchange gateway failed to get exchange request"
}

type ErrExchangeRequestMultipleValues struct{}

func (e ErrExchangeRequestMultipleValues) Error() string {
	return "exchange gateway fund multiple values in exchange request"
}

type ErrExchangeRequestNillValue struct{}

func (e ErrExchangeRequestNillValue) Error() string {
	return "exchange gateway fund nill value in exchange request"
}
