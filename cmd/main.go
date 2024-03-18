package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"

	"git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/config"
	"git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/db"
	"git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/handler"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "config folder")
	flag.Parse()
}

func main() {
	ctx := context.Background()
	var cfg config.Config
	err := confita.NewLoader(
		file.NewBackend(fmt.Sprintf("%s/default.yaml", configPath)),
		env.NewBackend(),
	).Load(ctx, &cfg)
	if err != nil {
		fmt.Printf("failed to parse config: %s\n", err.Error())
		return
	}

	postgres, err := db.NewDatabase(cfg.Postgres)
	if err != nil {
		fmt.Printf("failed to connect postgresql: %s\n", err.Error())
		return
	}

	err = postgres.Init(ctx)
	if err != nil {
		fmt.Printf("failed to migrate database: %s\n", err.Error())
		return
	}

	handlers := handler.NewHandler(postgres, fmt.Sprintf("%s:%s", cfg.Host, cfg.Port))
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: handlers.Router(),
	}

	go func() {
		fmt.Println("server started")
		_ = srv.ListenAndServe()
	}()

	// wait for interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// attempt a graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
