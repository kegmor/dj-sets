build:
	cd backend && GOOS=linux GOARCH=amd64 go build -o cmd/api/bootstrap ./cmd/api/
	cd backend && GOOS=linux GOARCH=amd64 go build -o cmd/youtube/bootstrap ./cmd/youtube/
deploy: build
	cd infra && cdk deploy --all