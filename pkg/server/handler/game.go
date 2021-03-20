package handler

import (
	"InfecShotAPI/pkg/dcontext"
	"InfecShotAPI/pkg/http/response"
	"InfecShotAPI/pkg/server/service"
	"encoding/json"
	"fmt"
	"log"
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
		log.Println(err)
		h.HttpResponse.BadRequest(writer, "Bad Request")
		return
	}
	// scoreが負の数のときエラーを返す
	if requestBody.Score < 0 {
		log.Println(fmt.Sprintf("score is minus. score=%d", requestBody.Score))
		h.HttpResponse.BadRequest(writer, "Bad Request")
		return
	}

	// ミドルウェアでコンテキストに格納したユーザidの取得
	ctx := request.Context()
	userID := dcontext.GetUserIDFromContext(ctx)
	if userID == "" {
		log.Println("userID from context is empty")
		h.HttpResponse.InternalServerError(writer, "Internal Server Error")
		return
	}

	// ゲーム終了時のロジック
	if err := h.GameService.FinishGame(&service.FinishGameRequest{
		UserId: userID,
		Score:  requestBody.Score,
	}); err != nil {
		log.Println(err)
		h.HttpResponse.InternalServerError(writer, "Internal Server Error")
		return
	}

	// 獲得コインをレスポンスとして返す
	h.HttpResponse.Success(writer, nil)
}
