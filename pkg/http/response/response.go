package response

import (
	"InfecShotAPI/pkg/derror"
	"InfecShotAPI/pkg/logging"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// HttpResponse レスポンス出力のための構造体
type HttpResponse struct{}

// HttpResponseInterface レスポンス出力のためのインターフェース
type HttpResponseInterface interface {
	Success(writer http.ResponseWriter, response interface{})
	Failed(writer http.ResponseWriter, request *http.Request, err error)
	BadRequest(writer http.ResponseWriter, message string)
	InternalServerError(writer http.ResponseWriter, message string)
}

// NewHttpResponse レスポンス出力のための構造体の初期化をする
func NewHttpResponse() *HttpResponse {
	return &HttpResponse{}
}

// Success HTTPコード:200 正常終了を処理する
func (hr *HttpResponse) Success(writer http.ResponseWriter, response interface{}) {
	if response == nil {
		return
	}
	data, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		hr.InternalServerError(writer, "marshal error")
		return
	}
	writer.Write(data)
}

// Failed リクエスト失敗時のエラー処理
func (hr *HttpResponse) Failed(writer http.ResponseWriter, request *http.Request, err error) {
	var appErr derror.ApplicationError
	if errors.As(err, &appErr) {
		switch appErr.Code {
		case http.StatusBadRequest:
			hr.BadRequest(writer, appErr.Msg)
		case http.StatusInternalServerError:
			hr.InternalServerError(writer, appErr.Msg)
		}
	} else {
		hr.InternalServerError(writer, "Unknown Internal Server Error")
	}
	logging.AccessLogging(request, err)
}

// BadRequest HTTPコード:400 BadRequestを処理する
func (hr *HttpResponse) BadRequest(writer http.ResponseWriter, message string) {
	httpError(writer, http.StatusBadRequest, message)
}

// InternalServerError HTTPコード:500 InternalServerErrorを処理する
func (hr *HttpResponse) InternalServerError(writer http.ResponseWriter, message string) {
	httpError(writer, http.StatusInternalServerError, message)
}

// httpError エラー用のレスポンス出力を行う
func httpError(writer http.ResponseWriter, code int, message string) {
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
