package model

import (
	"InfecShotAPI/pkg/derror"
	"database/sql"
)

// User userテーブルデータ
type User struct {
	ID        string
	AuthToken string
	Name      string
	HighScore int
}

type UserRepository struct {
	Conn *sql.DB
}

func NewUserRepository(conn *sql.DB) *UserRepository {
	return &UserRepository{
		Conn: conn,
	}
}

type UserRepositoryInterface interface {
	InsertUser(record *User) error
	SelectUserByAuthToken(authToken string) (*User, error)
	UpdateUserByPrimaryKey(record *User) error
	SelectUserByPrimaryKey(userID string) (*User, error)
	SelectUsersOrderByHighScoreAsc(limit int, offset int) ([]*User, error)
}

// インターフェースを満たしているかを確認
var _ UserRepositoryInterface = (*UserRepository)(nil)

// InsertUser データベースをレコードを登録する
func (r *UserRepository) InsertUser(record *User) error {
	stmt, err := r.Conn.Prepare("INSERT INTO user(id, auth_token, name, high_score) VALUES(?, ?, ?, ?)")
	if err != nil {
		return derror.DatabaseOperationError.Wrap(err)
	}
	_, err = stmt.Exec(record.ID, record.AuthToken, record.Name, record.HighScore)
	if err != nil {
		return derror.DatabaseOperationError.Wrap(err)
	}
	return nil
}

// UpdateUserByPrimaryKey 主キーを条件にレコードを更新する
func (r *UserRepository) UpdateUserByPrimaryKey(record *User) error {
	stmt, err := r.Conn.Prepare("UPDATE user SET name = ?, high_score = ? where id = ?")
	if err != nil {
		return derror.DatabaseOperationError.Wrap(err)
	}
	_, err = stmt.Exec(record.Name, record.HighScore, record.ID)
	if err != nil {
		return derror.DatabaseOperationError.Wrap(err)
	}
	return nil
}

// SelectUserByAuthToken auth_tokenを条件にレコードを取得する
func (r *UserRepository) SelectUserByAuthToken(authToken string) (*User, error) {
	row := r.Conn.QueryRow("SELECT * from user WHERE auth_token = ?", authToken)
	return convertToUser(row)
}

// SelectUserByPrimaryKey 主キーを条件にレコードを取得する
func (r *UserRepository) SelectUserByPrimaryKey(userID string) (*User, error) {
	row := r.Conn.QueryRow("SELECT * from user WHERE id = ?", userID)
	return convertToUser(row)
}

// SelectUsersOrderByHighScoreDesc ハイスコア順に指定順位から指定件数を取得する
func (r *UserRepository) SelectUsersOrderByHighScoreAsc(limit int, offset int) ([]*User, error) {
	stmt, err := r.Conn.Prepare("SELECT * FROM user WHERE high_score > 0 ORDER BY high_score Asc LIMIT ? OFFSET ?")
	if err != nil {
		return nil, derror.DatabaseOperationError.Wrap(err)
	}

	rows, err := stmt.Query(limit, offset-1)
	if err != nil {
		return nil, derror.DatabaseOperationError.Wrap(err)
	}

	return convertToUsers(rows)
}

// convertToUser rowデータをUserデータへ変換する
func convertToUser(row *sql.Row) (*User, error) {
	user := User{}
	err := row.Scan(&user.ID, &user.AuthToken, &user.Name, &user.HighScore)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, derror.DatabaseDataScanError.Wrap(err)
	}
	return &user, nil
}

// convertToUsers rowsデータをUserのスライスへ変換する
func convertToUsers(rows *sql.Rows) ([]*User, error) {
	defer rows.Close()

	var (
		users []*User
		err   error
	)

	for rows.Next() {
		user := User{}
		if err = rows.Scan(&user.ID, &user.AuthToken, &user.Name, &user.HighScore); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, derror.DatabaseDataScanError.Wrap(err)
		}
		users = append(users, &user)
	}
	return users, nil
}
