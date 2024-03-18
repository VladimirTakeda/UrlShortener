package handler_test

import (
	"fmt"
	"net/http"
)

func (s *TestSuite) TestDeleteNotFound() {
	code, _ := s.doRequest(s.deleteRequest("random-bytes"))
	s.Require().Equal(http.StatusNotFound, code)
}

func (s *TestSuite) TestDeleteOk() {
	// create short url
	code, rawBody := s.doRequest(s.shortRequest(newSimpleShortRequest("http://ya.ru")))
	s.Require().Equal(http.StatusOK, code)

	shortResponse, err := s.shortResponseFromBody(rawBody)
	s.Require().NoError(err)

	// delete
	code, _ = s.doRequest(s.deleteRequest(shortResponse.SecretKey))
	s.Require().Equal(http.StatusOK, code)

	// check info
	code, _ = s.doRequest(s.infoRequest(shortResponse.SecretKey))
	s.Require().Equal(http.StatusNotFound, code)
}

func (s *TestSuite) deleteRequest(secretKey string) *http.Request {
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/admin/%s", secretKey), nil)
	return req
}
