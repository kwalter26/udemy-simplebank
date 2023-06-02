# Go Simple Bank (Udemy Course) [![Go](https://github.com/kwalter26/udemy-simplebank/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/kwalter26/udemy-simplebank/actions/workflows/ci.yml) [![codecov](https://codecov.io/gh/kwalter26/udemy-simplebank/branch/main/graph/badge.svg?token=hbYZBzkiYa)](https://codecov.io/gh/kwalter26/udemy-simplebank)

### Migration up

```bash
migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
```

Create new migration

```bash
migrate create -ext sql -dir db/migration -seq add_sessions
```

### sqlc setup

[SITE](https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html)
