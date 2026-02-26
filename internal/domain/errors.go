package domain

import "github.com/pkg/errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrInternal        = errors.New("internal")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrAccessDenied    = errors.New("access denied")
	ErrNotAuthorized   = errors.New("not authorized")
	ErrExternalSystem  = errors.New("external system error")
)
