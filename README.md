# Promotion API

API that manages promotions. API save, deletes, gets promotions by request.

- Version 1.0.0

# Development setup

- Copy the `config-sample.yml` as `config.yml` and update its values as needed.

## Tests

- Run tests with coverage
```
$ go test ./... -coverprofile ./report/coverage.out
```

- Generate html coverage report
```
$ go tool cover -html=./report/coverage.out -o ./report/coverage.html
```

- Reports are in the `/report` folder

## Regenerate API docs (swagger)

- Make sure swag cli is installed
``
$ go get -u github.com/swaggo/swag/cmd/swag
``
- Run init in the root
``
$ swag init -g cmd/promotionapi/main.go