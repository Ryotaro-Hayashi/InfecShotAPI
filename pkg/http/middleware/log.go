package middleware

import (
	"InfecShotAPI/pkg/dcontext"
	"InfecShotAPI/pkg/logging"
	"context"
	"net/http"

	"github.com/google/uuid"
)

func AccessLogging(nextFunc http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		requestID, err := uuid.NewRandom()
		if err != nil {
			// TODO:Error構造体を作成
		}
		ctx = dcontext.SetRequestID(ctx, requestID.String())
		request = request.WithContext(ctx)
		logging.AccessLogging(request, err)

		nextFunc(writer, request)
	}
}
