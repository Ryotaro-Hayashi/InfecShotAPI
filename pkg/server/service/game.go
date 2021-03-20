package service

import (
	"InfecShotAPI/pkg/server/model"
	"errors"
	"fmt"
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
		return err
	}
	if user == nil {
		return errors.New(fmt.Sprintf("user not found. userID=%s", serviceRequest.UserId))
	}

	// ユーザのハイスコアとリクエストのスコアを比較
	if user.HighScore > serviceRequest.Score || user.HighScore == 0 {
		user.HighScore = serviceRequest.Score
		// ハイスコアを更新
		if err = s.UserRepository.UpdateUserByPrimaryKey(user); err != nil {
			return err
		}
	}

	return err
}
