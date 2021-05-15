package handler

import (
	"InfecShotAPI/pkg/http/response/mock_response"
	"InfecShotAPI/pkg/server/service/mock_service"

	"github.com/golang/mock/gomock"
)

type mock struct {
	mockHttpResponse *mock_response.MockHttpResponseInterface
	mockUserService  *mock_service.MockUserServiceInterface
}

func newMock(ctrl *gomock.Controller) *mock {
	return &mock{
		mockHttpResponse: mock_response.NewMockHttpResponseInterface(ctrl),
		mockUserService:  mock_service.NewMockUserServiceInterface(ctrl),
	}
}
