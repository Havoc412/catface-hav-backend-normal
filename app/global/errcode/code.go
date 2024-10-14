package errcode

const (
	ErrGeneral = (iota + 1) * 100000
	ErrAnimal
)

const (
	ErrGeneralStart = ErrGeneral + iota
	ErrInvalidData
	ErrInternalError
)
