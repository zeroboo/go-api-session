package apisession

import "fmt"

var ErrTooFast = fmt.Errorf("request too fast")
var ErrTooMany = fmt.Errorf("too many requests")
var ErrInvalidSession = fmt.Errorf("invalid session")
