package betman

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/sangchul-sim/totopang_kit/crawling"
	"github.com/suapapa/go_hangul/encoding/cp949"
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
	BaseUrl = "http://betman.co.kr"
	//BaseUrl                      = "http://localhost:9999/mock/batman"
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
	listKeyPageMap               = "list"
	detailKeyPageMap             = "detail"
	gameTypeRecord               = "record"
	gameTypeWinLose              = "winLose"
	gameTypeSoccerResult         = "soccerResult"
)

type pageMap map[string]map[string]string

var gameSchedulePageMap = pageMap{
	listKeyPageMap: {
		gameTypeWinLose:      targetPageGameScheduleList,
		gameTypeRecord:       targetPageGameScheduleList,
		gameTypeSoccerResult: targetPageGameScheduleList,
	},
	detailKeyPageMap: {
		gameTypeWinLose:      targetPageGameScheduleList,
		gameTypeRecord:       targetPageGameScheduleList,
		gameTypeSoccerResult: targetPageGameScheduleList,
	},
}

var hitResultPageMap = pageMap{
	listKeyPageMap: {
		gameTypeWinLose:      targetPageWinningResultList,
		gameTypeRecord:       targetPageWinningResultList,
		gameTypeSoccerResult: targetPageWinningResultList,
	},
	detailKeyPageMap: {
		gameTypeWinLose:      targetPageWinningResultProto,
		gameTypeRecord:       targetPageWinningResultProto,
		gameTypeSoccerResult: targetPageWinningResultToto,
	},
}

var GameType = map[string]string{
	"야구": "baseball",
	"농구": "basketball",
	"축구": "soccer",
	"배구": "volleyball",
	"골프": "golf",
}

func GetRecordGameHitResultNew(page int) {
	//doc, err := goquery.NewDocument("http://data.7m.com.cn/result_data/default_kr2.shtml?date=2018-02-22")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	//	e.Request.Visit(e.Attr("href"))
	//})
}

func GetRecordGameHitResult(page int) {
	listUrl, err := getRecordGameHitResultListUrl(page)
	if err != nil {
		panic(err)
	}
	fmt.Println(listUrl)

	request := crawling.NewRequest()
	b, err := request.
		SetMethod(crawling.RequestMethodGet).
		SetUrl(listUrl).
		SetAgentByOs(crawling.OsWindows).
		Do()

	//b, err := crawling.RequestUrl(listUrl, crawling.RequestMethodGet, nil)
	if err != nil {
		panic(err)
	}
	utf8b, err := cp949.From(b)
	if err != nil {
		panic(err)
	}
	for i, detail := range getRecordGameHitResultDetailParam(utf8b) {
		if resultPage, err := getHitResulPage(detailKeyPageMap, gameTypeRecord); err == nil {
			detailUrl := strings.Join([]string{BaseUrl, "/", resultPage, "?", detail.BuildQuery()}, "")
			fmt.Println(detailUrl)
			if i == 0 {
				b, err := request.
					SetMethod(crawling.RequestMethodGet).
					SetUrl(detailUrl).
					//SetAgentByOs(crawling.OsWindows).
					Do()
				//b, err := crawling.RequestUrl(detailUrl, crawling.RequestMethodGet, nil)
				if err != nil {
					panic(err)
				}
				utf8b, err := cp949.From(b)
				if err != nil {
					panic(err)
				}
				getRecordGameHitResultDetail(utf8b)
			}
		}
	}
}

func GetPage() {
	//fmt.Println("gameSchedulePageMap", gameSchedulePageMap)
	//
	//val, err := getHitResulPage("list", "WinLose")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println("getHitResulPage", val)
	//
	//listParam := GameScheduleListParam("a238dsfsd", "201902")
	//query := listParam.BuildQuery()
	//fmt.Println("GameScheduleListParam:", query)
	//
	//listParam2 := GameScheduleListParam("a238dsfsd", "201902")
	//query2 := listParam2.BuildQuery()
	//fmt.Println("GameScheduleListParam:", query2)
	//
	//detailParam := GameScheduleDetailParam("a238dsfsd", "round0002983", "201902", "", "")
	//query = detailParam.BuildQuery()
	//fmt.Println("GameScheduleDetailParam:", query)

	//dom.ExampleGetElementsByTagName()
	//return

	GetRecordGameHitResult(1)

	// http://www.betman.co.kr/winningResultProto.so?method=detail&gameId=G102&gameRound=180366&page=1&selectedGameId=G102
	// http://www.betman.co.kr/winningResultProto.so?method=detail&gameId=G102&gameRound=180366&selectedGameId=G102
}

// 기록식 적중결과 상세
func getRecordGameHitResultDetailUrl(gameRound string, page int) (string, error) {
	resultPage, err := getHitResulPage(detailKeyPageMap, gameTypeRecord)
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
func getHitResulPage(key1, key2 string) (string, error) {
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
