package server

import (
	"InfecShotAPI/pkg/server/handler"
	"InfecShotAPI/pkg/server/service"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestUserCreateIntegration(t *testing.T) {
	ctrl := gomock.NewController(t)
	// mockの初期化
	mock := newMockUUID(ctrl)
	testUserService := service.NewUserService(testUserRepository, mock.UUID)
	testUserHandler := handler.NewUserHandler(httpResponse, testUserService)

	mux := http.NewServeMux() // モックサーバー
	mux.HandleFunc("/test/user/create", post(testUserHandler.HandleUserCreate))
	server := httptest.NewServer(mux)
	defer server.Close()

	method := "POST"
	pattern := "/test/user/create"
	requestBody := strings.NewReader(`{"name": "test-user-name"}`)
	type want struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name   string
		before func(mock *mockUUID)
		want   want
		after  func(*http.Response)
	}{
		{
			name: "normal: get a user",
			before: func(mock *mockUUID) {
				mock.UUID.EXPECT().Get().Return("test-uuid", nil).Times(2)
			},
			want: want{
				statusCode: http.StatusOK,
				body: `{
							"token": "test-uuid"
						}`,
			},
			after: func(res *http.Response) {
				// シードの削除
				query := `DELETE FROM user WHERE id = "test-uuid"`
				_, err := testUserRepository.Conn.Exec(query)
				if err != nil {
					t.Errorf("db.TestConn.Exec() failed %s", err)
					return
				}
				defer res.Body.Close()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before(mock)

			// リクエスト
			req, err := http.NewRequest(method, server.URL+pattern, requestBody)
			if err != nil {
				t.Errorf("http.NewRequest faild: %v", err)
				return
			}

			// 実行してレスポンスを取得
			client := http.DefaultClient
			res, err := client.Do(req)
			if err != nil {
				t.Errorf("http.DefaultClient.Do() failed: %v", err)
				return
			}

			if res.StatusCode != tt.want.statusCode {
				t.Errorf("status code = %d, want %d", res.StatusCode, tt.want.statusCode)
			}

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("ioutil.ReadAll failed: %v", err)
				return
			}

			boolean, err := deepEqualString(string(body), tt.want.body)
			if err != nil {
				t.Errorf("response.DeepEqualString() failed %s", err)
			}
			if !boolean {
				t.Errorf("response body = \n%s\n, want \n%s\n", string(body), tt.want.body)
			}

			tt.after(res)
		})
	}
}
