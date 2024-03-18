package handler_test

import (
	"net/http"
)

func (s *TestSuite) TestHealthCheck() {
	code, _ := s.doRequest(s.healthCheckRequest())
	s.Require().Equal(http.StatusOK, code)
}

func (s *TestSuite) healthCheckRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/healthcheck/ping", nil)
	return req
}
