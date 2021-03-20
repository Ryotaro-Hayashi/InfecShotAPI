package handler

import (
	"InfecShotAPI/pkg/dcontext"
	"InfecShotAPI/pkg/derror"
	"InfecShotAPI/pkg/http/response"
	"InfecShotAPI/pkg/server/service"
	"encoding/json"
	"errors"
	"net/http"
)

type gameFinishRequest struct {
	Score int `json:"score"`
}

type GameHandler struct {
	HttpResponse response.HttpResponseInterface
	GameService  service.GameServiceInterface
}

func NewGameHandler(httpResponse response.HttpResponseInterface, gameService service.GameServiceInterface) *GameHandler {
	return &GameHandler{
		HttpResponse: httpResponse,
		GameService:  gameService,
	}
}

// HandleGameFinish インゲーム終了
func (h *GameHandler) HandleGameFinish(writer http.ResponseWriter, request *http.Request) {
	// リクエストbodyからスコアを取得
	var requestBody gameFinishRequest
	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
		// TODO:アプリケーションログ
		//log.Println(err)
		h.HttpResponse.Failed(writer, request, derror.BadRequestError.Wrap(err))
		return
	}
	// scoreが負の数のときエラーを返す
	if requestBody.Score < 0 {
		// TODO:アプリケーションログ
		//log.Println(fmt.Sprintf("score is minus. score=%d", requestBody.Score))
		h.HttpResponse.Failed(writer, request, derror.BadRequestError.Wrap(errors.New("score is minus")))
		return
	}

	// ミドルウェアでコンテキストに格納したユーザidの取得
	ctx := request.Context()
	userID := dcontext.GetUserIDFromContext(ctx)
	if userID == "" {
		// TODO:アプリケーションログ
		//log.Println("userID from context is empty")
		h.HttpResponse.Failed(writer, request, derror.InternalServerError.Wrap(errors.New("failed to find user")))
		return
	}

	// ゲーム終了時のロジック
	if err := h.GameService.FinishGame(&service.FinishGameRequest{
		UserId: userID,
		Score:  requestBody.Score,
	}); err != nil {
		// TODO:アプリケーションログ
		//log.Println(err)
		h.HttpResponse.Failed(writer, request, derror.StackError(err))
		return
	}

	// 獲得コインをレスポンスとして返す
	h.HttpResponse.Success(writer, request, nil)
}
