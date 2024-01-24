package models

import "errors"

var ErrInsertErr = errors.New("cannot insert row")
var ErrNoRows = errors.New("no rows")

var ErrParse = errors.New("cannot parse")
