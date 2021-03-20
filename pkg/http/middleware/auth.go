package middleware

import (
	"InfecShotAPI/pkg/dcontext"
	"InfecShotAPI/pkg/http/response"
	"InfecShotAPI/pkg/server/model"
	"context"
	"errors"
	"log"
	"net/http"
)

type Middleware struct {
	HttpResponse   response.HttpResponseInterface
	UserRepository model.UserRepositoryInterface
}

func NewMiddleware(httpResponse response.HttpResponseInterface, userRepository model.UserRepositoryInterface) *Middleware {
	return &Middleware{
		HttpResponse:   httpResponse,
		UserRepository: userRepository,
	}
}

// Authenticate ユーザ認証を行ってContextへユーザID情報を保存する
func (m *Middleware) Authenticate(nextFunc http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		ctx := request.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		// リクエストヘッダからx-token(認証トークン)を取得
		token := request.Header.Get("x-token")
		if token == "" {
			log.Println("x-token is empty")
			m.HttpResponse.BadRequest(writer, "Bad Request")
			return
		}

		user, err := m.UserRepository.SelectUserByAuthToken(token)
		if err != nil {
			log.Println(err)
			m.HttpResponse.InternalServerError(writer, "Internal Server Error")
			return
		}
		if user == nil {
			log.Println(errors.New("user not found"))
			m.HttpResponse.InternalServerError(writer, "Internal Server Error")
			return
		}

		// ユーザIDをContextへ保存して以降の処理に利用する
		ctx = dcontext.SetUserID(ctx, user.ID)

		// 次の処理
		nextFunc(writer, request.WithContext(ctx))
	}
}
