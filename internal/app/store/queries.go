package store

var postgre = struct {
	selectURL                string
	insertURL                string
	insertURLIgnoreConflicts string
}{
	selectURL:                "SELECT original_url as originalURL FROM urls WHERE short_url=$1",
	insertURL:                "INSERT INTO urls (short_url, original_url) VALUES($1, $2)",
	insertURLIgnoreConflicts: "INSERT INTO urls (short_url, original_url) VALUES ($1, $2) ON CONFLICT (short_url) DO NOTHING",
}
