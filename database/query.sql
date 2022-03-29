-- name: GetPlmn :one
SELECT * FROM plmn
WHERE id = ? LIMIT 1;

-- name: ListPlmn :many
SELECT * FROM plmn;

-- name: CreatePlmn :execresult
INSERT INTO plmn (
	mcc, mnc, iso, country, country_code, network
) VALUES (
  ?, ?, ?, ?, ?, ?
);
