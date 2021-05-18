package service

import (
	"InfecShotAPI/pkg/server/model"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGameService_FinishGame(t *testing.T) {
	type args struct {
		serviceRequest *FinishGameRequest
	}
	tests := []struct {
		name    string
		args    args
		before  func(mock *mockRepository, args args)
		wantErr bool
	}{
		{
			name: "normal: finish game with updating score",
			args: args{
				// DBのスコアより小さいのでスコア更新必要
				serviceRequest: &FinishGameRequest{
					UserId: "test-user-id",
					Score:  100,
				},
			},
			before: func(mock *mockRepository, args args) {
				mock.userRepository.EXPECT().SelectUserByPrimaryKey(args.serviceRequest.UserId).Return(&model.User{
					ID:        "test-user-id",
					AuthToken: "test-auth-token",
					Name:      "test-user-name",
					HighScore: 1000,
				}, nil)
				mock.userRepository.EXPECT().UpdateUserByPrimaryKey(&model.User{
					ID:        "test-user-id",
					AuthToken: "test-auth-token",
					Name:      "test-user-name",
					HighScore: 100,
				}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "normal: finish game with no updating",
			args: args{
				// DBのスコアより大きいのでスコア更新不要
				serviceRequest: &FinishGameRequest{
					UserId: "test-user-id",
					Score:  1000,
				},
			},
			before: func(mock *mockRepository, args args) {
				mock.userRepository.EXPECT().SelectUserByPrimaryKey(args.serviceRequest.UserId).Return(&model.User{
					ID:        "test-user-id",
					AuthToken: "test-auth-token",
					Name:      "test-user-name",
					HighScore: 1000,
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "normal: finish first game ",
			args: args{
				serviceRequest: &FinishGameRequest{
					UserId: "test-user-id",
					Score:  100,
				},
			},
			before: func(mock *mockRepository, args args) {
				mock.userRepository.EXPECT().SelectUserByPrimaryKey(args.serviceRequest.UserId).Return(&model.User{
					ID:        "test-user-id",
					AuthToken: "test-auth-token",
					Name:      "test-user-name",
					HighScore: 0,
				}, nil)
				mock.userRepository.EXPECT().UpdateUserByPrimaryKey(&model.User{
					ID:        "test-user-id",
					AuthToken: "test-auth-token",
					Name:      "test-user-name",
					HighScore: 100,
				}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "abnormal: failed to model.SelectUserByPrimaryKey()",
			args: args{
				serviceRequest: &FinishGameRequest{
					UserId: "test-user-id",
					Score:  100,
				},
			},
			before: func(mock *mockRepository, args args) {
				mock.userRepository.EXPECT().SelectUserByPrimaryKey(args.serviceRequest.UserId).Return(nil, errors.New("failed to model.SelectUserByPrimaryKey()"))
			},
			wantErr: true,
		},
		{
			name: "abnormal: nil user",
			args: args{
				serviceRequest: &FinishGameRequest{
					UserId: "test-user-id",
					Score:  100,
				},
			},
			before: func(mock *mockRepository, args args) {
				mock.userRepository.EXPECT().SelectUserByPrimaryKey(args.serviceRequest.UserId).Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name: "abnormal: failed to model.UpdateUserByPrimaryKey()",
			args: args{
				serviceRequest: &FinishGameRequest{
					UserId: "test-user-id",
					Score:  100,
				},
			},
			before: func(mock *mockRepository, args args) {
				mock.userRepository.EXPECT().SelectUserByPrimaryKey(args.serviceRequest.UserId).Return(&model.User{
					ID:        "test-user-id",
					AuthToken: "test-auth-token",
					Name:      "test-user-name",
					HighScore: 1000,
				}, nil)
				mock.userRepository.EXPECT().UpdateUserByPrimaryKey(&model.User{
					ID:        "test-user-id",
					AuthToken: "test-auth-token",
					Name:      "test-user-name",
					HighScore: 100,
				}).Return(errors.New("failed to model.UpdateUserByPrimaryKey()"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := newMockRepository(ctrl)
			tt.before(mock, tt.args)
			s := NewGameService(mock.userRepository)
			if err := s.FinishGame(tt.args.serviceRequest); (err != nil) != tt.wantErr {
				t.Errorf("FinishGame() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
