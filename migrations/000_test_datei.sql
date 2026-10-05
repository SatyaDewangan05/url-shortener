UPDATE urls
SET original_url = $1
WHERE short_code = $2;