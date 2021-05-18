package e2e

import (
	"InfecShotAPI/pkg/db"
	"InfecShotAPI/pkg/http/middleware"
	"InfecShotAPI/pkg/http/response"
	"InfecShotAPI/pkg/server/handler"
	"InfecShotAPI/pkg/server/model"
	"InfecShotAPI/pkg/server/service"
	"InfecShotAPI/pkg/utils/mock_utils"
	"encoding/json"
	"reflect"

	"github.com/golang/mock/gomock"
)

var (
	testHttpResponse   = response.NewHttpResponse()
	testUserRepository = model.NewUserRepository(db.Conn)
	testAuthMiddleware = middleware.NewAuthMiddleware(testHttpResponse, testUserRepository)

	testRankingService = service.NewRankingService(testUserRepository)
	testGameService    = service.NewGameService(testUserRepository)

	testRankingHandler = handler.NewRankingHandler(testHttpResponse, testRankingService)
	testGameHandler    = handler.NewGameHandler(testHttpResponse, testGameService)
)

type mockUUID struct {
	UUID *mock_utils.MockUUIDInterface
}

func newMockUUID(ctrl *gomock.Controller) *mockUUID {
	return &mockUUID{
		UUID: mock_utils.NewMockUUIDInterface(ctrl),
	}
}

// deepEqualString 文字列同士を比較する
func deepEqualString(str1, str2 string) (bool, error) {
	if str1 == str2 {
		return true, nil
	} else {
		var strInterface1 interface{}
		err := json.Unmarshal([]byte(str1), &strInterface1)
		if err != nil {
			return false, err
		}

		var strInterface2 interface{}
		err = json.Unmarshal([]byte(str2), &strInterface2)
		if err != nil {
			return false, err
		}

		if reflect.DeepEqual(strInterface1, strInterface2) {
			return true, nil
		} else {
			return false, nil
		}
	}
}
