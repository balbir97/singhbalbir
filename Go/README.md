```
go mod init build-resume
go get gopkg.in/yaml.v3
go mod tidy
go build main.go
go run main.go
go clean
```