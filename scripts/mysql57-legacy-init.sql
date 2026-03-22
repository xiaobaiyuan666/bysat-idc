UPDATE mysql.user SET plugin='mysql_native_password', authentication_string=PASSWORD('MySQL57Root!2026') WHERE User='root';
CREATE DATABASE IF NOT EXISTS zjmf_legacy374 DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
CREATE USER IF NOT EXISTS 'zjmf_legacy'@'127.0.0.1' IDENTIFIED BY 'ZjmfLegacy!2026';
GRANT ALL PRIVILEGES ON zjmf_legacy374.* TO 'zjmf_legacy'@'127.0.0.1';
CREATE USER IF NOT EXISTS 'zjmf_legacy'@'localhost' IDENTIFIED BY 'ZjmfLegacy!2026';
GRANT ALL PRIVILEGES ON zjmf_legacy374.* TO 'zjmf_legacy'@'localhost';
FLUSH PRIVILEGES;
