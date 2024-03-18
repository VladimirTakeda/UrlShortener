package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	database "git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/db"
	apierrors "git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/errors"
)

const (
	shortLinkLen = 5
	secretKeyLen = 8
)

const defaultDuration = 10
const defaultDurationUnit = time.Hour

func calcDurationFunc(durationPtr *string) (time.Duration, error) {
	if durationPtr == nil {
		return defaultDurationUnit, nil
	}
	duration := *durationPtr
	if duration == "SECONDS" {
		return time.Second, nil
	}
	if duration == "MINUTES" {
		return time.Minute, nil
	}
	if duration == "HOURS" {
		return time.Hour, nil
	}
	if duration == "DAYS" {
		return time.Hour * 24, nil
	}
	return time.Second, errors.New("Unknown TTL format")
}

func checkTtlLimit(duration time.Duration) error {
	if duration > time.Hour*48 {
		return errors.New("TTL limit exceeded")
	}
	return nil
}

var shortLinkFunc = func(baseUrl, suffix string) string { return fmt.Sprintf("http://%s/%s", baseUrl, suffix) }

func (h *Handler) Short(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request ShortLinkRequest
	err := h.parseBody(r.Body, &request)
	if err != nil {
		fmt.Printf("error when parse body: %v\n", err)
		h.renderer.RenderError(w, apierrors.BadRequest{})
		return
	}

	if err := request.Validate(); err != nil {
		fmt.Printf("error when validate request: %v\n", err)
		h.renderer.RenderError(w, apierrors.BadRequest{})
		return
	}

	// generate link and check if it's unique
	var shortSuffix string
	if request.VipKey != nil {
		fmt.Printf("Hello\n")
		shortSuffix = *request.VipKey
		err = h.checkUniqueShortSuffix(ctx, shortSuffix)
		if err != nil {
			if errors.Is(err, apierrors.SuffixAlreadyExistsError{}) {
				fmt.Printf("error when set short link: %v\n", err)
				h.renderer.RenderError(w, apierrors.BadRequest{})
				return
			}
		}
	} else {
		for {
			fmt.Printf("World\n")
			shortSuffix, err = h.generateShortSuffix(ctx)
			if err == nil {
				break
			}
			if !errors.Is(err, apierrors.SuffixAlreadyExistsError{}) {
				fmt.Printf("error when generate short link: %v\n", err)
				h.renderer.RenderError(w, apierrors.InternalError{})
				return
			}
		}
	}

	// check if long url has already been saved
	link, err := h.db.SelectByLink(ctx, request.LongUrl)
	switch {
	case err == nil:
		fmt.Printf("long url \"%s\" has been already shorten with suffix %s\n", link.Link, link.ShortSuffix)
		h.renderer.RenderJSON(w, ShortLinkResponse{ShortUrl: shortLinkFunc(h.url, link.ShortSuffix)})
		return
	case errors.Is(err, database.LinkNotFoundError):
	default:
		fmt.Printf("error when select long url: %v\n", err)
		h.renderer.RenderError(w, apierrors.InternalError{})
		return
	}

	// generate secret key and check if it's unique
	var secretKey string
	for {
		secretKey, err = h.generateSecretKey(ctx)
		if err == nil {
			break
		}

		if !errors.Is(err, apierrors.SecretKeyAlreadyExistsError{}) {
			fmt.Printf("error when generate secret key: %v\n", err)
			h.renderer.RenderError(w, apierrors.InternalError{})
			return
		}
	}

	if request.VipKey != nil {
		duration, localErr := calcDurationFunc(request.TTLUnit)
		if localErr != nil {
			fmt.Printf("error when calculate TTL: %v\n", err)
			h.renderer.RenderError(w, apierrors.BadRequest{})
			return
		}
		var ttl int
		if request.TTL == nil || *request.TTL == 0 {
			ttl = defaultDuration
		} else {
			ttl = *request.TTL
		}

		if checkTtlLimit(duration*time.Duration(ttl)) != nil {
			fmt.Printf("error when calculate TTL: %v\n", err)
			h.renderer.RenderError(w, apierrors.BadRequest{})
			return
		}

		err = h.db.SaveVip(ctx, shortSuffix, request.LongUrl, secretKey, time.Now().Add(duration*time.Duration(ttl)))
	} else {
		err = h.db.Save(ctx, shortSuffix, request.LongUrl, secretKey)
	}
	if err != nil {
		fmt.Printf("error when saving short link: %v\n", err)
		h.renderer.RenderError(w, apierrors.InternalError{})
		return
	}

	fmt.Printf("short link \"%s\" with suffix \"%s\" has been successfully saved\n", request.LongUrl, shortSuffix)
	h.renderer.RenderJSON(w, ShortLinkResponse{ShortUrl: shortLinkFunc(h.url, shortSuffix), SecretKey: secretKey})
}

func (h *Handler) generateShortSuffix(ctx context.Context) (string, error) {
	shortSuffix, err := h.generate(shortLinkLen)
	if err != nil {
		return "", err
	}

	// check if short suffix has already been used
	_, err = h.db.SelectBySuffix(ctx, shortSuffix)
	switch {
	case err == nil:
		return "", apierrors.SuffixAlreadyExistsError{}
	case errors.Is(err, database.SuffixNotFoundError):
		return shortSuffix, nil
	default:
		return "", err
	}
}

func (h *Handler) checkUniqueShortSuffix(ctx context.Context, shortSuffix string) error {
	// check if short suffix has already been used
	_, err := h.db.SelectBySuffix(ctx, shortSuffix)
	switch {
	case err == nil:
		return apierrors.SuffixAlreadyExistsError{}
	case errors.Is(err, database.SuffixNotFoundError):
		return nil
	default:
		return err
	}
}

func (h *Handler) generateSecretKey(ctx context.Context) (string, error) {
	secretKey, err := h.generate(secretKeyLen)
	if err != nil {
		return "", err
	}

	// check if secret key has already been used
	_, err = h.db.SelectBySecretKey(ctx, secretKey)
	switch {
	case err == nil:
		return "", apierrors.SecretKeyAlreadyExistsError{}
	case errors.Is(err, database.SecretKeyNotFoundError):
		return secretKey, nil
	default:
		return "", err
	}
}

func (h *Handler) generate(length int) (string, error) {
	// generate random bytes
	randomBytes := make([]byte, length)
	encodedBytes := make([]byte, hex.EncodedLen(length))
	n, err := rand.Read(randomBytes)
	if n != length {
		fmt.Printf("invalid bytes generated: expected %d, got %d\n", length, n)
		return "", apierrors.InternalError{}
	}
	if err != nil {
		return "", err
	}

	// encode to make human-readable string
	hex.Encode(encodedBytes, randomBytes)
	return string(encodedBytes), nil
}
