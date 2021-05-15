package handler

import (
	"InfecShotAPI/pkg/server/service"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
)

var addr = "80"

func TestUserHandler_HandleUserCreate(t *testing.T) {
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
				request: httptest.NewRequest("POST", url, strings.NewReader(`{"name": "test-user-name"}`)),
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
			name: "abnormal: failed to service.CreateUser()",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest("POST", url, strings.NewReader(`{"name": "test-user-name"}`)),
			},
			before: func(mock *mock, args args) {
				err := errors.New("failed to service.CreateUser()")
				mock.mockUserService.EXPECT().CreateUser(&service.CreateUserRequest{Name: "test-user-name"}).Return(nil, err)
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

// contextのmock化が必要
//func TestUserHandler_HandleUserGet(t *testing.T) {
//	url := "http://localhost:" + addr + "/user/get"
//	type args struct {
//		writer  http.ResponseWriter
//		request *http.Request
//	}
//	tests := []struct {
//		name   string
//		args   args
//		before func(mock *mock, args args)
//	}{
//		{
//			name: "normal: get a user",
//			args: args{
//				writer:  httptest.NewRecorder(),
//				request: httptest.NewRequest("GET", url, nil),
//			},
//			before: func(mock *mock, args args) {
//				mock.mockUserService.EXPECT().GetUser(&service.GetUserRequest{ID: "test-user-id"}).Return(&service.GetUserResponse{
//					ID:        "test-user-id",
//					Name:      "test-user-name",
//					HighScore: 1000,
//				}, nil)
//				mock.mockHttpResponse.EXPECT().Success(args.writer, args.request, &userGetResponse{
//					ID:        "test-user-id",
//					Name:      "test-user-name",
//					HighScore: 1000,
//				}).Return()
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			mock := newMock(ctrl)
//			tt.before(mock, tt.args)
//
//			writer := httptest.NewRecorder()
//
//			h := NewUserHandler(mock.mockHttpResponse, mock.mockUserService)
//			h.HandleUserGet(writer, tt.args.request)
//		})
//	}
//}
