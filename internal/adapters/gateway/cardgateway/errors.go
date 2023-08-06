package cardgateway

import "fmt"

type ErrFailedToGetCardRequest struct {
	httpStatus int
}

func (e ErrFailedToGetCardRequest) Error() string {
	return fmt.Sprintf("failed to get card request: %d", e.httpStatus)
}

type ErrPriceIsZero struct{}

func (e ErrPriceIsZero) Error() string {
	return "price is zero - card could be foil or non-foil, check register"
}

type ErrCardNotFound struct{}

func (e ErrCardNotFound) Error() string {
	return "card not found"
}
