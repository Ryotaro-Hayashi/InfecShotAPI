package service

import (
	"InfecShotAPI/pkg/server/model/mock_model"
	"InfecShotAPI/pkg/utils/mock_utils"

	"github.com/golang/mock/gomock"
)

type mockRepository struct {
	userRepository *mock_model.MockUserRepositoryInterface
}

func newMockRepository(ctrl *gomock.Controller) *mockRepository {
	return &mockRepository{
		userRepository: mock_model.NewMockUserRepositoryInterface(ctrl),
	}
}

type mockUUID struct {
	UUID *mock_utils.MockUUIDInterface
}

func newMockUUID(ctrl *gomock.Controller) *mockUUID {
	return &mockUUID{
		UUID: mock_utils.NewMockUUIDInterface(ctrl),
	}
}
