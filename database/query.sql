-- name: GetPlmn :one
SELECT * FROM plmn
WHERE id = ? LIMIT 1;

-- name: ListPlmn :many
SELECT * FROM plmn;

-- name: UpsertPlmn :execresult
INSERT INTO plmn (
	mcc,
  mnc,
  iso,
  country,
  country_code,
  network,
  created_at,
  updated_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?
) ON DUPLICATE KEY
  UPDATE
	mcc = ?,
  mnc = ?,
  iso = ?,
  country = ?,
  country_code = ?,
  network = ?,
  updated_at = ?;
