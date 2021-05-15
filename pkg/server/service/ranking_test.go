package service

import (
	"InfecShotAPI/pkg/server/model"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestRankingService_GetRankInfoList(t *testing.T) {
	type args struct {
		serviceRequest *GetRankInfoListRequest
	}
	tests := []struct {
		name    string
		args    args
		before  func(mock *mockRepository, args args)
		want    *GetRankInfoListResponse
		wantErr bool
	}{
		{
			name: "normal: get ranking information",
			args: args{
				serviceRequest: &GetRankInfoListRequest{
					Limit:  2,
					Offset: 1,
				},
			},
			before: func(mock *mockRepository, args args) {
				mock.userRepository.EXPECT().SelectUsersOrderByHighScoreAsc(args.serviceRequest.Limit, args.serviceRequest.Offset).Return([]*model.User{
					{
						ID:        "test-user-id-1",
						AuthToken: "test-auth-token-1",
						Name:      "test-user-name-1",
						HighScore: 1000,
					},
					{
						ID:        "test-user-id-2",
						AuthToken: "test-auth-token-2",
						Name:      "test-user-name-2",
						HighScore: 100,
					},
				}, nil)
			},
			want: &GetRankInfoListResponse{
				RankInfoList: []*RankInfo{
					{
						UserId:   "test-user-id-1",
						UserName: "test-user-name-1",
						Rank:     1,
						Score:    1000,
					},
					{
						UserId:   "test-user-id-2",
						UserName: "test-user-name-2",
						Rank:     2,
						Score:    100,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "abnormal: failed to model.SelectUsersOrderByHighScoreAsc()",
			args: args{
				serviceRequest: &GetRankInfoListRequest{
					Limit:  2,
					Offset: 1,
				},
			},
			before: func(mock *mockRepository, args args) {
				mock.userRepository.EXPECT().SelectUsersOrderByHighScoreAsc(args.serviceRequest.Limit, args.serviceRequest.Offset).Return(nil, errors.New("failed to model.SelectUsersOrderByHighScoreAsc()"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := newMockRepository(ctrl)
			tt.before(mock, tt.args)
			s := NewRankingService(mock.userRepository)
			got, err := s.GetRankInfoList(tt.args.serviceRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRankInfoList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRankInfoList() got = %v, want %v", got, tt.want)
			}
		})
	}
}
