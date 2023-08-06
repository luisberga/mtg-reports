package domain

type ErrCardAlreadyExists struct{}

func (e ErrCardAlreadyExists) Error() string {
	return "card already exists"
}

type ErrCardNotFound struct{}

func (e ErrCardNotFound) Error() string {
	return "card not found"
}

type ErrCardsPriceNotFound struct{}

func (e ErrCardsPriceNotFound) Error() string {
	return "cards price not found"
}

type ErrInvalidSetName struct{}

func (e ErrInvalidSetName) Error() string {
	return "invalid set name"
}
