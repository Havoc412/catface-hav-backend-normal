package errcode

const (
	ErrGeneral = (iota + 1) * 100000
	ErrAnimal
	ErrUser
	ErrEncounter
	ErrNlp
	ErrKnowledge
)

const (
	ErrGeneralStart = ErrGeneral + iota
	ErrInvalidData
	ErrInternalError

	ErrDataNoFound
)
