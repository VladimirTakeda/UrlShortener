package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	client *sql.DB
}

func NewDatabase(config Config) (*Database, error) {
	connInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Database,
	)
	client, err := sql.Open("postgres", connInfo)
	if err != nil {
		return nil, err
	}

	return &Database{client: client}, nil
}

func (db *Database) Init(ctx context.Context) error {
	_, err := db.client.ExecContext(ctx, initRequest)
	return err
}

// Drop drops all tables in database. USE ONLY FOR TESTS!
func (db *Database) Drop(ctx context.Context) error {
	_, err := db.client.ExecContext(ctx, dropRequest)
	return err
}

// Clean cleans all tables in database. USE ONLY FOR TESTS!
func (db *Database) Clean(ctx context.Context) error {
	_, err := db.client.ExecContext(ctx, cleanRequest)
	return err
}

func (db *Database) Save(ctx context.Context, shortSuffix, longLink, secretKey string) error {
	_, err := db.client.ExecContext(ctx, saveRequest, shortSuffix, longLink, secretKey)
	return err
}

func (db *Database) SaveVip(ctx context.Context, shortSuffix, longLink, secretKey string, expirationDate time.Time) error {
	_, err := db.client.ExecContext(ctx, saveVipRequest, shortSuffix, longLink, secretKey, expirationDate)
	return err
}

func (db *Database) SelectBySuffix(ctx context.Context, shortSuffix string) (*Link, error) {
	row := db.client.QueryRowContext(ctx, selectBySuffixRequest, shortSuffix, time.Now())

	var link Link
	err := row.Scan(&link.ShortSuffix, &link.Link, &link.SecretKey, &link.ShortLinkExpDate, &link.Clicks)
	switch {
	case err == sql.ErrNoRows:
		return nil, SuffixNotFoundError
	case err != nil:
		return nil, err
	}

	return &link, nil
}

func (db *Database) SelectByLink(ctx context.Context, longLink string) (*Link, error) {
	row := db.client.QueryRowContext(ctx, selectByLinkRequest, longLink)

	var link Link
	err := row.Scan(&link.ShortSuffix, &link.Link, &link.SecretKey, &link.ShortLinkExpDate, &link.Clicks)
	switch {
	case err == sql.ErrNoRows:
		return nil, LinkNotFoundError
	case err != nil:
		return nil, err
	}

	return &link, nil
}

func (db *Database) SelectBySecretKey(ctx context.Context, secretKey string) (*Link, error) {
	row := db.client.QueryRowContext(ctx, selectBySecretKeyRequest, secretKey)

	var link Link
	err := row.Scan(&link.ShortSuffix, &link.Link, &link.SecretKey, &link.ShortLinkExpDate, &link.Clicks)
	switch {
	case err == sql.ErrNoRows:
		return nil, SecretKeyNotFoundError
	case err != nil:
		return nil, err
	}

	return &link, nil
}

func (db *Database) DeleteBySecretKey(ctx context.Context, secretKey string) error {
	res, err := db.client.ExecContext(ctx, deleteBySecretKeyRequest, secretKey)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return SecretKeyNotFoundError
	}

	return nil
}

func (db *Database) IncrementClicksBySuffix(ctx context.Context, shortSuffix string) error {
	_, err := db.client.ExecContext(ctx, incrementClicksBySuffixRequest, shortSuffix)
	return err
}
