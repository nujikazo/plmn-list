CREATE TABLE plmn (
  id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  mcc VARCHAR(128) NOT NULL,
  mnc VARCHAR(128) NOT NULL,
  iso VARCHAR(8) NOT NULL,
  country VARCHAR(128) NOT NULL,
  country_code VARCHAR(128),
  network VARCHAR(128) NOT NULL,
  delete_flg TINYINT(1) NOT NULL DEFAULT 0
);