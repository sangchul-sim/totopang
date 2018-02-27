package crawling

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	RequestMethodGet    = "GET"
	RequestMethodPost   = "POST"
	RequestMethodPut    = "PUT"
	RequestMethodDelete = "DELETE"
	OsWindows           = "windows"
	OsMac               = "mac"
	OsAndroid           = "android"
)

var userAgents = map[string]string{
	OsWindows: "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.120 Safari/537.36",
	OsMac:     "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.120 Safari/537.36",
	OsAndroid: "Mozilla/5.0 (Linux; U; Android 4.1.2; ko-kr; SHV-E170K/KKJMK3 Build/JZO54K) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
}

// TODO method에 따라서 ContentType 변경할 것
type Request struct {
	Method    string
	UserAgent string
	//ContentType string
	Accept string
	Url    string
	form   url.Values
	req    *http.Request
	client *http.Client
}

func (r *Request) SetMethod(method string) *Request {
	r.Method = method
	return r
}

func (r *Request) SetAgentByOs(userOs string) *Request {
	if agent, ok := userAgents[userOs]; ok {
		r.UserAgent = agent
	}
	return r
}

//func (r *Request) SetContentType(contentType string) *Request {
//	r.ContentType = contentType
//	return r
//}

func (r *Request) SetUrl(url string) *Request {
	r.Url = url
	return r
}

func (r *Request) SetBody(data map[string]interface{}) *Request {
	for key, val := range data {
		r.form.Add(key, fmt.Sprintf("%v", val))
	}
	return r
}

func (r *Request) setHeader() {
	//r.req.Header.Add("Content-Type", r.ContentType)
	r.req.Header.Add("Accept", r.Accept)
	r.req.Header.Add("User-Agent", r.UserAgent)
	r.req.Header.Add("Referer", r.Url)
	r.req.PostForm = r.form
}

func (r *Request) setClient() {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	r.client = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
}

func (r *Request) setRequest() (err error) {
	r.req, err = http.NewRequest(r.Method, r.Url, strings.NewReader(r.form.Encode()))
	return
}

func (r *Request) Do() ([]byte, error) {
	if err := r.setRequest(); err != nil {
		return []byte{}, err
	}
	r.setClient()
	r.setHeader()

	resp, err := r.client.Do(r.req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("http status error: %d", resp.StatusCode)
	}

	//ret, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(ret))
	return ioutil.ReadAll(resp.Body)
}

func NewRequest() *Request {
	return &Request{
		//ContentType: "application/x-www-form-urlencoded",
		//ContentType: "text/html;charset=euc-kr",
		//Accept: "application/json",
		Accept: "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
	}
}

//func RequestUrl(reqUrl, method string, data map[string]interface{}) ([]byte, error) {
//	form := url.Values{}
//	for key, val := range data {
//		form.Add(key, fmt.Sprintf("%v", val))
//	}
//
//	req, err := http.NewRequest(method, reqUrl, strings.NewReader(form.Encode()))
//	var netTransport = &http.Transport{
//		Dial: (&net.Dialer{
//			Timeout: 5 * time.Second,
//		}).Dial,
//		TLSHandshakeTimeout: 5 * time.Second,
//	}
//	client := &http.Client{
//		Timeout:   time.Second * 10,
//		Transport: netTransport,
//	}
//	agent, _ := userAgent(OsWindows)
//	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
//	req.Header.Add("Accept", "application/json")
//	req.Header.Add("User-Agent", agent)
//	req.Header.Add("Referer", reqUrl)
//	req.PostForm = form
//
//	resp, err := client.Do(req)
//	if err != nil {
//		return []byte{}, err
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		return []byte{}, fmt.Errorf("http status error: %d", resp.StatusCode)
//	}
//
//	return ioutil.ReadAll(resp.Body)
//}

/**
$config['agents'] = array(
	'windows' => 'Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.120 Safari/537.36',
	'mac' => 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.120 Safari/537.36',
	'android' => 'Mozilla/5.0 (Linux; U; Android 4.1.2; ko-kr; SHV-E170K/KKJMK3 Build/JZO54K) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30',
);
*/

/**
"SERVER_SOFTWARE":  r.Header.Get("SERVER_SOFTWARE"),
"REQUEST_METHOD":   r.Method,
"HOST":             r.Host,
"X_REAL_IP":        r.Header.Get("X_REAL_IP"),
"X_FORWARDED_FOR":  r.Header.Get("X_FORWARDED_FOR"),
"CONNECTION":       r.Header.Get("CONNECTION"),
"USER-AGENT":       r.UserAgent(),
"ACCEPT":           r.Header.Get("ACCEPT"),
"ACCEPT-LANGUAGE":  r.Header.Get("ACCEPT-LANGUAGE"),
"ACCEPT-ENCODING":  r.Header.Get("ACCEPT-ENCODING"),
"X_REQUESTED_WITH": r.Header.Get("X_REQUESTED_WITH"),
"REFERER":          r.Referer(),
"PRAGMA":           r.Header.Get("PRAGMA"),
"CACHE_CONTROL":    r.Header.Get("CACHE_CONTROL"),
"REMOTE_ADDR":      r.RemoteAddr,
"REQUEST_TIME":     r.Header.Get("REQUEST_TIME"),
*/
