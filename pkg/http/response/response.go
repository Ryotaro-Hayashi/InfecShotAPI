//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock_$GOPACKAGE/mock_$GOFILE

package response

import (
	"InfecShotAPI/pkg/derror"
	"InfecShotAPI/pkg/logging"
	"encoding/json"
	"errors"
	"net/http"
)

// httpResponse レスポンス出力のための構造体
type httpResponse struct{}

// NewHttpResponse レスポンス出力のための構造体の初期化をする
func NewHttpResponse() HttpResponse {
	return &httpResponse{}
}

// HttpResponse レスポンス出力のためのインターフェース
type HttpResponse interface {
	Success(writer http.ResponseWriter, request *http.Request, response interface{})
	Failed(writer http.ResponseWriter, request *http.Request, err error)
}

var _ HttpResponse = (*httpResponse)(nil)

// Success HTTPコード:200 正常終了を処理する
func (hr *httpResponse) Success(writer http.ResponseWriter, request *http.Request, response interface{}) {
	if response == nil {
		writer.WriteHeader(http.StatusNoContent)
		logging.AccessLogging(request, nil)
		return
	}
	data, err := json.Marshal(response)
	if err != nil {
		hr.Failed(writer, request, derror.InternalServerError.Wrap(err))
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write(data)
	logging.AccessLogging(request, err)
}

// Failed リクエスト失敗時のエラー処理
func (hr *httpResponse) Failed(writer http.ResponseWriter, request *http.Request, err error) {
	var appErr derror.ApplicationError
	if errors.As(err, &appErr) {
		HttpError(writer, appErr.Code, appErr.Msg)
	} else {
		HttpError(writer, http.StatusInternalServerError, "Unknown Internal Server Error")
	}
	logging.ApplicationErrorLogging(request, err)
	logging.AccessLogging(request, err)
}

// HttpError エラー用のレスポンス出力を行う
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
