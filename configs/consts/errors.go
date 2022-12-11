package consts

import "errors"

const (
	// DuplicationCode aka UniqueConstraintCode
	DuplicationCode = "23505"
)

var (
	LongLanguageError   = errors.New("language cannot be longer than 2 symbols")
	UnknownCommandError = errors.New("unknown command")
	RootDomainError     = errors.New("root domain undefined")
	AdaptationError     = errors.New("no adaptation for this url")
)
