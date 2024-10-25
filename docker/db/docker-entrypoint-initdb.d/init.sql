CREATE USER tgtime_auth_user WITH PASSWORD 'tgtime_auth_password';
CREATE DATABASE tgtime_auth;
GRANT ALL PRIVILEGES ON DATABASE tgtime_auth TO tgtime_auth_user;