package handler

import (
	"InfecShotAPI/pkg/dcontext"
	"InfecShotAPI/pkg/server/service"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestUserHandler_HandleUserCreate(t *testing.T) {
	method := "POST"
	url := "http://localhost:" + addr + "/user/create"
	type args struct {
		writer  http.ResponseWriter
		request *http.Request
	}
	tests := []struct {
		name   string
		args   args
		before func(mock *mock, args args)
	}{
		{
			name: "normal: create a user",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(method, url, strings.NewReader(`{"name": "test-user-name"}`)),
			},
			before: func(mock *mock, args args) {
				mock.mockUserService.EXPECT().CreateUser(&service.CreateUserRequest{Name: "test-user-name"}).Return(&service.CreateUserResponse{
					Token: "test-token",
				}, nil)
				mock.mockHttpResponse.EXPECT().Success(args.writer, args.request, &userCreateResponse{
					Token: "test-token",
				}).Return()
			},
		},
		{
			name: "abnormal: failed to decode request body",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(method, url, strings.NewReader(`{"name": 100}`)),
			},
			before: func(mock *mock, args args) {
				mock.mockHttpResponse.EXPECT().Failed(args.writer, args.request, gomock.Any()).Return()
			},
		},
		{
			name: "abnormal: failed to service.CreateUser()",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(method, url, strings.NewReader(`{"name": "test-user-name"}`)),
			},
			before: func(mock *mock, args args) {
				mock.mockUserService.EXPECT().CreateUser(&service.CreateUserRequest{Name: "test-user-name"}).Return(nil, errors.New("failed to service.CreateUser()"))
				mock.mockHttpResponse.EXPECT().Failed(args.writer, args.request, gomock.Any()).Return()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := newMock(ctrl)
			tt.before(mock, tt.args)

			writer := httptest.NewRecorder()

			h := NewUserHandler(mock.mockHttpResponse, mock.mockUserService)
			h.HandleUserCreate(writer, tt.args.request)
		})
	}
}

func TestUserHandler_HandleUserGet(t *testing.T) {
	method := "GET"
	url := "http://localhost:" + addr + "/user/get"
	type args struct {
		writer  http.ResponseWriter
		request *http.Request
	}
	tests := []struct {
		name   string
		args   args
		userID string
		before func(mock *mock, args args, userID string)
	}{
		{
			name: "normal: get a user",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(method, url, nil),
			},
			userID: "test-user-id",
			before: func(mock *mock, args args, userID string) {
				mock.mockUserService.EXPECT().GetUser(&service.GetUserRequest{ID: "test-user-id"}).Return(&service.GetUserResponse{
					ID:        userID,
					Name:      "test-user-name",
					HighScore: 1000,
				}, nil)
				mock.mockHttpResponse.EXPECT().Success(args.writer, gomock.Any(), &userGetResponse{
					ID:        userID,
					Name:      "test-user-name",
					HighScore: 1000,
				}).Return()
			},
		},
		{
			name: "abnormal: userID from context is empty",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(method, url, nil),
			},
			before: func(mock *mock, args args, userID string) {
				mock.mockHttpResponse.EXPECT().Failed(args.writer, gomock.Any(), gomock.Any()).Return()
			},
		},
		{
			name: "abnormal: failed to service.GetUser()",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest("GET", url, nil),
			},
			userID: "test-user-id",
			before: func(mock *mock, args args, userID string) {
				mock.mockUserService.EXPECT().GetUser(&service.GetUserRequest{ID: "test-user-id"}).Return(nil, errors.New("failed to service.GetUser()"))
				mock.mockHttpResponse.EXPECT().Failed(args.writer, gomock.Any(), gomock.Any()).Return()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := newMock(ctrl)
			tt.before(mock, tt.args, tt.userID)

			writer := httptest.NewRecorder()

			h := NewUserHandler(mock.mockHttpResponse, mock.mockUserService)
			h.HandleUserGet(writer, tt.args.request.WithContext(dcontext.SetUserID(tt.args.request.Context(), tt.userID)))
		})
	}
}
