package derror

import (
	"fmt"
	"net/http"

	"golang.org/x/xerrors"
)

type ApplicationError struct {
	Err   error
	Code  int
	Msg   string
	Level string
	Frame xerrors.Frame
}

func NewApplicationError(code int, msg string, level string) *ApplicationError {
	return &ApplicationError{
		Err:   nil,
		Code:  code,
		Msg:   msg,
		Level: level,
	}
}

func (e ApplicationError) Error() string {
	return e.Err.Error()
}

func (e ApplicationError) Wrap(originalError error) error {
	// ラップするエラーを渡す
	e.Err = originalError
	return &ApplicationError{
		Err:   e,
		Frame: xerrors.Caller(1),
	}
}

// めっちゃ重要
func (e ApplicationError) Unwrap() error {
	return e.Err
}

// fmt.Formatterを実装
func (e *ApplicationError) Format(s fmt.State, v rune) {
	xerrors.FormatError(e, s, v)
}

// xerrors.Formatterを実装
func (e *ApplicationError) FormatError(p xerrors.Printer) (next error) {
	//p.Print(e.Error())
	e.Frame.Format(p)
	return e.Err
}

func StackError(err error) error {
	return &ApplicationError{
		Err:   err,
		Frame: xerrors.Caller(1),
	}
}

var (
	InternalServerError    = NewApplicationError(http.StatusInternalServerError, "internal server error", "error")
	DatabaseOperationError = NewApplicationError(http.StatusInternalServerError, "database operation error", "error")
	DatabaseDataScanError  = NewApplicationError(http.StatusInternalServerError, "database data scan error", "error")
	GenerateRequestIdError = NewApplicationError(http.StatusInternalServerError, "generate requestID error", "error")
	BadRequestError        = NewApplicationError(http.StatusBadRequest, "bad request error", "warn")
)
