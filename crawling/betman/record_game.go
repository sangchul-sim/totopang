package betman

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/sangchul-sim/totopang_kit/crawling"
)

type ProtoRecordDividend struct {
	DividendID   int
	DividendRate float64
	HomeScore    string
	AwayScore    string
	Hit          bool
}

func (p ProtoRecordDividend) Json() string {
	b, err := json.Marshal(&p)
	if err != nil {
		return ""
	}
	return string(b)
}

type ProtoRecordMatch struct {
	Url       string
	RoundNo   int
	GameType  string
	MatchTime time.Time
	//MatchStr   string
	GameTitle  string
	GameResult string
}

type ProtoRecordMatchResult struct {
	match     *ProtoRecordMatch
	dividends []*ProtoRecordDividend
}

func (p ProtoRecordMatch) Json() string {
	b, err := json.Marshal(&p)
	if err != nil {
		return ""
	}
	return string(b)
}

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

//func toChar(i int) rune {
//	return rune('A' - 1 + i)
//}

func LetterToNum(letter string) int {
	for i, a := range letters {
		if letter == string(a) {
			return i + 1
		}
	}
	return 0
}

func NumToLetter(i int) string {
	return letters[i-1 : i]
}

/**
 * 기록식 적중결과 리스트
 * @param integer $nPage
 */
func GetRecordGameHitResultListUrl(page int) (string, error) {
	resultPage, err := GetHitResulPage(ListKeyPageMap, GameTypeRecord)
	if err != nil {
		return "", err
	}
	queryString := NewHitResultListParam(GameIDRecord, strconv.Itoa(page)).BuildQuery()
	return BaseUrl + "/" + resultPage + "?" + queryString, nil
}

/**
// wanted
						<tr>
							<td>1</td>
							<td>H1~3</td>
							<td class="pnt">
								<span>6.80</span>
								<img src="images/icon/ico_chkoff.gif" alt="" />
							</td>
						</tr>
*/
func getRecordGameHitResultDividend(tr *colly.HTMLElement, tdCount int) (dividends []*ProtoRecordDividend) {
	tr.ForEach("div.viwRap div.lst tbody tr", func(n int, tr *colly.HTMLElement) {
		var (
			dividend ProtoRecordDividend
			count    int
		)
		tr.ForEach("td", func(k int, td *colly.HTMLElement) {
			text := strings.TrimSpace(td.Text)
			switch k {
			case 0: // dividend id
				if dividendID, err := strconv.Atoi(text); err == nil {
					dividend.DividendID = dividendID
				}
			case 1: // score
				if strings.Contains(text, "-") {
					score := strings.Split(text, "-")
					dividend.HomeScore = score[0]
					dividend.AwayScore = score[1]
				} else {
					dividend.HomeScore = text
					dividend.AwayScore = ""
				}

			case 2: // hit
				if dividendRate, err := strconv.ParseFloat(td.ChildText("span"), 64); err == nil {
					dividend.DividendRate = dividendRate
				}
				imgSrc := td.ChildAttr("img", "src")
				if strings.Contains(imgSrc, "ico_chkon") {
					dividend.Hit = true
				}
			}
			count = k + 1
		})
		if tdCount == count && dividend.DividendID > 0 && dividend.DividendRate > 0 {
			dividends = append(dividends, &dividend)
		}
	})
	return
}

/**
// wanted
		<tr>
			<td>A</td>
			<td><img src="images/icon/ico_item_basketball.gif" alt="농구" /></td>
			<td>18-02-26 (월)</td>
			<td>19:00</td>
			<td class="sbj">
				<!--
				<a href="javascript:showBuyableGameDetail('A','winResultProtoRecordDetail.html');">
				프리미어리그 아스널 - 애스턴빌라 최종점수는?
				</a>
				 -->
				 WKBL 삼성생명-신한은행 최종점수차
			</td>
			<td class="bgRes01">5번. (H13~15) 13.2배</td>
		</tr>

// not wanted
						<tr>
							<td>1</td>
							<td>H1~3</td>
							<td class="pnt">
								<span>6.80</span>
								<img src="images/icon/ico_chkoff.gif" alt="" />
							</td>
						</tr>
*/
func getRecordGameHitResultMatch(tr *colly.HTMLElement, tdCount int) *ProtoRecordMatch {
	var (
		match ProtoRecordMatch
		count int
	)
	tr.ForEach("td", func(n int, td *colly.HTMLElement) {
		text := strings.TrimSpace(td.Text)
		switch n {
		case 0: // 게임 A~Z
			//match.RoundNo = strconv.Itoa(LetterToNum(text))
			match.RoundNo = LetterToNum(text)
		case 1: // 종목
			alt := td.ChildAttr("img", "alt")
			if val, ok := GameType[alt]; ok {
				match.GameType = val
			}
		case 2: // 경기일, Sun Feb 11 23:15:00 KST 2018
			if matchTime, err := time.Parse("Mon Jan 02 15:04:05 MST 2006", text); err == nil {
				match.MatchTime = matchTime
			}
		case 3: // 경기시각
		case 4: // 게임주제
			match.GameTitle = text
		case 5: // 프로토결과
			match.GameResult = text
		}
		count = n + 1
	})
	if tdCount == count && match.GameType != "" && match.MatchTime.String() != "" {
		return &match
	}
	return nil
}

func GetRecordGameHitResult(page int) {
	listUrl, err := GetRecordGameHitResultListUrl(page)
	if err != nil {
		panic(err)
	}
	fmt.Println(listUrl)

	req := crawling.NewRequest()
	req.SetAgentByOs(crawling.OsWindows)

	c := colly.NewCollector(
		colly.UserAgent(req.UserAgent),
	)
	detailCollector := c.Clone()

	c.OnHTML("#contents tbody tr td a", func(a *colly.HTMLElement) {
		if resultPage, err := GetHitResulPage(DetailKeyPageMap, GameTypeRecord); err == nil {
			detailUrl := strings.Join([]string{
				BaseUrl,
				"/",
				resultPage,
				"?",
				NewURLParamFromURL(a.Attr("href")).BuildQuery(),
			}, "")
			detailCollector.Visit(detailUrl)
		}
	})

	var results []*ProtoRecordMatchResult
	detailCollector.OnHTML("#tblSort tbody", func(d *colly.HTMLElement) {
		var result ProtoRecordMatchResult
		d.ForEach("tr", func(n int, tr *colly.HTMLElement) {
			switch n {
			case 0:
				if match := getRecordGameHitResultMatch(tr, 6); match != nil {
					result.match = match
					result.match.Url = d.Request.URL.String()
				}
			case 1:
				result.dividends = getRecordGameHitResultDividend(tr, 3)
			}
		})
		if result.match != nil {
			results = append(results, &result)
		}
	})
	c.Visit(listUrl)

	for i, _ := range results {
		fmt.Println("match:", results[i].match.Json())
		for _, dividend := range results[i].dividends {
			fmt.Println("dividends:", dividend.Json())
		}
		fmt.Println("\n\n")
	}
}
