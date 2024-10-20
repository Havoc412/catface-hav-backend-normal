package errcode

const (
	ErrGeneral = (iota + 1) * 100000
	ErrAnimal
	ErrUser
)

const (
	ErrGeneralStart = ErrGeneral + iota
	ErrInvalidData
	ErrInternalError

	ErrDataNoFound
)
