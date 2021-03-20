package service

import (
	"InfecShotAPI/pkg/derror"
	"InfecShotAPI/pkg/server/model"
	"errors"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Name string
}

type createUserResponse struct {
	Token string
}

type GetUserRequest struct {
	ID string
}

type getUserResponse struct {
	ID        string
	Name      string
	HighScore int
}

type UserService struct {
	UserRepository model.UserRepositoryInterface
}

func NewUserService(userRepository model.UserRepositoryInterface) *UserService {
	return &UserService{
		UserRepository: userRepository,
	}
}

type UserServiceInterface interface {
	CreateUser(serviceRequest *CreateUserRequest) (*createUserResponse, error)
	GetUser(serviceRequest *GetUserRequest) (*getUserResponse, error)
}

var _ UserServiceInterface = (*UserService)(nil)

// CreateUser ユーザ情報作成のロジック
func (s *UserService) CreateUser(serviceRequest *CreateUserRequest) (*createUserResponse, error) {
	// UUIDでユーザIDを生成する
	userID, err := uuid.NewRandom()
	if err != nil {
		return nil, derror.InternalServerError.Wrap(err)
	}

	// UUIDで認証トークンを生成する
	authToken, err := uuid.NewRandom()
	if err != nil {
		return nil, derror.InternalServerError.Wrap(err)
	}

	// データベースにユーザデータを登録する
	if err = s.UserRepository.InsertUser(&model.User{
		ID:        userID.String(),
		AuthToken: authToken.String(),
		Name:      serviceRequest.Name,
		HighScore: 0,
	}); err != nil {
		return nil, derror.StackError(err)
	}

	return &createUserResponse{Token: authToken.String()}, nil
}

// CreateUser ユーザ情報取得のロジック
func (s *UserService) GetUser(serviceRequest *GetUserRequest) (*getUserResponse, error) {
	user, err := s.UserRepository.SelectUserByPrimaryKey(serviceRequest.ID)
	if err != nil {
		return nil, derror.StackError(err)
	}
	if user == nil {
		return nil, derror.BadRequestError.Wrap(errors.New("failed to find user"))
	}

	return &getUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		HighScore: user.HighScore,
	}, nil
}
