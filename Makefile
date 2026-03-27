build:
	cd backend && GOOS=linux GOARCH=amd64 go build -o cmd/api/bootstrap cmd/api/main.go
	cd backend && GOOS=linux GOARCH=amd64 go build -o cmd/youtube/bootstrap cmd/youtube/main.go

deploy: build
	cd infra && cdk deploy --all