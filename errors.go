package switchboard

import (
	"errors"
)

var ErrInvalidArgumentType = errors.New("struct field type is not a supported option type")
