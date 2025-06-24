package store

var postgre = struct {
	selectURL                string
	selectURLsByUserID       string
	insertURL                string
	insertURLWithUser        string
	insertURLIgnoreConflicts string
	setDeletedRows           string
	setDeleteOnceRow         string
	selectDBStatistics       string
}{
	selectURL:                "SELECT original_url as originalURL, deleted FROM urls WHERE short_url=$1",
	selectURLsByUserID:       "SELECT short_url as hash, original_url as originalURL FROM urls WHERE user_id=$1",
	insertURL:                "INSERT INTO urls (short_url, original_url) VALUES($1, $2) RETURNING id",
	insertURLWithUser:        "INSERT INTO urls (user_id, short_url, original_url) VALUES($3, $1, $2) RETURNING id",
	insertURLIgnoreConflicts: "INSERT INTO urls (user_id, short_url, original_url) VALUES ($3, $1, $2) ON CONFLICT (short_url) DO NOTHING",
	setDeletedRows:           "update urls set deleted=true where user_id=$1 and short_url in ($$$)",
	setDeleteOnceRow:         "update urls set deleted=true where user_id=$1 and short_url=$2",
	selectDBStatistics:       "select sum(urls) urls, count(users) users from (select sum(1) urls, 0 users from urls union select 0 urls, 1 users from urls group by user_id)",
}
