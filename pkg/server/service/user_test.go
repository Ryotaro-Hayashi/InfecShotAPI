package service

import (
	"InfecShotAPI/pkg/server/model"
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
		want    *createUserResponse
		wantErr bool
	}{
		{
			name: "test",
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
			want: &createUserResponse{
				Token: "test-uuid",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := newMockRepository(ctrl)
			mockUUID := newMockUUID(ctrl)
			tt.before(mock, mockUUID, tt.args)

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

//func TestUserService_GetUser(t *testing.T) {
//	type fields struct {
//		UserRepository model.UserRepositoryInterface
//	}
//	type args struct {
//		serviceRequest *GetUserRequest
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *getUserResponse
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &UserService{
//				UserRepository: tt.fields.UserRepository,
//			}
//			got, err := s.GetUser(tt.args.serviceRequest)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
