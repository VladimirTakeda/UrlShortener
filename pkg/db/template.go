package db

const (
	initRequest = `
		CREATE TABLE IF NOT EXISTS links (
			short_suffix TEXT PRIMARY KEY,
			link TEXT,
			secret_key TEXT UNIQUE,
			expiration_date timestamp,
			clicks INTEGER NOT NULL DEFAULT 0
		);
	`

	dropRequest = `
		DROP TABLE IF EXISTS links;
	`

	cleanRequest = `
		DELETE FROM links;
	`

	saveRequest = `
		INSERT INTO links(short_suffix, link, secret_key) 
			VALUES ($1, $2, $3);
	`

	saveVipRequest = `
		INSERT INTO links(short_suffix, link, secret_key, expiration_date) 
			VALUES ($1, $2, $3, $4);
	`

	selectBySuffixRequest = `
		SELECT * FROM links 
			WHERE short_suffix=$1 AND expiration_date IS NULL OR expiration_date > $2;
	`

	selectByLinkRequest = `
		SELECT * FROM links 
			WHERE link=$1;
	`

	selectBySecretKeyRequest = `
		SELECT * FROM links 
			WHERE secret_key=$1;
	`

	deleteBySecretKeyRequest = `
		DELETE FROM links 
			WHERE secret_key=$1;
	`

	incrementClicksBySuffixRequest = `
		UPDATE links
			SET clicks = clicks+1
			WHERE short_suffix=$1;
	`
)
