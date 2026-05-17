package ports

import "errors"

// ErrNotImplemented is returned by stub adapters until a real implementation lands.
var ErrNotImplemented = errors.New("not implemented")
