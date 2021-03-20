package handler

import (
	"2103_proto_f_server/pkg/dcontext"
	"2103_proto_f_server/pkg/http/response"
	"2103_proto_f_server/pkg/server/service"
	"encoding/json"
	"log"
	"net/http"
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
	// リクエストBodyから作成後情報を取得
	var requestBody userCreateRequest
	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
		log.Println(err)
		h.HttpResponse.BadRequest(writer, "Bad Request")
		return
	}

	// ユーザ情報作成のロジック
	res, err := h.UserService.CreateUser(&service.CreateUserRequest{Name: requestBody.Name})
	if err != nil {
		log.Println(err)
		h.HttpResponse.InternalServerError(writer, "Internal Server Error")
		return
	}

	// 生成した認証トークンを返却
	h.HttpResponse.Success(writer, &userCreateResponse{Token: res.Token})
}

// HandleUserGet ユーザ情報取得処理
func (h *UserHandler) HandleUserGet(writer http.ResponseWriter, request *http.Request) {
	// Contextから認証済みのユーザIDを取得
	ctx := request.Context()
	userID := dcontext.GetUserIDFromContext(ctx)
	if userID == "" {
		log.Println("userID from context is empty")
		h.HttpResponse.InternalServerError(writer, "Internal Server Error")
		return
	}

	// ユーザ情報取得のロジック
	res, err := h.UserService.GetUser(&service.GetUserRequest{ID: userID})
	if err != nil {
		log.Println(err)
		h.HttpResponse.InternalServerError(writer, "Internal Server Error")
		return
	}

	// レスポンスに必要な情報を詰めて返却
	h.HttpResponse.Success(writer, &userGetResponse{
		ID:        res.ID,
		Name:      res.Name,
		HighScore: res.HighScore,
	})
}
