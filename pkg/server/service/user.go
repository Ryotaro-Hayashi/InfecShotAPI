//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock_$GOPACKAGE/mock_$GOFILE

package service

import (
	"InfecShotAPI/pkg/derror"
	"InfecShotAPI/pkg/server/model"
	"InfecShotAPI/pkg/utils"
	"errors"
)

type CreateUserRequest struct {
	Name string
}

type CreateUserResponse struct {
	Token string
}

type GetUserRequest struct {
	ID string
}

type GetUserResponse struct {
	ID        string
	Name      string
	HighScore int
}

type userService struct {
	UserRepository model.UserRepository
	UUID           utils.UUID
}

func NewUserService(userRepository model.UserRepository, uuid utils.UUID) UserService {
	return &userService{
		UserRepository: userRepository,
		UUID:           uuid,
	}
}

type UserService interface {
	CreateUser(serviceRequest *CreateUserRequest) (*CreateUserResponse, error)
	GetUser(serviceRequest *GetUserRequest) (*GetUserResponse, error)
}

var _ UserService = (*userService)(nil)

// CreateUser ユーザ情報作成のロジック
func (s *userService) CreateUser(serviceRequest *CreateUserRequest) (*CreateUserResponse, error) {
	// UUIDでユーザIDを生成する
	userID, err := s.UUID.Get()
	if err != nil {
		return nil, derror.InternalServerError.Wrap(err)
	}

	// UUIDで認証トークンを生成する
	authToken, err := s.UUID.Get()
	if err != nil {
		return nil, derror.InternalServerError.Wrap(err)
	}

	// データベースにユーザデータを登録する
	if err = s.UserRepository.InsertUser(&model.User{
		ID:        userID,
		AuthToken: authToken,
		Name:      serviceRequest.Name,
		HighScore: 0,
	}); err != nil {
		return nil, derror.StackError(err)
	}

	return &CreateUserResponse{Token: authToken}, nil
}

// GetUser ユーザ情報取得のロジック
func (s *userService) GetUser(serviceRequest *GetUserRequest) (*GetUserResponse, error) {
	user, err := s.UserRepository.SelectUserByPrimaryKey(serviceRequest.ID)
	if err != nil {
		return nil, derror.StackError(err)
	}
	if user == nil {
		return nil, derror.InternalServerError.Wrap(errors.New("empty set"))
	}

	return &GetUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		HighScore: user.HighScore,
	}, nil
}
