@REM SET GOOS=linux
SET GOOS=windows
SET GOARCH=amd64
go build -ldflags "-s -w"