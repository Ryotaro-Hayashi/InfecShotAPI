package handler

import (
	"InfecShotAPI/pkg/server/service"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestRankingHandler_HandleRankingList(t *testing.T) {
	method := "GET"
	baseUrl := "http://localhost:" + addr + "/ranking/list"
	type args struct {
		writer  http.ResponseWriter
		request *http.Request
	}
	tests := []struct {
		name   string
		args   args
		before func(mock *mock, args args)
	}{
		{
			name: "normal: get ranking information",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(method, baseUrl+"?start=1", nil),
			},
			before: func(mock *mock, args args) {
				mock.mockRankingService.EXPECT().GetRankInfoList(&service.GetRankInfoListRequest{
					Offset: 1,
					Limit:  10,
				}).Return(&service.GetRankInfoListResponse{RankInfoList: []*service.RankInfo{
					{
						UserId:   "test-user-id-1",
						UserName: "test-user-name-1",
						Rank:     1,
						Score:    100,
					},
					{
						UserId:   "test-user-id-2",
						UserName: "test-user-name-2",
						Rank:     2,
						Score:    1000,
					},
				}}, nil)
				mock.mockHttpResponse.EXPECT().Success(args.writer, args.request, rankingListResponse{Ranks: []*rank{
					{
						UserId:   "test-user-id-1",
						UserName: "test-user-name-1",
						Rank:     1,
						Score:    100,
					},
					{
						UserId:   "test-user-id-2",
						UserName: "test-user-name-2",
						Rank:     2,
						Score:    1000,
					},
				}}).Return()
			},
		},
		{
			name: "abnormal: failed to strconv.Atoi()",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(method, baseUrl+"?notstart=1", nil),
			},
			before: func(mock *mock, args args) {
				mock.mockHttpResponse.EXPECT().Failed(args.writer, args.request, gomock.Any()).Return()
			},
		},
		{
			name: "abnormal: validation error",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(method, baseUrl+"?start=-1", nil),
			},
			before: func(mock *mock, args args) {
				mock.mockHttpResponse.EXPECT().Failed(args.writer, args.request, gomock.Any()).Return()
			},
		},
		{
			name: "abnormal: failed to service.GetRankInfoList()",
			args: args{
				writer:  httptest.NewRecorder(),
				request: httptest.NewRequest(method, baseUrl+"?start=1", nil),
			},
			before: func(mock *mock, args args) {
				mock.mockRankingService.EXPECT().GetRankInfoList(&service.GetRankInfoListRequest{
					Offset: 1,
					Limit:  10,
				}).Return(nil, errors.New("failed to service.GetRankInfoList()"))
				mock.mockHttpResponse.EXPECT().Failed(args.writer, args.request, gomock.Any()).Return()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := newMock(ctrl)
			tt.before(mock, tt.args)

			writer := httptest.NewRecorder()

			h := NewRankingHandler(mock.mockHttpResponse, mock.mockRankingService)
			h.HandleRankingList(writer, tt.args.request)
		})
	}
}
