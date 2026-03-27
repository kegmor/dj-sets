build:
	cd backend && GOOS=linux GOARCH=amd64 go build -o cmd/api/bootstrap cmd/api/main.go

deploy: build
	cd infra && cdk deploy --all