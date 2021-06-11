package middleware

import (
	"InfecShotAPI/pkg/dcontext"
	"InfecShotAPI/pkg/derror"
	"InfecShotAPI/pkg/http/response"
	"InfecShotAPI/pkg/logging"
	"context"
	"net/http"

	"github.com/google/uuid"
)

type accessMiddleware struct {
	HttpResponse response.HttpResponse
}

func NewAccessMiddleware(httpResponse response.HttpResponse) *accessMiddleware {
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
			logging.ApplicationErrorLogging(request, derror.InternalServerError.Wrap(err))
		}
		if requestID != uuid.Nil {
			ctx = dcontext.SetRequestID(ctx, requestID.String())
			request = request.WithContext(ctx)
		}

		nextFunc(writer, request)
	}
}
