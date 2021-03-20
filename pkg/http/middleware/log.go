package middleware

import (
	"InfecShotAPI/pkg/dcontext"
	"InfecShotAPI/pkg/derror"
	"InfecShotAPI/pkg/http/response"
	"context"
	"net/http"

	"github.com/google/uuid"
)

type accessMiddleware struct {
	HttpResponse response.HttpResponseInterface
}

func NewAccessMiddleware(httpResponse response.HttpResponseInterface) *accessMiddleware {
	return &accessMiddleware{
		HttpResponse: httpResponse,
	}
}

func (m *accessMiddleware) Access(nextFunc http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		requestID, err := uuid.NewRandom()
		if err != nil {
			err = derror.InternalServerError.Wrap(err)
			// TODO:アプリケーションログの出力
		}
		if requestID != uuid.Nil {
			ctx = dcontext.SetRequestID(ctx, requestID.String())
			request = request.WithContext(ctx)
		}

		nextFunc(writer, request)
	}
}
