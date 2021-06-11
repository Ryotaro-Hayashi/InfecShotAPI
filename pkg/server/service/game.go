//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock_$GOPACKAGE/mock_$GOFILE

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

type gameService struct {
	UserRepository model.UserRepository
}

func NewGameService(userRepository model.UserRepository) GameService {
	return &gameService{
		UserRepository: userRepository,
	}
}

type GameService interface {
	FinishGame(serviceRequest *FinishGameRequest) error
}

var _ GameService = (*gameService)(nil)

// FinishGame ゲーム終了時のロジック
func (s *gameService) FinishGame(serviceRequest *FinishGameRequest) error {
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

	return nil
}
