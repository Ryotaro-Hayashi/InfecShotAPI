package service

import (
	"InfecShotAPI/pkg/server/model"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestUserService_CreateUser(t *testing.T) {
	type args struct {
		serviceRequest *CreateUserRequest
	}
	tests := []struct {
		name    string
		args    args
		before  func(mock *mockRepository, mockUUID *mockUUID, args args)
		want    *CreateUserResponse
		wantErr bool
	}{
		{
			name: "normal: create a user",
			args: args{
				serviceRequest: &CreateUserRequest{
					Name: "test-user-name",
				},
			},
			// model層のmockとuuidのmock
			before: func(mock *mockRepository, mockUUID *mockUUID, args args) {
				mock.userRepository.EXPECT().InsertUser(&model.User{
					ID:        "test-uuid",
					AuthToken: "test-uuid",
					Name:      args.serviceRequest.Name,
					HighScore: 0,
				}).Return(nil)
				mockUUID.UUID.EXPECT().Get().Return("test-uuid", nil).Times(2)
			},
			want: &CreateUserResponse{
				Token: "test-uuid",
			},
			wantErr: false,
		},
		{
			name: "abnormal: error in generate uuid",
			args: args{
				serviceRequest: &CreateUserRequest{
					Name: "test-user-name",
				},
			},
			before: func(mock *mockRepository, mockUUID *mockUUID, args args) {
				mockUUID.UUID.EXPECT().Get().Return("", errors.New("failed to generate uuid"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "abnormal: error in model.InsertUser()",
			args: args{
				serviceRequest: &CreateUserRequest{
					Name: "test-user-name",
				},
			},
			before: func(mock *mockRepository, mockUUID *mockUUID, args args) {
				mock.userRepository.EXPECT().InsertUser(&model.User{
					ID:        "test-uuid",
					AuthToken: "test-uuid",
					Name:      args.serviceRequest.Name,
					HighScore: 0,
				}).Return(errors.New("failed to InsertUser()"))
				mockUUID.UUID.EXPECT().Get().Return("test-uuid", nil).Times(2)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			// mockの初期化
			mock := newMockRepository(ctrl)
			mockUUID := newMockUUID(ctrl)
			tt.before(mock, mockUUID, tt.args) // mockの作成

			s := NewUserService(mock.userRepository, mockUUID.UUID)
			got, err := s.CreateUser(tt.args.serviceRequest)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	type args struct {
		serviceRequest *GetUserRequest
	}
	tests := []struct {
		name    string
		args    args
		before  func(mock *mockRepository, args args)
		want    *GetUserResponse
		wantErr bool
	}{
		{
			name: "normal: get a user",
			args: args{
				serviceRequest: &GetUserRequest{
					ID: "test-user-id",
				},
			},
			before: func(mock *mockRepository, args args) {
				mock.userRepository.EXPECT().SelectUserByPrimaryKey(args.serviceRequest.ID).Return(&model.User{
					ID:        "test-user-id",
					AuthToken: "test-auth-token",
					Name:      "test-user-name",
					HighScore: 100,
				}, nil)
			},
			want: &GetUserResponse{
				ID:        "test-user-id",
				Name:      "test-user-name",
				HighScore: 100,
			},
			wantErr: false,
		},
		{
			name: "abnormal: failed to model.SelectUserByPrimaryKey()",
			args: args{
				serviceRequest: &GetUserRequest{
					ID: "test-user-id",
				},
			},
			before: func(mock *mockRepository, args args) {
				mock.userRepository.EXPECT().SelectUserByPrimaryKey(args.serviceRequest.ID).Return(nil,
					errors.New("failed to model.SelectUserByPrimaryKey()"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "abnormal: nil user",
			args: args{
				serviceRequest: &GetUserRequest{
					ID: "test-user-id",
				},
			},
			before: func(mock *mockRepository, args args) {
				mock.userRepository.EXPECT().SelectUserByPrimaryKey(args.serviceRequest.ID).Return(nil, nil)
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
			s := NewUserService(mock.userRepository, nil)
			got, err := s.GetUser(tt.args.serviceRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
