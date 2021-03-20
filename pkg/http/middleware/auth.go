package middleware

import (
	"InfecShotAPI/pkg/dcontext"
	"InfecShotAPI/pkg/derror"
	"InfecShotAPI/pkg/http/response"
	"InfecShotAPI/pkg/server/model"
	"context"
	"errors"
	"net/http"
)

type authMiddleware struct {
	HttpResponse   response.HttpResponseInterface
	UserRepository model.UserRepositoryInterface
}

func NewAuthMiddleware(httpResponse response.HttpResponseInterface, userRepository model.UserRepositoryInterface) *authMiddleware {
	return &authMiddleware{
		HttpResponse:   httpResponse,
		UserRepository: userRepository,
	}
}

// Authenticate ユーザ認証を行ってContextへユーザID情報を保存する
func (m *authMiddleware) Authenticate(nextFunc http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		ctx := request.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		// リクエストヘッダからx-token(認証トークン)を取得
		token := request.Header.Get("x-token")
		if token == "" {
			// TODO:アプリケーションログ
			m.HttpResponse.Failed(writer, request, derror.BadRequestError.Wrap(errors.New("failed to get token")))
			return
		}

		user, err := m.UserRepository.SelectUserByAuthToken(token)
		if err != nil {
			// TODO:アプリケーションログ
			//log.Println(err)
			m.HttpResponse.Failed(writer, request, derror.StackError(err))
			return
		}
		if user == nil {
			// TODO:アプリケーションログ
			//log.Println(errors.New("user not found"))
			m.HttpResponse.Failed(writer, request, derror.InternalServerError.Wrap(errors.New("empty set user")))
			return
		}

		// ユーザIDをContextへ保存して以降の処理に利用する
		ctx = dcontext.SetUserID(ctx, user.ID)

		// 次の処理
		nextFunc(writer, request.WithContext(ctx))
	}
}
