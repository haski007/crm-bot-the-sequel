package repository

import "errors"

var (
	ErrDocAlreadyExists = errors.New("DOCUMENT_ALREADY_EXISTS")
	ErrDocDoesNotExist  = errors.New("DOCUMENT_DOES_NOT_EXIST")

	ErrYouHaveNoRights = errors.New("У вас недостаточно прав для этой операции.\n" +
		"Обратитесь к вашему начальству")
)

type BotRepository interface {
	InitConn()
}
