package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/google/uuid"
)

type API struct {
	VERSION      string
	KeyVersion   string
	KEY          string
	CAPABILITIES string
}

type HttpResponse struct {
	Err       error
	ResStatus int
	Req       *http.Request
	Res       *http.Response
	Body      string
	Headers   http.Header
	Cookies   *cookiejar.Jar
}

func MakeHttpResponse(Response *http.Response, Request *http.Request, jar *cookiejar.Jar, Error error) HttpResponse {

	var res = ""
	var StatusCode = 0
	var Headers http.Header = nil

	if Response != nil {
		var reader io.ReadCloser
		switch Response.Header.Get("Content-Encoding") {
		case "gzip":
			reader, _ = gzip.NewReader(Response.Body)
			defer reader.Close()
		default:
			reader = Response.Body
		}
		body, _ := ioutil.ReadAll(reader)
		res = string(body)

		if Response.Header != nil {
			Headers = Response.Header
		}

		if Response.StatusCode != 0 {
			StatusCode = Response.StatusCode
		}
	}

	return HttpResponse{ResStatus: StatusCode, Res: Response, Req: Request, Body: res, Headers: Headers, Cookies: jar, Err: Error}
}

func createKeyValuePairs(m http.Header) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		_, _ = fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

func HMACSHA256(message string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func GetAPI() API {

	IG_VERSION := "10.26.0"
	IG_SIG_KEY := "4f8732eb9ba7d1c8e8897a75d6474d4eb3f5279137431b2aafb71fafe2abe178"
	SIG_KEY_VERSION := "4"
	X_IG_Capabilities := "3brTvw=="

	_API := API{VERSION: IG_VERSION, KEY: IG_SIG_KEY, KeyVersion: SIG_KEY_VERSION, CAPABILITIES: X_IG_Capabilities}

	return _API
}

func IR(iurl string, signedbody map[string]string, payload string,
	Headers map[string]string, api API, proxy string,
	ptype string, cookie *cookiejar.Jar, usecookies bool) HttpResponse {

	_url := iurl

	if ((!strings.Contains(_url, "https")) || (!strings.Contains(_url, "http"))) && _url[0] != '/' {
		_url = "https://i.instagram.com/api/v1/" + _url
	} else if ((!strings.Contains(_url, "https")) || (!strings.Contains(_url, "http"))) && _url[0] == '/' {
		_url = "https://i.instagram.com/api/v1" + _url
	}

	_api := API{}
	if api == (API{}) {
		_api = GetAPI()
	} else {
		_api = api
	}

	_payload := ""
	if signedbody != nil {
		_data, _ := json.Marshal(signedbody)
		_json := string(_data)
		_signed := fmt.Sprintf("%v.%s", HMACSHA256(_api.KEY, _json), _json)
		_payload = "ig_sig_key_version=" + _api.KeyVersion + "&signed_body=" + _signed
	} else if payload != "" {
		_payload = payload
	}

	var req *http.Request
	if _payload != "" {
		req, _ = http.NewRequest("POST", _url, bytes.NewBuffer([]byte(_payload)))
	} else {
		req, _ = http.NewRequest("GET", _url, nil)
	}

	req.Header.Set("User-Agent", "Instagram "+_api.VERSION+" Android (19/4.4.2; 480dpi; 1080x1920; samsung; SM-N900T; hltetmo; qcom; en_US)")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie2", "$Version=1")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-IG-Connection-Type", "WIFI")
	req.Header.Set("X-IG-Capabilities", _api.CAPABILITIES)
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "Keep-Alive")

	if Headers != nil {
		var keys []string
		for key := range Headers {
			keys = append(keys, key)
		}
		var values []string
		for _, value := range Headers {
			values = append(values, value)
		}

		for i := 0; i < len(keys); i++ {
			req.Header.Set(keys[i], values[i])
		}
	}

	jar := cookie
	transport := http.Transport{}
	if proxy != "" {
		proxyUrl, _ := url.Parse(ptype + "://" + proxy)
		transport.Proxy = http.ProxyURL(proxyUrl) // set proxy proxyType://proxyIp:proxyPort
	}
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //set ssl
	client := &http.Client{}
	if usecookies {
		client = &http.Client{Jar: jar}
	}
	client.Transport = &transport
	resp, err := client.Do(req)
	if err != nil {
		return MakeHttpResponse(resp, req, jar, err)
	}
	defer resp.Body.Close()
	return MakeHttpResponse(resp, req, jar, err)
}

func MakeList(chars []string, l int) []string {
	var list []string
	var clearList []string
	var n = len(chars)
	ml(chars, "", n, l, &list)
	for _, v := range list {
		if v[:1] == "." || v[(len(v)-1):] == "." {
		} else {
			clearList = append(clearList, v)
		}
	}
	return clearList
}

func ml(chars []string, prefix string, n int, l int, list *[]string) {
	var copied []string
	if l == 0 {
		copied = *list
		copied = append(copied, prefix)
		*list = copied
		return
	}
	for i := 0; i < n; i++ {
		newPrefix := prefix + chars[i]
		ml(chars, newPrefix, n, l-1, list)
	}
}

func CreateUsernames(chars []string, length int) []string {
	t := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "_", "."}
	l := 3
	if length != 0 {
		l = length
	}
	if chars != nil {
		t = chars
	}
	return MakeList(t, l)
}

func StringWithCharset(length int, charset string) string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	var letters = []rune(charset)
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[seededRand.Intn(len(letters))]
	}
	return string(b)
}

func unique(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// GetProxies function
func GetProxies() ([]string, []string, []string) {

	var req *http.Request
	req, _ = http.NewRequest("GET", "https://raw.githubusercontent.com/fate0/proxylist/master/proxy.list", nil)
	req.Header.Set("Host", "raw.githubusercontent.com")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Accept-Language", "ar,en-US;q=0.7,en;q=0.3")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:67.0) Gecko/20100101 Firefox/67.0")
	req.Header.Set("Connection", "keep-alive")

	transport := http.Transport{}
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //set ssl
	client := &http.Client{}
	client.Transport = &transport
	resp, err := client.Do(req)
	_ = err

	var reader io.ReadCloser
	var response = ""

	if resp != nil {
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, _ = gzip.NewReader(resp.Body)
			defer reader.Close()
		default:
			reader = resp.Body
		}
		body, _ := ioutil.ReadAll(reader)
		response = string(body)
	}

	_Proxies := strings.Split(response, "}\n{")
	var HTTPSProxies []string
	var HTTPProxies []string

	for i := 0; i < len(_Proxies); i++ {
		if i == 0 {

			current := _Proxies[i] + "}"
			_proxy := make(map[string]interface{})
			json.Unmarshal([]byte(current), &_proxy)

			_type := fmt.Sprintf("%v", _proxy["type"])
			_ip := fmt.Sprintf("%v", _proxy["host"])
			_port := fmt.Sprintf("%v", _proxy["port"])

			if _type == "https" {
				HTTPSProxies = append(HTTPSProxies, _type+"://"+_ip+":"+_port)
			}
			if _type == "http" {
				HTTPProxies = append(HTTPProxies, _type+"://"+_ip+":"+_port)
			}

			continue
		}
		if i == len(_Proxies)-1 {

			current := "{" + _Proxies[i]
			_proxy := make(map[string]interface{})
			json.Unmarshal([]byte(current), &_proxy)

			_type := fmt.Sprintf("%v", _proxy["type"])
			_ip := fmt.Sprintf("%v", _proxy["host"])
			_port := fmt.Sprintf("%v", _proxy["port"])

			if _type == "https" {
				HTTPSProxies = append(HTTPSProxies, _type+"://"+_ip+":"+_port)
			}
			if _type == "http" {
				HTTPProxies = append(HTTPProxies, _type+"://"+_ip+":"+_port)
			}

			break
		}

		current := "{" + _Proxies[i] + "}"
		_proxy := make(map[string]interface{})
		json.Unmarshal([]byte(current), &_proxy)

		_type := fmt.Sprintf("%v", _proxy["type"])
		_ip := fmt.Sprintf("%v", _proxy["host"])
		_port := fmt.Sprintf("%v", _proxy["port"])

		if _type == "https" {
			HTTPSProxies = append(HTTPSProxies, _type+"://"+_ip+":"+_port)
		}
		if _type == "http" {
			HTTPProxies = append(HTTPProxies, _type+"://"+_ip+":"+_port)
		}

	}
	return HTTPSProxies, HTTPProxies, append(HTTPProxies, HTTPSProxies...)
}

func ssliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func login(us string, ps string, proxy string, proxy_type string, FakeCookies bool, InstaAPI API) HttpResponse {
	url := "https://i.instagram.com/api/v1/accounts/login/"

	jar, _ := cookiejar.New(nil)

	u, _ := uuid.NewUUID()
	guid := u.String()

	post := make(map[string]string)
	post["phone_id"] = guid
	post["_csrftoken"] = "missing"
	post["username"] = us
	post["password"] = ps
	post["device_id"] = guid
	post["guid"] = guid
	post["login_attempt_count"] = "0"

	if FakeCookies {
		jar = APICreateCookies(-1, false, "", "")
		return IR(url, post, "", nil, InstaAPI, proxy, proxy_type, jar, true)
	}
	return IR(url, post, "", nil, InstaAPI, proxy, proxy_type, jar, true)
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// writeLines writes the lines to the given file.
func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		_, _ = fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func RandRange(min int, max int) int {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	_max := max + 1
	return seededRand.Intn(_max-min) + min
}

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

func b(bin int) bool {
	if bin == 0 {
		return false
	}
	if bin == 1 {
		return true
	}
	return false
}

func RandomChoice() bool {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	min := 0
	max := 2
	return b(seededRand.Intn(max-min) + min)
}

func GenerateRandomString(length int, numbers bool, uppers bool, lowers bool, symbols bool, NotStartWithNumber bool) string {
	var _uppers = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var _lowers = "abcdefghijklmnopqrstuvwxyz"
	var _numbers = "0123456789"
	var _symbols = "!@#$%^&*()_-+=\\/,.?<>|"
	var charset = ""
	if numbers {
		charset += charset + _numbers
	}
	if lowers {
		charset += charset + _lowers
	}
	if uppers {
		charset += charset + _uppers
	}
	if symbols {
		charset += charset + _symbols
	}
	if charset == "" {
		log.Println("Are an idoit ?")
		charset = "Yes, Im an Idoit"
	}
	rnd := StringWithCharset(length, charset)
	if NotStartWithNumber {
		matched, _ := regexp.MatchString("[0-9]", string(rnd[0]))
		if matched {
			chars := string(_lowers)
			if RandomChoice() {
				chars = string(_uppers)
			}
			rnd = replaceAtIndex(rnd, rune(chars[RandRange(0, len(chars)-1)]), 0)
		}
	}
	return rnd
}

// APICreateCookies function
func APICreateCookies(URLType int, Legit bool, dsUserId string, dsUser string) *cookiejar.Jar {

	var mid string
	var csrftoken string
	var sessionid string
	var ds_user_id string = dsUserId
	var ds_user string = dsUser
	var rur string

	if Legit {
		CLegit := GetLegitCookies()
		_url, _ := url.Parse("https://www.instagram.com/")
		var Cookies = CLegit.Cookies(_url)
		for i := 0; i < len(Cookies); i++ {
			if strings.Contains(strings.ToLower(Cookies[i].Name), "csrftoken") {
				csrftoken = Cookies[i].Value
			}
			if strings.Contains(strings.ToLower(Cookies[i].Name), "mid") {
				mid = Cookies[i].Value
			}
		}
	} else {
		mid = "XSy" + GenerateRandomString(25, true, true, true, false, false) // was without XSx and was 27 and true start with char
		csrftoken = GenerateRandomString(32, true, true, true, false, true)
		NewChar := '_'
		if RandomChoice() {
			NewChar = '-'
		}
		if RandomChoice() {
			mid = replaceAtIndex(mid, NewChar, (RandRange(4, 27))) // was 1, 28
		}
	}

	Random_rur := RandRange(0, 4)

	if Random_rur == 0 {
		rur = "PRN"
	}
	if Random_rur == 1 {
		rur = "ASH"
	}
	if Random_rur == 2 {
		rur = "ATN"
	}
	if Random_rur == 3 {
		rur = "FRC"
	}
	if Random_rur == 4 {
		rur = "FTW"
	}

	if ds_user_id == "" {
		ds_user_id = GenerateRandomString(10, true, false, false, false, false)
	}

	sessionid = ds_user_id + "%3A" + GenerateRandomString(14, true, true, true, false, false) + "%3A" + GenerateRandomString(2, true, false, false, false, false)

	var cookies []*http.Cookie

	cookie0 := &http.Cookie{
		Name:     "sessionid",
		Value:    sessionid,
		Path:     "/",
		Domain:   "instagram.com",
		Secure:   true,
		HttpOnly: true,
	}
	cookies = append(cookies, cookie0)

	cookie1 := &http.Cookie{
		Name:   "mid",
		Value:  mid,
		Path:   "/",
		Domain: "instagram.com",
		Secure: true,
	}
	cookies = append(cookies, cookie1)

	cookie2 := &http.Cookie{
		Name:   "csrftoken",
		Value:  csrftoken,
		Path:   "/",
		Domain: "instagram.com",
		Secure: true,
	}
	cookies = append(cookies, cookie2)

	cookie4 := &http.Cookie{
		Name:   "ds_user_id",
		Value:  ds_user_id,
		Path:   "/",
		Domain: "instagram.com",
		Secure: true,
	}
	cookies = append(cookies, cookie4)

	cookie3 := &http.Cookie{
		Name:     "rur",
		Value:    rur,
		Path:     "/",
		Domain:   "instagram.com",
		Secure:   true,
		HttpOnly: true,
	}
	cookies = append(cookies, cookie3)

	cookie5 := &http.Cookie{
		Name:     "shbts",
		Value:    "1" + GenerateRandomString(9, true, false, false, false, false) + "." + GenerateRandomString(7, true, false, false, false, false),
		Path:     "/",
		Domain:   "instagram.com",
		Secure:   true,
		HttpOnly: true,
	}
	cookies = append(cookies, cookie5)

	if ds_user == "" {
		ds_user = GenerateRandomString(RandRange(3, 10), true, false, true, false, true)
	}

	cookie6 := &http.Cookie{
		Name:     "ds_user",
		Value:    ds_user,
		Path:     "/",
		Domain:   "instagram.com",
		Secure:   true,
		HttpOnly: true,
	}
	cookies = append(cookies, cookie6)

	cookie7 := &http.Cookie{
		Name:     "shbid",
		Value:    "11" + GenerateRandomString(3, true, false, false, false, false),
		Path:     "/",
		Domain:   "instagram.com",
		Secure:   true,
		HttpOnly: true,
	}
	cookies = append(cookies, cookie7)

	u, _ := url.Parse("https://www.instagram.com/")

	jar, _ := cookiejar.New(nil)
	jar.SetCookies(u, cookies)

	return jar
}

// CreateCookies function
func CreateCookies(URLType int, Legit bool, dsUserId string) http.CookieJar {

	var mid string
	var csrftoken string
	var urlgen string
	var sessionid string
	var ds_user_id string = dsUserId
	var rur string

	if Legit {
		CLegit := GetLegitCookies()
		_url, _ := url.Parse("https://www.instagram.com/")
		var Cookies = CLegit.Cookies(_url)
		for i := 0; i < len(Cookies); i++ {
			if strings.Contains(strings.ToLower(Cookies[i].Name), "csrftoken") {
				csrftoken = Cookies[i].Value
			}
			if strings.Contains(strings.ToLower(Cookies[i].Name), "mid") {
				mid = Cookies[i].Value
			}
		}
	} else {
		mid = "XShb" + GenerateRandomString(24, b(1), b(1), true, false, true) // was without XShy and was 27
		csrftoken = GenerateRandomString(32, b(1), b(1), true, false, true)
		NewChar := '_'
		if RandomChoice() {
			NewChar = '-'
		}
		if RandomChoice() {
			mid = replaceAtIndex(mid, NewChar, (RandRange(4, 27))) // was 1, 28
		}
	}

	Random_rur := RandRange(0, 4)

	if Random_rur == 0 {
		rur = "PRN"
	}
	if Random_rur == 1 {
		rur = "ASH"
	}
	if Random_rur == 2 {
		rur = "ATN"
	}
	if Random_rur == 3 {
		rur = "FRC"
	}
	if Random_rur == 4 {
		rur = "FTW"
	}

	if ds_user_id == "" {
		ds_user_id = GenerateRandomString(10, true, false, false, false, false)
	}

	sessionid = ds_user_id + "%3A" + GenerateRandomString(14, true, true, true, false, false) + "%3A" + GenerateRandomString(2, true, false, false, false, false)

	begin, _ := url.QueryUnescape("%22%7B%5C%22")
	end, _ := url.QueryUnescape("%5C%22%3A%20")
	final := end + "25019}:1hl"                                        // was 1h
	ffinal := "}:1hl"                                                  // was 1h
	lastone := GenerateRandomString(3, true, true, true, false, false) // was 4
	eend, _ := url.QueryUnescape("25019%5C054%20%5C%22")
	theend := GenerateRandomString(27, true, true, true, false, false)
	comma, _ := url.QueryUnescape("%22")
	cr := '_'

	if RandomChoice() {
		cr = '-'
	}
	if RandomChoice() {
		theend = replaceAtIndex(theend, cr, (RandRange(1, 26)))
	}

	Random_urlgen := RandRange(0, 4)
	if URLType != -1 {
		Random_urlgen = URLType
	}

	if Random_urlgen == 0 {
		urlgen = fmt.Sprintf("%s%s%s%s%s%s%s:%s%s", begin, randomdata.IpV4Address(), end, eend, randomdata.IpV6Address(), final, lastone, theend, comma)
	}
	if Random_urlgen == 1 {
		urlgen = fmt.Sprintf("%s%s%s25019%s%s:%s%s", begin, randomdata.IpV4Address(), end, ffinal, lastone, theend, comma)
	}
	if Random_urlgen == 2 {
		urlgen = fmt.Sprintf("%s%s%s25019%s%s:%s%s", begin, randomdata.IpV6Address(), end, ffinal, lastone, theend, comma)
	}
	if Random_urlgen == 3 {
		urlgen = fmt.Sprintf("%s%s%s%s%s%s%s:%s%s", begin, randomdata.IpV6Address(), end, eend, randomdata.IpV6Address(), final, lastone, theend, comma)
	}
	if Random_urlgen == 4 {
		urlgen = fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s:%s%s", begin, randomdata.IpV6Address(), end, eend, randomdata.IpV6Address(), end, eend, randomdata.IpV6Address(), final, lastone, theend, comma)
	}

	var jar http.CookieJar
	var cookies []*http.Cookie

	cookie0 := &http.Cookie{
		Name:     "sessionid",
		Value:    sessionid,
		Path:     "/",
		Domain:   ".instagram.com",
		Secure:   true,
		HttpOnly: true,
	}
	cookies = append(cookies, cookie0)

	cookie1 := &http.Cookie{
		Name:   "mid",
		Value:  mid,
		Path:   "/",
		Domain: ".instagram.com",
		Secure: true,
	}
	cookies = append(cookies, cookie1)

	cookie2 := &http.Cookie{
		Name:   "csrftoken",
		Value:  csrftoken,
		Path:   "/",
		Domain: ".instagram.com",
		Secure: true,
	}
	cookies = append(cookies, cookie2)

	cookie4 := &http.Cookie{
		Name:   "ds_user_id",
		Value:  ds_user_id,
		Path:   "/",
		Domain: ".instagram.com",
		Secure: true,
	}
	cookies = append(cookies, cookie4)

	cookie3 := &http.Cookie{
		Name:     "rur",
		Value:    rur,
		Path:     "/",
		Domain:   ".instagram.com",
		Secure:   true,
		HttpOnly: true,
	}
	cookies = append(cookies, cookie3)

	cookie5 := &http.Cookie{
		Name:     "urlgen",
		Value:    urlgen,
		Path:     "/",
		Domain:   ".instagram.com",
		Secure:   true,
		HttpOnly: true,
	}
	cookies = append(cookies, cookie5)

	u, _ := url.Parse("https://www.instagram.com/")

	jar.SetCookies(u, cookies)

	return jar
}

// GetLegitCookies function
func GetLegitCookies() http.CookieJar {
	jar := CreateCookies(-1, false, "")
	//jar, _ := cookiejar.New(nil)
	var req *http.Request
	req, _ = http.NewRequest("GET", "https://www.instagram.com/", nil)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")
	transport := http.Transport{}
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //set ssl
	client := &http.Client{Jar: jar}
	client.Transport = &transport
	client.Do(req)
	return jar
}

func GetProfile(jar cookiejar.Jar, api API) (map[string]string, HttpResponse) {
	res := IR("accounts/current_user/?edit=true", nil, "", nil, api, "", "", &jar, true)
	var profile = make(map[string]string)

	var username = ""
	_username := regexp.MustCompile("\"username\": \"(.*?)\",").FindStringSubmatch(res.Body)
	if _username != nil {
		username = _username[1]
	}
	var biography = ""
	_biography := regexp.MustCompile("\"biography\": \"(.*?)\",").FindStringSubmatch(res.Body)
	if _biography != nil {
		biography = _biography[1]
	}

	var fullName = ""
	_fullName := regexp.MustCompile("\"full_name\": \"(.*?)\",").FindStringSubmatch(res.Body)
	if _fullName != nil {
		fullName = _fullName[1]
	}

	var phoneNumber = ""
	_phoneNumber := regexp.MustCompile("\"phone_number\": \"(.*?)\",").FindStringSubmatch(res.Body)
	if _phoneNumber != nil {
		phoneNumber = _phoneNumber[1]
	}

	var email = ""
	_email := regexp.MustCompile("\"email\": \"(.*?)\"").FindStringSubmatch(res.Body)
	if _email != nil {
		email = _email[1]
	}
	var gender = ""
	_gender := regexp.MustCompile("\"gender\": \"(.*?)\",").FindStringSubmatch(res.Body)
	if _gender != nil {
		gender = _gender[1]
	}

	var externalUrl = ""
	_externalUrl := regexp.MustCompile("\"external_url\": \"(.*?)\",").FindStringSubmatch(res.Body)
	if _externalUrl != nil {
		externalUrl = _externalUrl[1]
	}

	var isVerified = ""
	_isVerified := regexp.MustCompile("\"is_verified\": \"(.*?)\",").FindStringSubmatch(res.Body)
	if _isVerified != nil {
		isVerified = _isVerified[1]
	}

	profile["username"] = username
	profile["biography"] = biography
	profile["full_name"] = fullName
	profile["phone_number"] = phoneNumber
	profile["email"] = email
	profile["gender"] = gender
	profile["external_url"] = externalUrl
	profile["is_verified"] = isVerified

	return profile, res
}
