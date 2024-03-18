package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/handler"
)

func (s *TestSuite) TestInfoNotFound() {
	code, _ := s.doRequest(s.infoRequest("random-bytes"))
	s.Require().Equal(http.StatusNotFound, code)
}

func (s *TestSuite) TestInfoOk() {
	// create short url
	code, rawBody := s.doRequest(s.shortRequest(newSimpleShortRequest("http://ya.ru")))
	s.Require().Equal(http.StatusOK, code)

	shortResponse, err := s.shortResponseFromBody(rawBody)
	s.Require().NoError(err)

	// check info
	code, rawBody = s.doRequest(s.infoRequest(shortResponse.SecretKey))
	s.Require().Equal(http.StatusOK, code)

	infoResponse, err := s.infoResponseFromBody(rawBody)
	s.Require().NoError(err)

	s.Require().Equal(0, infoResponse.Clicks)
}

func (s *TestSuite) TestInfoOkWithRedirects() {
	// create short url
	code, rawBody := s.doRequest(s.shortRequest(newSimpleShortRequest("http://ya.ru")))
	s.Require().Equal(http.StatusOK, code)

	shortResponse, err := s.shortResponseFromBody(rawBody)
	s.Require().NoError(err)

	shortUrl, err := url.Parse(shortResponse.ShortUrl)
	s.Require().NoError(err)

	// make some redirects
	for i := 0; i < 15; i++ {
		code, _ = s.doRequest(s.redirectRequest(shortUrl.Path))
		s.Require().Equal(http.StatusTemporaryRedirect, code)
	}

	// check info
	code, rawBody = s.doRequest(s.infoRequest(shortResponse.SecretKey))
	s.Require().Equal(http.StatusOK, code)

	infoResponse, err := s.infoResponseFromBody(rawBody)
	s.Require().NoError(err)

	s.Require().Equal(15, infoResponse.Clicks)
}

func (s *TestSuite) infoRequest(secretKey string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/admin/%s", secretKey), nil)
	return req
}

func (s *TestSuite) infoResponseFromBody(body string) (handler.InfoResponse, error) {
	var res handler.InfoResponse
	err := json.Unmarshal([]byte(body), &res)
	return res, err
}
