package e2e

import (
	"InfecShotAPI/pkg/server"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRankingListIntegration(t *testing.T) {
	// モックサーバー
	mux := http.NewServeMux()
	mux.HandleFunc("/test/ranking/list", server.Get(testRankingHandler.HandleRankingList))
	server := httptest.NewServer(mux)
	defer server.Close()

	type request struct {
		method  string
		pattern string
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
			name: "normal: get a rank",
			before: func() {
				// シードの作成
				query := `INSERT INTO user(id, auth_token, name, high_score) VALUES ("test-user-id-1", "test-auth-token-1", "test-user-name-1", 100), ("test-user-id-2", "test-auth-token-2", "test-user-name-2", 1000), ("test-user-id-3", "test-auth-token-3", "test-user-name-3", 10000)`
				_, err := testUserRepository.Conn.Exec(query)
				if err != nil {
					t.Errorf("db.TestConn.Exec() failed %s", err)
					return
				}
			},
			request: request{
				method:  "GET",
				pattern: "/test/ranking/list?start=1",
			},
			after: func(res *http.Response) {
				// シードの削除
				query := `DELETE FROM user WHERE id in ("test-user-id-1", "test-user-id-2", "test-user-id-3")`
				_, err := testUserRepository.Conn.Exec(query)
				if err != nil {
					t.Errorf("db.TestConn.Exec() failed %s", err)
					return
				}
				defer res.Body.Close()
			},
			want: want{
				statusCode: http.StatusOK,
				body: `{
						  "ranks": [
							{
							  "userId": "test-user-id-1",
							  "userName": "test-user-name-1",
							  "rank": 1,
							  "score": 100
							},
							{
							  "userId": "test-user-id-2",
							  "userName": "test-user-name-2",
							  "rank": 2,
							  "score": 1000
							},
							{
							  "userId": "test-user-id-3",
							  "userName": "test-user-name-3",
							  "rank": 3,
							  "score": 10000
							}
						  ]
						}`,
			},
		},
		{
			name: "normal: get a rank with paging",
			before: func() {
				// シードの作成
				query := `INSERT INTO user(id, auth_token, name, high_score) VALUES ("test-user-id-1", "test-auth-token-1", "test-user-name-1", 100), ("test-user-id-2", "test-auth-token-2", "test-user-name-2", 1000), ("test-user-id-3", "test-auth-token-3", "test-user-name-3", 10000)`
				_, err := testUserRepository.Conn.Exec(query)
				if err != nil {
					t.Errorf("db.TestConn.Exec() failed %s", err)
					return
				}
			},
			request: request{
				method:  "GET",
				pattern: "/test/ranking/list?start=2",
			},
			after: func(res *http.Response) {
				// シードの削除
				query := `DELETE FROM user WHERE id in ("test-user-id-1", "test-user-id-2", "test-user-id-3")`
				_, err := testUserRepository.Conn.Exec(query)
				if err != nil {
					t.Errorf("db.TestConn.Exec() failed %s", err)
					return
				}
				defer res.Body.Close()
			},
			want: want{
				statusCode: http.StatusOK,
				body: `{
						  "ranks": [
							{
							  "userId": "test-user-id-2",
							  "userName": "test-user-name-2",
							  "rank": 2,
							  "score": 1000
							},
							{
							  "userId": "test-user-id-3",
							  "userName": "test-user-name-3",
							  "rank": 3,
							  "score": 10000
							}
						  ]
						}`,
			},
		},
		{
			name: "abnormal: method not allowed",
			before: func() {
			},
			request: request{
				method:  "POST",
				pattern: "/test/ranking/list?start=1",
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
			name: "abnormal: invalid query parameter",
			before: func() {
				// シードの作成
				query := `INSERT INTO user(id, auth_token, name, high_score) VALUES ("test-user-id-1", "test-auth-token-1", "test-user-name-1", 100), ("test-user-id-2", "test-auth-token-2", "test-user-name-2", 1000), ("test-user-id-3", "test-auth-token-3", "test-user-name-3", 10000)`
				_, err := testUserRepository.Conn.Exec(query)
				if err != nil {
					t.Errorf("db.TestConn.Exec() failed %s", err)
				}
			},
			request: request{
				method:  "GET",
				pattern: "/test/ranking/list?notstart=1",
			},
			after: func(res *http.Response) {
				// シードの削除
				query := `DELETE FROM user WHERE id in ("test-user-id-1", "test-user-id-2", "test-user-id-3")`
				_, err := testUserRepository.Conn.Exec(query)
				if err != nil {
					t.Errorf("db.TestConn.Exec() failed %s", err)
				}
				defer res.Body.Close()
			},
			want: want{
				statusCode: http.StatusBadRequest,
				body: `{
							"code": 400,
							"message": "Bad Request"
						}`,
			},
		},
		{
			name: "abnormal: invalid query value",
			before: func() {
				// シードの作成
				query := `INSERT INTO user(id, auth_token, name, high_score) VALUES ("test-user-id-1", "test-auth-token-1", "test-user-name-1", 100), ("test-user-id-2", "test-auth-token-2", "test-user-name-2", 1000), ("test-user-id-3", "test-auth-token-3", "test-user-name-3", 10000)`
				_, err := testUserRepository.Conn.Exec(query)
				if err != nil {
					t.Errorf("db.TestConn.Exec() failed %s", err)
				}
			},
			request: request{
				method:  "GET",
				pattern: "/test/ranking/list?start=-1",
			},
			after: func(res *http.Response) {
				// シードの削除
				query := `DELETE FROM user WHERE id in ("test-user-id-1", "test-user-id-2", "test-user-id-3")`
				_, err := testUserRepository.Conn.Exec(query)
				if err != nil {
					t.Errorf("db.TestConn.Exec() failed %s", err)
				}
				defer res.Body.Close()
			},
			want: want{
				statusCode: http.StatusBadRequest,
				body: `{
							"code": 400,
							"message": "Bad Request"
						}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before() // シード作成

			// リクエスト
			req, err := http.NewRequest(tt.request.method, server.URL+tt.request.pattern, nil)
			if err != nil {
				t.Errorf("http.NewRequest faild: %v", err)
				return
			}

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
