package handler_test

import (
	"net/http"
	"net/url"
)

func (s *TestSuite) TestRedirectNotFound() {
	code, _ := s.doRequest(s.redirectRequest("/api/v1/random-bytes"))
	s.Require().Equal(http.StatusNotFound, code)
}

func (s *TestSuite) TestRedirectOk() {
	// create short url
	code, rawBody := s.doRequest(s.shortRequest(newSimpleShortRequest("http://ya.ru")))
	s.Require().Equal(http.StatusOK, code)

	shortResponse, err := s.shortResponseFromBody(rawBody)
	s.Require().NoError(err)

	shortUrl, err := url.Parse(shortResponse.ShortUrl)
	s.Require().NoError(err)

	// check redirect
	code, _ = s.doRequest(s.redirectRequest(shortUrl.Path))
	s.Require().Equal(http.StatusTemporaryRedirect, code)

	// check counters
	code, rawBody = s.doRequest(s.infoRequest(shortResponse.SecretKey))
	s.Require().Equal(http.StatusOK, code)

	infoResponse, err := s.infoResponseFromBody(rawBody)
	s.Require().NoError(err)

	s.Require().Equal(1, infoResponse.Clicks)
}

func (s *TestSuite) redirectRequest(shortPath string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, shortPath, nil)
	return req
}
