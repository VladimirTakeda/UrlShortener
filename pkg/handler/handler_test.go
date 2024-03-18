package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/stretchr/testify/suite"

	"git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/config"
	"git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/db"
	"git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/handler"
)

func TestHandlers(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

type TestSuite struct {
	suite.Suite

	ctx    context.Context
	db     *db.Database
	router chi.Router
}

func (s *TestSuite) SetupSuite() {
	s.ctx = context.Background()
	var cfg config.Config
	err := confita.NewLoader(
		env.NewBackend(),
	).Load(s.ctx, &cfg)
	s.Require().NoError(err, "failed to parse config")

	postgres, err := db.NewDatabase(cfg.Postgres)
	s.Require().NoError(err, "failed to connect postgresql")
	s.db = postgres
	s.TearDownSuite()

	err = postgres.Init(s.ctx)
	if err != nil {
		fmt.Printf("failed to migrate database: %s\n", err.Error())
		return
	}

	handlers := handler.NewHandler(postgres, fmt.Sprintf("%s:%s", cfg.Host, cfg.Port))
	s.router = handlers.Router()
}

func (s *TestSuite) TearDownSuite() {
	if s.db == nil {
		return
	}

	err := s.db.Drop(s.ctx)
	if err != nil {
		fmt.Printf("failed to drop database: %s\n", err.Error())
	}
}

func (s *TestSuite) TearDownTest() {
	if s.db == nil {
		return
	}

	err := s.db.Clean(s.ctx)
	if err != nil {
		fmt.Printf("failed to clean database: %s\n", err.Error())
	}
}

func (s *TestSuite) doRequest(r *http.Request) (int, string) {
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, r)
	return rr.Code, rr.Body.String()
}
