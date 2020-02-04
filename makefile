all:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" subjectmilter.go

compress:
	upx subjectmilter
