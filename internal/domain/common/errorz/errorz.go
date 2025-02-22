package errorz

import "errors"

var (
	Forbidden  = errors.New("forbidden")
	NotFound   = errors.New("not found")
	BadRequest = errors.New("bad request")
)
