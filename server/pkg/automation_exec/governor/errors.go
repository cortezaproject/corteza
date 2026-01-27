package governor

import "errors"

var (
	ErrInvalidOps             = errors.New("governor: ops must be positive non zero value")
	ErrExecutionExists        = errors.New("governor: execution already initialized")
	ErrStepTooExpensive       = errors.New("governor: step ops exceed configured limits")
	ErrExecutableTooExpensive = errors.New("governor: an executable step exceed configured limits")
)
