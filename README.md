# payment-service-go
Payment Service using Go

## Development with LiveReload using Air

* Install [Air](https://github.com/cosmtrek/air)
```shell
$ export GOPATH=$HOME/go
$ export PATH="$PATH:$GOPATH/bin"
$ curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
$ air -v
$ air
```

* Run tests
```shell
$ go test ./...
```

* Linting
```shell
$ go install github.com/mgechev/revive@latest
$ reevive
```