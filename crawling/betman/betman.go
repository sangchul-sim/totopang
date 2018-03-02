package betman

import (
	"errors"
	"strconv"
)

/**
승부식	: winLose
기록식	: record
축구승무패	: soccerResult

게임일정 : game schedule
적중결과 : hit result

// 승부식 적중결과 리스트
http://betman.co.kr/winningResultList.so?page=1&gameId=G101

// 승부식 적중결과 상세
http://betman.co.kr/winningResultProto.so?method=detail&gameId=G101&gameRound=180002&page=1&selectedGameId=G101

// 승부식 게임일정 리스트
http://betman.co.kr/GameScheduleList.so?method=schedule&viewmethod=list&gameId=G101&league=&team=&selectedLeague=%EC%A0%84%EC%B2%B4&yearMonth=201802

// 승부식 게임 상세
http://betman.co.kr/GameScheduleList.so?method=detailSchedule&gameId=G101&gameRound=180011&seqNo=26751&viewmethod=list&yearMonth=201802&gameIds=G101&league=&team=&page=
*/

const (
	//BaseUrl                      = "http://localhost:9999/mock/batman"
	BaseUrl                      = "http://betman.co.kr"
	GameIDWinLose                = "G101" // 승부식
	GameIDRecord                 = "G102" // 기록식
	GameIDSoccerResult           = "G011" // 축구승무패
	targetPageWinningResultToto  = "winningResultToto.so"
	targetPageWinningResultProto = "winningResultProto.so"
	targetPageGameScheduleList   = "GameScheduleList.so"
	targetPageWinningResultList  = "winningResultList.so"
	paramMethodScheduleList      = "schedule"
	paramMethodScheduleDetail    = "detailSchedule"
	paramMethodHitResultDetail   = "detail"
	paramSelectedLeagueAll       = "전체"
	paramViewMethodList          = "list"
	ListKeyPageMap               = "list"
	DetailKeyPageMap             = "detail"
	GameTypeRecord               = "record"
	GameTypeWinLose              = "winLose"
	GameTypeSoccerResult         = "soccerResult"
)

type pageMap map[string]map[string]string

var gameSchedulePageMap = pageMap{
	ListKeyPageMap: {
		GameTypeWinLose:      targetPageGameScheduleList,
		GameTypeRecord:       targetPageGameScheduleList,
		GameTypeSoccerResult: targetPageGameScheduleList,
	},
	DetailKeyPageMap: {
		GameTypeWinLose:      targetPageGameScheduleList,
		GameTypeRecord:       targetPageGameScheduleList,
		GameTypeSoccerResult: targetPageGameScheduleList,
	},
}

var hitResultPageMap = pageMap{
	ListKeyPageMap: {
		GameTypeWinLose:      targetPageWinningResultList,
		GameTypeRecord:       targetPageWinningResultList,
		GameTypeSoccerResult: targetPageWinningResultList,
	},
	DetailKeyPageMap: {
		GameTypeWinLose:      targetPageWinningResultProto,
		GameTypeRecord:       targetPageWinningResultProto,
		GameTypeSoccerResult: targetPageWinningResultToto,
	},
}

var GameType = map[string]string{
	"야구": "baseball",
	"농구": "basketball",
	"축구": "soccer",
	"배구": "volleyball",
	"골프": "golf",
}

// 기록식 적중결과 상세
func getRecordGameHitResultDetailUrl(gameRound string, page int) (string, error) {
	resultPage, err := GetHitResulPage(DetailKeyPageMap, GameTypeRecord)
	if err != nil {
		return "", err
	}
	queryString := NewHitResultDetailParam(GameIDRecord, gameRound, strconv.Itoa(page)).BuildQuery()
	return BaseUrl + "/" + resultPage + "?" + queryString, nil
}

func GameSchedulePage(key1, key2 string) (string, error) {
	keyMap, ok := gameSchedulePageMap[key1]
	if !ok {
		return "", errors.New(key1 + " not found")
	}
	val, ok := keyMap[key2]
	if !ok {
		return "", errors.New(key2 + " not found")
	}
	return val, nil
}

// key1:list or detail
// key2:WinLose or record or soccerResult
func GetHitResulPage(key1, key2 string) (string, error) {
	keyMap, ok := hitResultPageMap[key1]
	if !ok {
		return "", errors.New(key1 + " not found")
	}
	val, ok := keyMap[key2]
	if !ok {
		return "", errors.New(key2 + " not found")
	}
	return val, nil
}

func NewHitResultListParam(gameID, p string) *UrlParam {
	param := newUrlParam()
	param.GameId = gameID
	param.Page = p
	return param
}

func NewHitResultDetailParam(gameID, roundID, page string) *UrlParam {
	param := newUrlParam()
	param.Method = paramMethodHitResultDetail
	param.GameId = gameID
	param.GameRound = roundID
	param.SelectedGameId = gameID
	param.Page = page
	return param
}

func NewGameScheduleListParam(gameID, yearMonth string) *UrlParam {
	return &UrlParam{
		Method:         paramMethodScheduleList,
		ViewMethod:     paramViewMethodList,
		GameId:         gameID,
		SelectedLeague: paramSelectedLeagueAll,
		YearMonth:      yearMonth,
	}
}

func NewGameScheduleDetailParam(gameID, roundID, yearMonth, year, outerRound string) *UrlParam {
	return &UrlParam{
		Method:         paramMethodScheduleDetail,
		ViewMethod:     paramViewMethodList,
		GameId:         gameID,
		GameIds:        gameID,
		GameRound:      roundID,
		SelectedLeague: paramSelectedLeagueAll,
		YearMonth:      yearMonth,
		OuterRound:     outerRound,
		SaleYear:       year,
	}
}
