package handler

import (
	"fmt"
	"net/url"

	"git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/db"
)

type ShortLinkRequest struct {
	LongUrl string  `json:"long_url"`
	VipKey  *string `json:"vip_key"`
	TTL     *int    `json:"ttl"`
	TTLUnit *string `json:"ttl_unit"`
}

func (r ShortLinkRequest) Validate() error {
	if r.LongUrl == "" {
		return fmt.Errorf("invalid long url")
	}

	longUrl, err := url.Parse(r.LongUrl)
	if err != nil {
		return err
	}

	if longUrl.Scheme == "" {
		return fmt.Errorf("schema should be provided in long url")
	}

	return nil
}

type ShortLinkResponse struct {
	ShortUrl  string `json:"short_url"`
	SecretKey string `json:"secret_key"`
}

type InfoResponse struct {
	LongUrl  string `json:"long_url"`
	ShortUrl string `json:"short_url"`
	Clicks   int    `json:"clicks"`
}

func (r *InfoResponse) FromLink(link *db.Link, baseUrl string) {
	if link == nil {
		return
	}

	r.ShortUrl = shortLinkFunc(baseUrl, link.ShortSuffix)
	r.LongUrl = link.Link
	r.Clicks = link.Clicks
}
