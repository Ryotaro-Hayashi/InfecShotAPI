package db

import (
	"InfecShotAPI/pkg/logging"
	"database/sql"
	"fmt"
	"os"

	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
)

// Driver名
const driverName = "mysql"

// Conn 各repositoryで利用するDB接続(Connection)情報
var Conn *sql.DB

func init() {
	/* ===== データベースへ接続する. ===== */
	// ユーザ
	user := os.Getenv("MYSQL_USER")
	// パスワード
	password := os.Getenv("MYSQL_PASSWORD")
	// 接続先ホスト
	//host := os.Getenv("MYSQL_HOST")
	// 接続先ポート
	port := os.Getenv("MYSQL_PORT")
	// 接続先データベース
	database := os.Getenv("MYSQL_DATABASE")

	// 接続情報は以下のように指定する.
	// user:password@tcp(host:port)/database
	var err error
	Conn, err = sql.Open(driverName,
		fmt.Sprintf("%s:%s@tcp(mysql:%s)/%s", user, password, port, database))
	logging.ApplicationLogger.Info(fmt.Sprintf("%s:%s@tcp(mysql:%s)/%s\n", user, password, port, database))
	if err != nil {
		logging.ApplicationLogger.Fatal("failed to sql.Open()", zap.Error(err))
	}
	if err = Conn.Ping(); err != nil {
		logging.ApplicationLogger.Fatal("failed to sql.DB.Ping()", zap.Error(err))
	}
}
