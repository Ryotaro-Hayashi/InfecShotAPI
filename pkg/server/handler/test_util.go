package handler

import (
	"InfecShotAPI/pkg/http/response/mock_response"
	"InfecShotAPI/pkg/server/service/mock_service"

	"github.com/golang/mock/gomock"
)

var addr = "80"

type mock struct {
	mockHttpResponse   *mock_response.MockHttpResponseInterface
	mockUserService    *mock_service.MockUserServiceInterface
	mockRankingService *mock_service.MockRankingServiceInterface
	mockGameService    *mock_service.MockGameServiceInterface
}

func newMock(ctrl *gomock.Controller) *mock {
	return &mock{
		mockHttpResponse:   mock_response.NewMockHttpResponseInterface(ctrl),
		mockUserService:    mock_service.NewMockUserServiceInterface(ctrl),
		mockRankingService: mock_service.NewMockRankingServiceInterface(ctrl),
		mockGameService:    mock_service.NewMockGameServiceInterface(ctrl),
	}
}
