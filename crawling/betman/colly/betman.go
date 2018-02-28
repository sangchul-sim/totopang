package colly

import (
	"fmt"
	"strings"

	"time"

	"github.com/gocolly/colly"
	"github.com/sangchul-sim/totopang_kit/crawling"
	"github.com/sangchul-sim/totopang_kit/crawling/betman"
)

var req *crawling.Request

func init() {
	req = crawling.NewRequest()
	req.
		//SetMethod(crawling.RequestMethodGet).
		//SetUrl(listUrl).
		SetAgentByOs(crawling.OsWindows)
}

func getRecordGameHitResultDetailByDividend(tr *colly.HTMLElement) (dividends []*betman.ProtoRecordMatchDividend) {
	tr.ForEach("div.viwRap div.lst tbody tr", func(n int, tr *colly.HTMLElement) {
		var (
			homeAwayStr string
			homeAway    []string
			dividend    betman.ProtoRecordMatchDividend
		)
		tr.ForEach("td", func(k int, td *colly.HTMLElement) {
			text := strings.TrimSpace(td.Text)
			switch k {
			case 0: // dividend id
				dividend.DividendID = text
			case 1: // score
				homeAwayStr = text
				if strings.Contains(homeAwayStr, "-") {
					homeAway = strings.Split(homeAwayStr, "-")
					dividend.HomeScore = homeAway[0]
					dividend.AwayScore = homeAway[1]
				} else {
					dividend.HomeScore = homeAwayStr
					dividend.AwayScore = ""
				}

			case 2: // hit
				dividend.DividendRate = td.ChildText("span")
				imgSrc := td.ChildAttr("img", "src")
				if strings.Contains(imgSrc, "ico_chkon") {
					dividend.IsHit = true
				}
			}
			dividends = append(dividends, &dividend)
		})
	})
	return
}

func getRecordGameHitResultDetailByMatch(tr *colly.HTMLElement) *betman.ProtoRecordMatch {
	var match betman.ProtoRecordMatch
	tr.ForEach("td", func(n int, td *colly.HTMLElement) {
		text := strings.TrimSpace(td.Text)
		switch n {
		case 0: // 게임 A~Z
			match.RoundNo = text
		case 1: // 종목
			alt := td.ChildAttr("img", "alt")
			if val, ok := betman.GameType[alt]; ok {
				match.GameType = val
			}
		case 2: // 경기일, Sun Feb 11 23:15:00 KST 2018
			match.MatchTime, _ = time.Parse("Mon Jan 02 15:04:05 MST 2006", text)
		case 3: // 경기시각
		case 4: // 게임주제
			match.GameTitle = text
		case 5: // 프로토결과
			match.GameResult = text
		}
	})
	return &match
}

func getRecordGameHitResultDetail(url string) (
	match *betman.ProtoRecordMatch,
	dividends []*betman.ProtoRecordMatchDividend,
) {
	c := colly.NewCollector(
		colly.UserAgent(req.UserAgent),
	)
	//var match *betman.ProtoRecordMatch
	//var dividends []*betman.ProtoRecordMatchDividend
	//var idx int
	c.OnHTML("#tblSort tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(n int, tr *colly.HTMLElement) {
			switch n {
			case 0:
				match = getRecordGameHitResultDetailByMatch(tr)
			case 1:
				dividends = getRecordGameHitResultDetailByDividend(tr)

			}
		})
	})

	// OnResponse X
	//
	//c.OnScraped(func(r *colly.Response) {
	//	fmt.Println("match:", match)
	//	for _, dividend := range dividends {
	//		fmt.Println("dividend:", dividend)
	//	}
	//
	//})

	c.Visit(url)
	c.Wait()

	return
}

// TODO https://github.com/gocolly/colly/blob/master/_examples/reddit/reddit.go 참고해서 async 로 구현할 것
func GetRecordGameHitResultNew(page int) {
	listUrl, err := betman.GetRecordGameHitResultListUrl(page)
	if err != nil {
		panic(err)
	}
	fmt.Println(listUrl)

	c := colly.NewCollector(
		colly.UserAgent(req.UserAgent),
	)

	var detailUrls []string
	c.OnHTML("#contents tbody tr td a", func(e *colly.HTMLElement) {
		if resultPage, err := betman.GetHitResulPage(betman.DetailKeyPageMap, betman.GameTypeRecord); err == nil {
			detailUrl := strings.Join([]string{
				betman.BaseUrl,
				"/",
				resultPage,
				"?",
				betman.NewUrlParamFromQuery(e.Attr("href")).BuildQuery(),
			}, "")
			detailUrls = append(detailUrls, detailUrl)
		}
	})

	// Before making a request print "Visiting ..."
	//c.OnRequest(func(r *colly.Request) {
	//	fmt.Println("Visiting", r.URL.String())
	//
	//})

	c.Visit(listUrl)
	c.Wait()

	for _, url := range detailUrls {
		match, dividends := getRecordGameHitResultDetail(url)
		fmt.Println("match:", match)
		for _, dividend := range dividends {
			fmt.Println("dividend:", dividend)
		}
	}
}
