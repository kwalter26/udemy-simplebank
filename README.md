# Go Simple Bank (Udemy Course) [![Go](https://github.com/kwalter26/udemy-simplebank/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/kwalter26/udemy-simplebank/actions/workflows/ci.yml)
### Migration up 

```bash
migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
```

### sqlc setup

[SITE](https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html)
