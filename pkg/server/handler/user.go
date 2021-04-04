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

type UserHandler struct {
	HttpResponse response.HttpResponseInterface
	UserService  service.UserServiceInterface
}

func NewUserHandler(httpResponse response.HttpResponseInterface, userService service.UserServiceInterface) *UserHandler {
	return &UserHandler{
		HttpResponse: httpResponse,
		UserService:  userService,
	}
}

type userCreateRequest struct {
	Name string `json:"name"`
}

type userCreateResponse struct {
	Token string `json:"token"`
}

type userGetResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	HighScore int    `json:"highScore"`
}

// HandleUserCreate ユーザ情報作成処理
func (h *UserHandler) HandleUserCreate(writer http.ResponseWriter, request *http.Request) {
	requestID := dcontext.GetRequestIDFromContext(request.Context())
	logging.ApplicationLogger.Info("start creating user", zap.String("requestID", requestID))

	// リクエストBodyから作成後情報を取得
	var requestBody userCreateRequest
	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
		h.HttpResponse.Failed(writer, request, derror.BadRequestError.Wrap(err))
		return
	}
	logging.ApplicationLogger.Info("succeed in decoding request body",
		zap.String("requestID", requestID),
		zap.Any("requestBody", requestBody))

	// ユーザ情報作成のロジック
	res, err := h.UserService.CreateUser(&service.CreateUserRequest{Name: requestBody.Name})
	if err != nil {
		h.HttpResponse.Failed(writer, request, derror.StackError(err))
		return
	}
	logging.ApplicationLogger.Debug("succeed in creating user", zap.String("requestID", requestID))

	// 生成した認証トークンを返却
	h.HttpResponse.Success(writer, request, &userCreateResponse{Token: res.Token})
	logging.ApplicationLogger.Info("finished creating user", zap.String("requestID", requestID))
}

// HandleUserGet ユーザ情報取得処理
func (h *UserHandler) HandleUserGet(writer http.ResponseWriter, request *http.Request) {
	requestID := dcontext.GetRequestIDFromContext(request.Context())
	userID := dcontext.GetUserIDFromContext(request.Context())
	logging.ApplicationLogger.Info("start getting user",
		zap.String("requestID", requestID),
		zap.String("userID", userID))

	if userID == "" {
		h.HttpResponse.Failed(writer, request, derror.InternalServerError.Wrap(errors.New("userID from context is empty")))
		return
	}
	logging.ApplicationLogger.Debug("succeed in getting userID from context", zap.String("requestID", requestID))

	// ユーザ情報取得のロジック
	res, err := h.UserService.GetUser(&service.GetUserRequest{ID: userID})
	if err != nil {
		h.HttpResponse.Failed(writer, request, derror.StackError(err))
		return
	}
	logging.ApplicationLogger.Debug("succeed in getting user", zap.String("requestID", requestID))

	// レスポンスに必要な情報を詰めて返却
	h.HttpResponse.Success(writer, request, &userGetResponse{
		ID:        res.ID,
		Name:      res.Name,
		HighScore: res.HighScore,
	})
	logging.ApplicationLogger.Info("finished getting user", zap.String("requestID", requestID))
}
