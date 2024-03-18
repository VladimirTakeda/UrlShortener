package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/handler"
)

func (s *TestSuite) TestShortOk() {
	code, rawBody := s.doRequest(s.shortRequest(newSimpleShortRequest("http://ya.ru")))
	s.Require().Equal(http.StatusOK, code)

	response, err := s.shortResponseFromBody(rawBody)
	s.Require().NoError(err)

	shortUrl, err := url.Parse(response.ShortUrl)
	s.Require().NoError(err)
	s.Require().Len(shortUrl.Path, 18)

	s.Require().Len(response.SecretKey, 16)
}

func (s *TestSuite) TestShortInvalidLongUrl() {
	// url without schema
	code, _ := s.doRequest(s.shortRequest(newSimpleShortRequest("ya.ru")))
	s.Require().Equal(http.StatusBadRequest, code)
}

func (s *TestSuite) TestShortDuplicate() {
	// main request
	code, rawBody := s.doRequest(s.shortRequest(newSimpleShortRequest("http://ya.ru")))
	s.Require().Equal(http.StatusOK, code)

	response, err := s.shortResponseFromBody(rawBody)
	s.Require().NoError(err)

	shortUrl, err := url.Parse(response.ShortUrl)
	s.Require().NoError(err)
	s.Require().Len(shortUrl.Path, 18)

	s.Require().Len(response.SecretKey, 16)

	// duplicated request
	code, rawBody = s.doRequest(s.shortRequest(newSimpleShortRequest("http://ya.ru")))
	s.Require().Equal(http.StatusOK, code)

	response, err = s.shortResponseFromBody(rawBody)
	s.Require().NoError(err)

	shortUrl, err = url.Parse(response.ShortUrl)
	s.Require().NoError(err)
	s.Require().Len(shortUrl.Path, 18)

	s.Require().Len(response.SecretKey, 0)
}

func (s *TestSuite) TestShortVIPOk() {
	code, rawBody := s.doRequest(s.shortRequest(newFullShortRequest("http://ya.ru", "my-vip-link", 0, "")))
	s.Require().Equal(http.StatusOK, code)

	response, err := s.shortResponseFromBody(rawBody)
	s.Require().NoError(err)

	shortUrl, err := url.Parse(response.ShortUrl)
	s.Require().NoError(err)
	s.Require().Equal("/api/v1/my-vip-link", shortUrl.Path)

	s.Require().Len(response.SecretKey, 16)
}

func (s *TestSuite) TestShortVIPDuplicate() {
	// main request
	code, rawBody := s.doRequest(s.shortRequest(newFullShortRequest("http://ya.ru", "my-vip-link", 0, "")))
	s.Require().Equal(http.StatusOK, code)

	response, err := s.shortResponseFromBody(rawBody)
	s.Require().NoError(err)

	shortUrl, err := url.Parse(response.ShortUrl)
	s.Require().NoError(err)
	s.Require().Equal("/api/v1/my-vip-link", shortUrl.Path)

	s.Require().Len(response.SecretKey, 16)

	// duplicated request
	code, _ = s.doRequest(s.shortRequest(newFullShortRequest("http://ya.ru", "my-vip-link", 0, "")))
	s.Require().Equal(http.StatusBadRequest, code)
}

func (s *TestSuite) TestShortVIPTTL() {
	code, rawBody := s.doRequest(s.shortRequest(newFullShortRequest("http://ya.ru", "my-vip-link", 48, "")))
	s.Require().Equal(http.StatusOK, code)

	response, err := s.shortResponseFromBody(rawBody)
	s.Require().NoError(err)

	shortUrl, err := url.Parse(response.ShortUrl)
	s.Require().NoError(err)
	s.Require().Equal("/api/v1/my-vip-link", shortUrl.Path)

	s.Require().Len(response.SecretKey, 16)
}

func (s *TestSuite) TestShortVIPExceedingTTL() {
	// 49 hours > 2 days == max ttl
	code, _ := s.doRequest(s.shortRequest(newFullShortRequest("http://ya.ru", "my-vip-link", 49, "")))
	s.Require().Equal(http.StatusBadRequest, code)
}

func (s *TestSuite) TestShortVIPCheckTTL() {
	// add link with ttl
	code, rawBody := s.doRequest(s.shortRequest(newFullShortRequest("http://ya.ru", "my-vip-link", 2, "SECONDS")))
	s.Require().Equal(http.StatusOK, code)

	response, err := s.shortResponseFromBody(rawBody)
	s.Require().NoError(err)

	shortUrl, err := url.Parse(response.ShortUrl)
	s.Require().NoError(err)
	s.Require().Equal("/api/v1/my-vip-link", shortUrl.Path)

	s.Require().Len(response.SecretKey, 16)

	// check redirect
	code, _ = s.doRequest(s.redirectRequest("/api/v1/my-vip-link"))
	s.Require().Equal(http.StatusTemporaryRedirect, code)

	time.Sleep(3 * time.Second)

	// check redirect after ttl expired
	code, _ = s.doRequest(s.redirectRequest("/api/v1/my-vip-link"))
	s.Require().Equal(http.StatusNotFound, code)
}

type shortRequest struct {
	LongUrl string
	VipKey  string
	TTL     int
	TTLUnit string
}

func newSimpleShortRequest(longUrl string) shortRequest {
	return shortRequest{
		LongUrl: longUrl,
	}
}

func newFullShortRequest(longUrl, vipKey string, ttl int, ttlUnit string) shortRequest {
	return shortRequest{
		LongUrl: longUrl,
		VipKey:  vipKey,
		TTL:     ttl,
		TTLUnit: ttlUnit,
	}
}

var shortRequestTemplate = template.Must(template.New("srt").Parse(`
{
	"long_url": "{{ .LongUrl }}"
{{- if .VipKey}}
	, "vip_key": "{{ .VipKey }}"
{{- end}}
{{- if .TTL}}
	, "ttl": {{ .TTL }}
{{- end}}
{{- if .TTLUnit}}
	, "ttl_unit": "{{ .TTLUnit }}"
{{- end}}
}`))

func (s *TestSuite) shortRequest(r shortRequest) *http.Request {
	var buffer bytes.Buffer
	err := shortRequestTemplate.Execute(&buffer, r)
	s.Require().NoError(err)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/make_shorter", &buffer)
	return req
}

func (s *TestSuite) shortResponseFromBody(body string) (handler.ShortLinkResponse, error) {
	var res handler.ShortLinkResponse
	err := json.Unmarshal([]byte(body), &res)
	return res, err
}
