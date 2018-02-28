package betman

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"encoding/json"

	"github.com/andybalholm/cascadia"
	"github.com/sangchul-sim/godom"
	"golang.org/x/net/html"
)

/**
<div id="contents">
	<table class="dataH01">
		<tbody>
			<tr>
				<td><img src=...></td>
				<td class="gname">
					<a href="winningResultProto.so?method=detail&amp;gameId=G102&amp;gameRound=180356&amp;page=2&amp;selectedGameId=G102" class="lkBlue">
										프로토 기록식 5회차
										 (R)
								</a>
				</td>
			</tr>
		</tbody>
	</table>

*/
func getRecordGameHitResultDetailParam(b []byte) (details []*UrlParam) {
	doc, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	s, err := cascadia.Compile("#contents tbody tr td a")
	if err != nil {
		panic(fmt.Errorf("error compiling %s", err))
	}
	for _, a := range s.MatchAll(doc) {
		attr := godom.NewGoQuery(a).GetAttributeByKey("href")
		if attr.Val != "" {
			details = append(details, NewUrlParamFromQuery(attr.Val))
		}
	}
	return
}

type ProtoRecordMatchDividend struct {
	DividendID   string
	DividendRate string
	HomeScore    string
	AwayScore    string
	IsHit        bool
}

type ProtoRecordMatch struct {
	RoundNo    string
	GameType   string
	MatchTime  time.Time
	MatchStr   string
	GameTitle  string
	GameResult string
}

func (p ProtoRecordMatch) Json() string {
	b, err := json.Marshal(&p)
	if err != nil {
		return ""
	}
	return string(b)
}

func getRecordGameHitResultDetailByDividend(tr *html.Node) (dividends []*ProtoRecordMatchDividend) {
	s, err := cascadia.Compile("div.viwRap div.lst tbody tr")
	if err != nil {
		panic(err)
	}

	for _, tr := range s.MatchAll(tr) {
		var (
			homeAwayStr string
			homeAway    []string
			dividend    ProtoRecordMatchDividend
		)
		for k, td := range godom.NewGoQuery(tr).GetElementsByTagName("td") {
			gq := godom.NewGoQuery(td)
			text := strings.TrimSpace(gq.GetInnerText())
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
				span := gq.GetElementsByTagName("span")[0]
				dividend.DividendRate = strings.TrimSpace(godom.NewGoQuery(span).GetInnerText())
				img := gq.GetElementsByTagName("img")[0]
				attr := godom.NewGoQuery(img).GetAttributeByKey("src")
				if strings.Contains(attr.Val, "ico_chkon") {
					dividend.IsHit = true
				}
			}
		}
		dividends = append(dividends, &dividend)
	}

	return
}

func getRecordGameHitResultDetailByMatch(tr *html.Node) *ProtoRecordMatch {
	var match ProtoRecordMatch
	for j, td := range godom.NewGoQuery(tr).GetElementsByTagName("td") {
		text := strings.TrimSpace(godom.NewGoQuery(td).GetInnerText())
		switch j {
		case 0: // 게임 A~Z
			match.RoundNo = text
		case 1: // 종목
			attr := godom.NewGoQuery(td.FirstChild).GetAttributeByKey("alt")
			if val, ok := GameType[attr.Val]; ok {
				match.GameType = val
			}
		case 2: // 경기일, Sun Feb 11 23:15:00 KST 2018
			match.MatchTime, _ = time.Parse("Mon Jan 02 15:04:05 MST 2006", td.FirstChild.Data)
		case 3: // 경기시각
		case 4: // 게임주제
			match.GameTitle = text
		case 5: // 프로토결과
			match.GameResult = text
		}
	}
	return &match
}

func getRecordGameHitResultDetail(b []byte) {
	doc, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		panic(err)
	}

	var match *ProtoRecordMatch
	var dividends []*ProtoRecordMatchDividend
	s, err := cascadia.Compile("#tblSort tbody tr")
	if err != nil {
		panic(err)
	}
	for i, tr := range s.MatchAll(doc) {
		switch i {
		case 0:
			match = getRecordGameHitResultDetailByMatch(tr)
		case 1:
			dividends = getRecordGameHitResultDetailByDividend(tr)
		}
	}

	fmt.Println("match:", match)
	for _, dividend := range dividends {
		fmt.Println("dividend:", dividend)
	}
}

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

var letterArr = [...]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
	"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

//func toChar(i int) rune {
//	return rune('A' - 1 + i)
//}

// TODO letterArr 대신 letters 를 이용하는 방법 강구할 것
func LetterToNum(letter string) int {
	for i, val := range letterArr {
		if letter == val {
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
