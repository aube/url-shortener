package store

var postgre = struct {
	selectURL                string
	selectURLsByUserID       string
	insertURL                string
	insertURLWithUser        string
	insertURLIgnoreConflicts string
}{
	selectURL:                "SELECT original_url as originalURL FROM urls WHERE short_url=$1",
	selectURLsByUserID:       "SELECT original_url as originalURL FROM urls WHERE user_id=$1",
	insertURL:                "INSERT INTO urls (short_url, original_url) VALUES($1, $2) RETURNING id",
	insertURLWithUser:        "INSERT INTO urls (user_id, short_url, original_url) VALUES($3, $1, $2) RETURNING id",
	insertURLIgnoreConflicts: "INSERT INTO urls (user_id, short_url, original_url) VALUES ($3, $1, $2) ON CONFLICT (short_url) DO NOTHING",
}
