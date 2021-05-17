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

	pattern := "/test/user/create"
	type request struct {
		method string
		body   *strings.Reader
	}
	type want struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name    string
		before  func(mock *mockUUID)
		request request
		want    want
		after   func(*http.Response)
	}{
		{
			name: "normal: create a user",
			before: func(mock *mockUUID) {
				mock.UUID.EXPECT().Get().Return("test-uuid", nil).Times(2)
			},
			request: request{
				method: "POST",
				body:   strings.NewReader(`{"name": "test-user-name"}`),
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
		{
			name: "abnormal: bad request",
			before: func(mock *mockUUID) {
			},
			request: request{
				method: "POST",
				body:   strings.NewReader(`{"name": 100}`),
			},
			want: want{
				statusCode: http.StatusBadRequest,
				body: `{
							"code":400,
							"message": "Bad Request"
						}`,
			},
			after: func(res *http.Response) {
				defer res.Body.Close()
			},
		},
		{
			name: "abnormal: method not allowed",
			before: func(mock *mockUUID) {
			},
			request: request{
				method: "GET",
				body:   strings.NewReader(`{"name": "test-user-name"}`),
			},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				body: `{
							"code":405,
							"message": "Method Not Allowed"
						}`,
			},
			after: func(res *http.Response) {
				defer res.Body.Close()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before(mock)

			// リクエスト
			req, err := http.NewRequest(tt.request.method, server.URL+pattern, tt.request.body)
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

func TestUserGetIntegration(t *testing.T) {
	testUserService := service.NewUserService(testUserRepository, nil)
	testUserHandler := handler.NewUserHandler(httpResponse, testUserService)

	mux := http.NewServeMux() // モックサーバー
	mux.HandleFunc("/test/user/get", get(testAuthMiddleware.Authenticate(testUserHandler.HandleUserGet)))
	server := httptest.NewServer(mux)
	defer server.Close()

	pattern := "/test/user/get"
	type request struct {
		method string
		token  string
	}
	type want struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name    string
		before  func()
		request request
		want    want
		after   func(*http.Response)
	}{
		{
			name: "normal: get a user",
			before: func() {
				// シードの作成
				query := `INSERT INTO user(id, auth_token, name, high_score) VALUES ("test-user-id", "test-auth-token", "test-user-name", 100)`
				_, err := testUserRepository.Conn.Exec(query)
				if err != nil {
					t.Errorf("db.TestConn.Exec() failed %s", err)
					return
				}
			},
			request: request{
				method: "GET",
				token:  "test-auth-token",
			},
			want: want{
				statusCode: http.StatusOK,
				body: `{
							"id": "test-user-id",
							"name": "test-user-name",
							"highScore": 100
						}`,
			},
			after: func(res *http.Response) {
				// シードの削除
				query := `DELETE FROM user WHERE id = "test-user-id"`
				_, err := testUserRepository.Conn.Exec(query)
				if err != nil {
					t.Errorf("db.TestConn.Exec() failed %s", err)
					return
				}
				defer res.Body.Close()
			},
		},
		{
			name: "abnormal: invalid token",
			before: func() {
				// シードの作成
				query := `INSERT INTO user(id, auth_token, name, high_score) VALUES ("test-user-id", "test-auth-token", "test-user-name", 100)`
				_, err := testUserRepository.Conn.Exec(query)
				if err != nil {
					t.Errorf("db.TestConn.Exec() failed %s", err)
					return
				}
			},
			request: request{
				method: "GET",
				token:  "test-invalid-auth-token",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				body: `{
							"code":400,
							"message": "Bad Request"
						}`,
			},
			after: func(res *http.Response) {
				// シードの削除
				query := `DELETE FROM user WHERE id = "test-user-id"`
				_, err := testUserRepository.Conn.Exec(query)
				if err != nil {
					t.Errorf("db.TestConn.Exec() failed %s", err)
					return
				}
				defer res.Body.Close()
			},
		},
		{
			name: "abnormal: method not allowed",
			before: func() {
				// シードの作成
				query := `INSERT INTO user(id, auth_token, name, high_score) VALUES ("test-user-id", "test-auth-token", "test-user-name", 100)`
				_, err := testUserRepository.Conn.Exec(query)
				if err != nil {
					t.Errorf("db.TestConn.Exec() failed %s", err)
					return
				}
			},
			request: request{
				method: "POST",
				token:  "test-auth-token",
			},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				body: `{
							"code":405,
							"message": "Method Not Allowed"
						}`,
			},
			after: func(res *http.Response) {
				// シードの削除
				query := `DELETE FROM user WHERE id = "test-user-id"`
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
			tt.before()

			// リクエスト
			req, err := http.NewRequest(tt.request.method, server.URL+pattern, nil)
			if err != nil {
				t.Errorf("http.NewRequest faild: %v", err)
				return
			}
			req.Header.Set("x-token", tt.request.token)

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
