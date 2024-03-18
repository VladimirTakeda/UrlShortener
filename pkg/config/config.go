package config

import "git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/db"

type Config struct {
	Host string `config:"APP_HOST" yaml:"host"`
	Port string `config:"APP_PORT" yaml:"port"`

	Postgres db.Config `config:"postgres"`
}
