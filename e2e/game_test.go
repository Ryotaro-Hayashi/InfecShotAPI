package e2e

import (
	"InfecShotAPI/pkg/server"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGameFinishIntegration(t *testing.T) {
	// モックサーバー
	mux := http.NewServeMux()
	mux.HandleFunc("/test/game/finish", server.Post(testAuthMiddleware.Authenticate(testGameHandler.HandleGameFinish)))
	server := httptest.NewServer(mux)
	defer server.Close()

	type request struct {
		method  string
		pattern string
		body    *strings.Reader
		token   string
	}
	type want struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name    string
		before  func()
		request request
		after   func(*http.Response)
		want    want
	}{
		{
			name: "normal: finish game with updating score",
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
				method:  "POST",
				pattern: "/test/game/finish",
				body:    strings.NewReader(`{"score": 1000}`),
				token:   "test-auth-token",
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
			want: want{
				statusCode: http.StatusNoContent,
			},
		},
		{
			name: "normal: finish game with no updating score",
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
				method:  "POST",
				pattern: "/test/game/finish",
				body:    strings.NewReader(`{"score": 10}`),
				token:   "test-auth-token",
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
			want: want{
				statusCode: http.StatusNoContent,
			},
		},
		{
			name: "normal: finish game with fist play",
			before: func() {
				// シードの作成
				query := `INSERT INTO user(id, auth_token, name, high_score) VALUES ("test-user-id", "test-auth-token", "test-user-name", 0)`
				_, err := testUserRepository.Conn.Exec(query)
				if err != nil {
					t.Errorf("db.TestConn.Exec() failed %s", err)
					return
				}
			},
			request: request{
				method:  "POST",
				pattern: "/test/game/finish",
				body:    strings.NewReader(`{"score": 1000}`),
				token:   "test-auth-token",
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
			want: want{
				statusCode: http.StatusNoContent,
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
				method:  "POST",
				pattern: "/test/game/finish",
				body:    strings.NewReader(`{"score": 1000}`),
				token:   "test-invalid-auth-token",
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
			want: want{
				statusCode: http.StatusBadRequest,
				body: `{
							"code":400,
							"message": "Bad Request"
						}`,
			},
		},
		{
			name: "abnormal: method not allowed",
			before: func() {
			},
			request: request{
				method:  "GET",
				pattern: "/test/game/finish",
				body:    strings.NewReader(`{"score": 1000}`),
				token:   "test-invalid-auth-token",
			},
			after: func(res *http.Response) {
			},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				body: `{
							"code":405,
							"message": "Method Not Allowed"
						}`,
			},
		},
		{
			name: "abnormal: bad request",
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
				method:  "POST",
				pattern: "/test/game/finish",
				body:    strings.NewReader(`{"score": -1000}`),
				token:   "test-auth-token",
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
			want: want{
				statusCode: http.StatusBadRequest,
				body: `{
							"code":400,
							"message": "Bad Request"
						}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before() // シード作成

			// リクエスト
			req, err := http.NewRequest(tt.request.method, server.URL+tt.request.pattern, tt.request.body)
			if err != nil {
				t.Errorf("http.NewRequest faild: %v", err)
				return
			}
			req.Header.Set("x-token", tt.request.token)

			// 実行してレスポンスを取得
			client := http.DefaultClient
			res, err := client.Do(req)
			if err != nil {
				t.Errorf("http.DefaultClient.Do failed: %v", err)
				return
			}

			if res.StatusCode != tt.want.statusCode {
				t.Errorf("status code = %d, want %d", res.StatusCode, tt.want.statusCode)
			}

			tt.after(res)
		})
	}
}
