package server

import (
	"InfecShotAPI/pkg/db"
	"InfecShotAPI/pkg/http/middleware"
	"InfecShotAPI/pkg/http/response"
	"InfecShotAPI/pkg/server/handler"
	"InfecShotAPI/pkg/server/model"
	"InfecShotAPI/pkg/server/service"
	"InfecShotAPI/pkg/utils"
	"fmt"
	"log"
	"net/http"
)

var (
	httpResponse = response.NewHttpResponse()

	accessMiddleware = middleware.NewAccessMiddleware(httpResponse)
	userRepository   = model.NewUserRepository(db.Conn)
	UUID             = utils.NewUUID()
	authMiddleware   = middleware.NewAuthMiddleware(httpResponse, userRepository)

	userService = service.NewUserService(userRepository, UUID)

	gameService    = service.NewGameService(userRepository)
	rankingService = service.NewRankingService(userRepository)

	userHandler    = handler.NewUserHandler(httpResponse, userService)
	gameHandler    = handler.NewGameHandler(httpResponse, gameService)
	rankingHandler = handler.NewRankingHandler(httpResponse, rankingService)
)

// Serve HTTPサーバを起動する
func Serve(addr string) {
	http.HandleFunc("/user/create", accessMiddleware.Access(post(userHandler.HandleUserCreate)))
	http.HandleFunc("/user/get", accessMiddleware.Access(get(authMiddleware.Authenticate(userHandler.HandleUserGet))))
	http.HandleFunc("/game/finish", accessMiddleware.Access(post(authMiddleware.Authenticate(gameHandler.HandleGameFinish))))
	http.HandleFunc("/ranking/list", accessMiddleware.Access(get(rankingHandler.HandleRankingList)))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})

	/* ===== サーバの起動 ===== */
	log.Println("Server running...")
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Listen and serve failed. %+v", err)
	}
}

// get GETリクエストを処理する
func get(apiFunc http.HandlerFunc) http.HandlerFunc {
	return httpMethod(apiFunc, http.MethodGet)
}

// post POSTリクエストを処理する
func post(apiFunc http.HandlerFunc) http.HandlerFunc {
	return httpMethod(apiFunc, http.MethodPost)
}

// httpMethod 指定したHTTPメソッドでAPIの処理を実行する
func httpMethod(apiFunc http.HandlerFunc, method string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		// CORS対応
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Access-Control-Allow-Headers", "Content-Type,Accept,Origin,x-token")

		// プリフライトリクエストは処理を通さない
		if request.Method == http.MethodOptions {
			return
		}
		// 指定のHTTPメソッドでない場合はエラー
		if request.Method != method {
			response.HttpError(writer, http.StatusMethodNotAllowed, "Method Not Allowed")
			return
		}

		// 共通のレスポンスヘッダを設定
		writer.Header().Add("Content-Type", "application/json")
		apiFunc(writer, request)
	}
}
