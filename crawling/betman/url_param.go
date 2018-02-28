package betman

import (
	"net/url"
	"reflect"
	"strings"
)

// TODO struct private 으로 변경
// new func
// 아래 func 도 new로 변경
// TODO UrlParam 을 정형화되지 않은 다양한 형태의 struct 를 변환할수 있도록
type UrlParam struct {
	Method         string `urlparam:"method,omitempty"`
	ViewMethod     string `urlparam:"viewmethod,omitempty"`
	GameId         string `urlparam:"gameId,omitempty"`
	GameIds        string `urlparam:"gameIds,omitempty"`
	GameRound      string `urlparam:"gameRound,omitempty"`
	League         string `urlparam:"league,omitempty"`
	Team           string `urlparam:"team,omitempty"`
	SelectedLeague string `urlparam:"selectedLeague,omitempty"`
	YearMonth      string `urlparam:"yearMonth,omitempty"`
	OuterRound     string `urlparam:"outerRound,omitempty"`
	SaleYear       string `urlparam:"saleYear,omitempty"`
	SelectedGameId string `urlparam:"selectedGameId,omitempty"`
	Page           string `urlparam:"page,omitempty"`
}

func newUrlParam() *UrlParam {
	return &UrlParam{}
}

// structTagKey will return struct tag key
func (p UrlParam) structTagKey() string {
	return "urlparam"
}

// BuildQuery will generate query string
func (p *UrlParam) BuildQuery() string {
	var (
		i        int
		query    string
		queryMap = make(map[string]string)
	)

	ptrT := reflect.TypeOf(p).Elem()
	ptrV := reflect.ValueOf(p).Elem()
	for i := 0; i < ptrV.NumField(); i++ {
		tag := ptrT.Field(i).Tag.Get(p.structTagKey())
		name, opts := parseTag(tag)
		if name == "" {
			name = ptrT.Field(i).Name
		}
		val := ptrV.Field(i).String()
		if strings.Contains(opts, "omitempty") && val == "" {
			continue
		}
		queryMap[name] = url.QueryEscape(val)
	}
	for name, val := range queryMap {
		var connector string
		switch i {
		case 0: // nothing to do
		default:
			connector = "&"
		}
		query += connector + name + "=" + val
		i++
	}
	return query
}

// SetByTagName returns the key of struct tag
func (p *UrlParam) SetByTagName(key, val string) {
	ptrT := reflect.TypeOf(p).Elem()
	ptrV := reflect.ValueOf(p).Elem()
	for i := 0; i < ptrV.NumField(); i++ {
		tag := ptrT.Field(i).Tag.Get(p.structTagKey())
		name, _ := parseTag(tag)
		if name == "" {
			name = ptrT.Field(i).Name
		}
		if name == key {
			ptrV.Field(i).SetString(val)
			break
		}
	}
}

func NewUrlParamFromQuery(query string) *UrlParam {
	urlParam := newUrlParam()
	var seperator = "?"
	if strings.Contains(query, seperator) {
		queries := strings.Split(query, seperator)
		query = queries[1]
	}
	for _, str := range strings.Split(query, "&") {
		var key, val string
		for i, q := range strings.Split(str, "=") {
			switch i {
			case 0:
				key = q
			case 1:
				val = q
			}
		}
		urlParam.SetByTagName(key, val)
	}
	return urlParam
}

func parseTag(tag string) (string, string) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tag[idx+1:]
	}
	return tag, ""
}
