package caql

import "errors"

var (
	ErrStack     = errors.New("unexpected operator stack")
	ErrUndefined = errors.New("variable not defined")
)
