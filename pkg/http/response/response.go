package response

import (
	"InfecShotAPI/pkg/derror"
	"InfecShotAPI/pkg/logging"
	"encoding/json"
	"errors"
	"net/http"
)

// HttpResponse レスポンス出力のための構造体
type HttpResponse struct{}

// HttpResponseInterface レスポンス出力のためのインターフェース
type HttpResponseInterface interface {
	Success(writer http.ResponseWriter, request *http.Request, response interface{})
	Failed(writer http.ResponseWriter, request *http.Request, err error)
}

// NewHttpResponse レスポンス出力のための構造体の初期化をする
func NewHttpResponse() *HttpResponse {
	return &HttpResponse{}
}

// Success HTTPコード:200 正常終了を処理する
func (hr *HttpResponse) Success(writer http.ResponseWriter, request *http.Request, response interface{}) {
	if response == nil {
		return
	}
	data, err := json.Marshal(response)
	if err != nil {
		// TODO:アプリケーションログ
		//log.Println(err)
		hr.Failed(writer, request, derror.InternalServerError.Wrap(err))
		return
	}
	writer.Write(data)
	logging.AccessLogging(request, err)
}

// Failed リクエスト失敗時のエラー処理
func (hr *HttpResponse) Failed(writer http.ResponseWriter, request *http.Request, err error) {
	var appErr derror.ApplicationError
	if errors.As(err, &appErr) {
		switch appErr.Code {
		case http.StatusBadRequest:
			hr.badRequest(writer, appErr.Msg)
		case http.StatusInternalServerError:
			hr.internalServerError(writer, appErr.Msg)
		}
	} else {
		hr.internalServerError(writer, "Unknown Internal Server Error")
	}
	logging.AccessLogging(request, err)
}

// BadRequest HTTPコード:400 BadRequestを処理する
func (hr *HttpResponse) badRequest(writer http.ResponseWriter, message string) {
	HttpError(writer, http.StatusBadRequest, message)
}

// InternalServerError HTTPコード:500 InternalServerErrorを処理する
func (hr *HttpResponse) internalServerError(writer http.ResponseWriter, message string) {
	HttpError(writer, http.StatusInternalServerError, message)
}

// httpError エラー用のレスポンス出力を行う
func HttpError(writer http.ResponseWriter, code int, message string) {
	data, _ := json.Marshal(errorResponse{
		Code:    code,
		Message: message,
	})
	writer.WriteHeader(code)
	if data != nil {
		writer.Write(data)
	}
}

// errorResponse エラー時の構造体
type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
