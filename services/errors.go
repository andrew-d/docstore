package services

var (
	ErrDocumentNotFound = NotFoundError{"document"}
	ErrFileNotFound     = NotFoundError{"file"}
	ErrTagNotFound      = NotFoundError{"tag"}

	ErrFileExists = ExistsError{"file"}
)

func IsExists(err error) bool {
	_, ok := err.(ExistsError)
	return ok
}

func IsNotFound(err error) bool {
	_, ok := err.(NotFoundError)
	return ok
}

type ExistsError struct {
	Type string
}

func (e ExistsError) Error() string {
	return e.Type + " exists"
}

type NotFoundError struct {
	Type string
}

func (e NotFoundError) Error() string {
	return e.Type + " not found"
}

var _ error = ExistsError{}
var _ error = NotFoundError{}
