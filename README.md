
# Go Clean GRPC
Go Distributed Tracing Using OpenTelemetry and Jaeger
## Stack
- Chi (net/http)
- MongoDB
## Run
Start the server using go run
```bash
  go run cmds/app/main.go
```
Start the server using [air](https://github.com/cosmtrek/air)
```bash
  make run
```
## Unit Test
Run Unit testing
```bash
  make test
```
Run Coverage
```bash
  make test/cover
```