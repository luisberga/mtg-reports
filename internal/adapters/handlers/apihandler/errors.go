package apihandler

type ErrInternalErr struct{}

func (e ErrInternalErr) Error() string {
	return "internal error"
}
