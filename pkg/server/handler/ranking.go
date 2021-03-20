package handler

import (
	"InfecShotAPI/pkg/derror"
	"InfecShotAPI/pkg/http/response"
	"InfecShotAPI/pkg/server/service"
	"errors"
	"net/http"
	"strconv"
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
	HttpResponse   response.HttpResponseInterface
	RankingService service.RankingServiceInterface
}

func NewRankingHandler(httpResponse response.HttpResponseInterface, rankingService service.RankingServiceInterface) *RankingHandler {
	return &RankingHandler{
		HttpResponse:   httpResponse,
		RankingService: rankingService,
	}
}

// HandleRankingList ランキング情報取得
func (h *RankingHandler) HandleRankingList(writer http.ResponseWriter, request *http.Request) {
	// クエリストリングから開始順位の受け取り
	param := request.URL.Query().Get("start")
	start, err := strconv.Atoi(param)
	if err != nil {
		// TODO:アプリケーションログ
		//log.Println(err)
		h.HttpResponse.Failed(writer, request, derror.BadRequestError.Wrap(err))
		return
	}
	// startが0以下のときエラーを返す
	if start <= 0 {
		// TODO:アプリケーションログ
		//log.Println(fmt.Sprintf("start rank is 0 or less. start=%d", start))
		h.HttpResponse.Failed(writer, request, derror.BadRequestError.Wrap(errors.New("failed to get start rank")))
		return
	}

	// ランキング情報取得のロジック
	res, err := h.RankingService.GetRankInfoList(&service.GetRankInfoListRequest{
		Offset: start,
		Limit:  10,
	})
	if err != nil {
		// TODO:アプリケーションログ
		//log.Println(err)
		h.HttpResponse.Failed(writer, request, derror.StackError(err))
		return
	}

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
}
