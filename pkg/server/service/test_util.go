package service

import (
	"InfecShotAPI/pkg/server/model/mock_model"
	"InfecShotAPI/pkg/utils/mock_utils"

	"github.com/golang/mock/gomock"
)

type mockRepository struct {
	userRepository *mock_model.MockUserRepository
}

func newMockRepository(ctrl *gomock.Controller) *mockRepository {
	return &mockRepository{
		userRepository: mock_model.NewMockUserRepository(ctrl),
	}
}

type mockUUID struct {
	UUID *mock_utils.MockUUID
}

func newMockUUID(ctrl *gomock.Controller) *mockUUID {
	return &mockUUID{
		UUID: mock_utils.NewMockUUID(ctrl),
	}
}
