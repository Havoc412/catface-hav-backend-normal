package errcode

const (
	ErrGeneral = (iota + 1) * 100000
	ErrAnimal
	ErrUser
	ErrEncounter
	ErrNlp
	ErrKnowledge
	ErrSubService
)

const (
	ErrGeneralStart = ErrGeneral + iota
	ErrInvalidData
	ErrInternalError

	ErrDataNoFound
)
