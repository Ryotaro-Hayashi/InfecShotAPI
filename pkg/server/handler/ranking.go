package handler

import (
	"InfecShotAPI/pkg/dcontext"
	"InfecShotAPI/pkg/derror"
	"InfecShotAPI/pkg/http/response"
	"InfecShotAPI/pkg/logging"
	"InfecShotAPI/pkg/server/service"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type rankingListResponse struct {
	Ranks []*rank `json:"ranks"`
}

// rank ランキング情報
type rank struct {
	UserId   string `json:"userId"`
	UserName string `json:"userName"`
	Rank     int    `json:"rank"`
	Score    int    `json:"score"`
}

type RankingHandler struct {
	HttpResponse   response.HttpResponse
	RankingService service.RankingService
}

func NewRankingHandler(httpResponse response.HttpResponse, rankingService service.RankingService) *RankingHandler {
	return &RankingHandler{
		HttpResponse:   httpResponse,
		RankingService: rankingService,
	}
}

// HandleRankingList ランキング情報取得
func (h *RankingHandler) HandleRankingList(writer http.ResponseWriter, request *http.Request) {
	requestID := dcontext.GetRequestIDFromContext(request.Context())
	logging.ApplicationLogger.Info("start getting rank",
		zap.String("requestID", requestID))

	// クエリストリングから開始順位の受け取り
	param := request.URL.Query().Get("start")
	start, err := strconv.Atoi(param)
	if err != nil {
		h.HttpResponse.Failed(writer, request, derror.BadRequestError.Wrap(err))
		return
	}
	logging.ApplicationLogger.Info("succeed in getting query",
		zap.String("requestID", requestID),
		zap.String("query", fmt.Sprintf("start=%d", start)))
	// startが0以下のときエラーを返す
	if start <= 0 {
		h.HttpResponse.Failed(writer, request, derror.BadRequestError.Wrap(errors.New("failed to get start rank")))
		return
	}
	logging.ApplicationLogger.Debug("succeed in getting start rank",
		zap.String("requestID", requestID),
		zap.Int("startRank", start))

	// ランキング情報取得のロジック
	res, err := h.RankingService.GetRankInfoList(&service.GetRankInfoListRequest{
		Offset: start,
		Limit:  10,
	})
	if err != nil {
		h.HttpResponse.Failed(writer, request, derror.StackError(err))
		return
	}
	logging.ApplicationLogger.Debug("succeed in getting rank", zap.String("requestID", requestID))

	// レスポンスの整形
	var ranks []*rank
	for _, rankInfo := range res.RankInfoList {
		rank := &rank{
			UserId:   rankInfo.UserId,
			UserName: rankInfo.UserName,
			Rank:     rankInfo.Rank,
			Score:    rankInfo.Score,
		}
		ranks = append(ranks, rank)
	}

	h.HttpResponse.Success(writer, request, rankingListResponse{Ranks: ranks})
	logging.ApplicationLogger.Info("finished getting rank", zap.String("requestID", requestID))
}
