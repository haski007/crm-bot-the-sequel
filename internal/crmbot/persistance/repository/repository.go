package repository

import "errors"

var (
	ErrDocAlreadyExists = errors.New("DOCUMENT_ALREADY_EXISTS")
	ErrDocDoesNotExist  = errors.New("DOCUMENT_DOES_NOT_EXIST")
)

type BotRepository interface {
	InitConn()
}
