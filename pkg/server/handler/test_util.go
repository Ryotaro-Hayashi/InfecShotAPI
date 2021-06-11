package handler

import (
	"InfecShotAPI/pkg/http/response/mock_response"
	"InfecShotAPI/pkg/server/service/mock_service"

	"github.com/golang/mock/gomock"
)

var addr = "80"

type mock struct {
	mockHttpResponse   *mock_response.MockHttpResponse
	mockUserService    *mock_service.MockUserService
	mockRankingService *mock_service.MockRankingService
	mockGameService    *mock_service.MockGameService
}

func newMock(ctrl *gomock.Controller) *mock {
	return &mock{
		mockHttpResponse:   mock_response.NewMockHttpResponse(ctrl),
		mockUserService:    mock_service.NewMockUserService(ctrl),
		mockRankingService: mock_service.NewMockRankingService(ctrl),
		mockGameService:    mock_service.NewMockGameService(ctrl),
	}
}
