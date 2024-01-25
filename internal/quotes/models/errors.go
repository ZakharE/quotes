package models

import "errors"

var ErrInsertErr = errors.New("cannot insert row")
var ErrNoRows = errors.New("no rows")
var ErrTransaction = errors.New("problem during transaction")
var ErrTransactionRollback = errors.New("transaction was rolledback")

var ErrParse = errors.New("cannot parse")
