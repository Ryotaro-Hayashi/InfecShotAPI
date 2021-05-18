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

func TestGameHandler_HandleGameFinish(t *testing.T) {
	method := "POST"
	url := "http://localhost:" + addr + "/game/finish"
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
			name: "normal: finish game",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(method, url, strings.NewReader(`{"score": 100}`)),
			},
			userID: "test-user-id",
			before: func(mock *mock, args args, userID string) {
				mock.mockGameService.EXPECT().FinishGame(&service.FinishGameRequest{
					UserId: userID,
					Score:  100,
				}).Return(nil)
				mock.mockHttpResponse.EXPECT().Success(args.writer, gomock.Any(), nil).Return()
			},
		},
		{
			name: "abnormal: failed to decode request body",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(method, url, strings.NewReader(`{"score": "100"}`)),
			},
			before: func(mock *mock, args args, userID string) {
				mock.mockHttpResponse.EXPECT().Failed(args.writer, gomock.Any(), gomock.Any()).Return()
			},
		},
		{
			name: "abnormal: validation error",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(method, url, strings.NewReader(`{"score": -100}`)),
			},
			before: func(mock *mock, args args, userID string) {
				mock.mockHttpResponse.EXPECT().Failed(args.writer, gomock.Any(), gomock.Any()).Return()
			},
		},
		{
			name: "abnormal: userID from context is empty",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(method, url, strings.NewReader(`{"score": 100}`)),
			},
			before: func(mock *mock, args args, userID string) {
				mock.mockHttpResponse.EXPECT().Failed(args.writer, gomock.Any(), gomock.Any()).Return()
			},
		},
		{
			name: "abnormal: failed to service.FinishGame()",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(method, url, strings.NewReader(`{"score": 100}`)),
			},
			userID: "test-user-id",
			before: func(mock *mock, args args, userID string) {
				mock.mockGameService.EXPECT().FinishGame(&service.FinishGameRequest{
					UserId: userID,
					Score:  100,
				}).Return(errors.New("failed to service.FinishGame()"))
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

			h := NewGameHandler(mock.mockHttpResponse, mock.mockGameService)
			h.HandleGameFinish(writer, tt.args.request.WithContext(dcontext.SetUserID(tt.args.request.Context(), tt.userID)))
		})
	}
}
