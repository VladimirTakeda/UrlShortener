package db

import "fmt"

var (
	LinkNotFoundError = fmt.Errorf("link doesn't exist")

	SuffixNotFoundError = fmt.Errorf("short suffix doesn't exist")

	SecretKeyNotFoundError = fmt.Errorf("secret key doesn't exist")
)
