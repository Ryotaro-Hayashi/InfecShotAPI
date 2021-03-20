package service

import (
	"InfecShotAPI/pkg/derror"
	"InfecShotAPI/pkg/server/model"
	"errors"
)

type FinishGameRequest struct {
	UserId string
	Score  int
}

type GameService struct {
	UserRepository model.UserRepositoryInterface
}

func NewGameService(userRepository model.UserRepositoryInterface) *GameService {
	return &GameService{
		UserRepository: userRepository,
	}
}

type GameServiceInterface interface {
	FinishGame(serviceRequest *FinishGameRequest) error
}

var _ GameServiceInterface = (*GameService)(nil)

// GameFinish ゲーム終了時のロジック
func (s *GameService) FinishGame(serviceRequest *FinishGameRequest) error {
	// ゲーム終了前のユーザ情報の取得
	user, err := s.UserRepository.SelectUserByPrimaryKey(serviceRequest.UserId)
	if err != nil {
		return derror.StackError(err)
	}
	if user == nil {
		return derror.BadRequestError.Wrap(errors.New("failed to find user"))
	}

	// ユーザのハイスコアとリクエストのスコアを比較
	if user.HighScore > serviceRequest.Score || user.HighScore == 0 {
		user.HighScore = serviceRequest.Score
		// ハイスコアを更新
		if err = s.UserRepository.UpdateUserByPrimaryKey(user); err != nil {
			return derror.StackError(err)
		}
	}

	return err
}
