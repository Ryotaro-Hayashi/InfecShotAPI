package handler

import (
	"InfecShotAPI/pkg/dcontext"
	"InfecShotAPI/pkg/derror"
	"InfecShotAPI/pkg/http/response"
	"InfecShotAPI/pkg/logging"
	"InfecShotAPI/pkg/server/service"
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
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
	requestID := dcontext.GetRequestIDFromContext(request.Context())
	logging.ApplicationLogger.Info("start finishing game", zap.String("requestID", requestID))

	// リクエストbodyからスコアを取得
	var requestBody gameFinishRequest
	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
		h.HttpResponse.Failed(writer, request, derror.BadRequestError.Wrap(err))
		return
	}
	logging.ApplicationLogger.Info("succeed in decoding request body", zap.String("requestID", requestID), zap.Any("requestBody", requestBody))
	// scoreが負の数のときエラーを返す
	if requestBody.Score < 0 {
		h.HttpResponse.Failed(writer, request, derror.BadRequestError.Wrap(errors.New("score is minus")))
		return
	}
	logging.ApplicationLogger.Debug("succeed in getting score", zap.String("requestID", requestID), zap.Int("score", requestBody.Score))

	// ミドルウェアでコンテキストに格納したユーザidの取得
	ctx := request.Context()
	userID := dcontext.GetUserIDFromContext(ctx)
	if userID == "" {
		h.HttpResponse.Failed(writer, request, derror.InternalServerError.Wrap(errors.New("userID from context is empty")))
		return
	}
	logging.ApplicationLogger.Info("succeed in getting userID from context", zap.String("requestID", requestID), zap.String("userID", userID))

	// ゲーム終了時のロジック
	if err := h.GameService.FinishGame(&service.FinishGameRequest{
		UserId: userID,
		Score:  requestBody.Score,
	}); err != nil {
		h.HttpResponse.Failed(writer, request, derror.StackError(err))
		return
	}
	logging.ApplicationLogger.Debug("succeed in finishing game", zap.String("requestID", requestID))

	h.HttpResponse.Success(writer, request, nil)
	logging.ApplicationLogger.Info("finished game", zap.String("requestID", requestID))
}
