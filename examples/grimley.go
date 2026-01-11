version: "0.1"

project:
  name: test
  environment: local

database:
  engine: postgres
  host: localhost
  port: 5432
  name: test
  ssl: false
  credentials:
    user: test
    password_env: TEST_DB_PASSWORD

domains: {}
