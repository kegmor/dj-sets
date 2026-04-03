build:
	cd backend && GOOS=linux GOARCH=amd64 go build -o cmd/api/bootstrap ./cmd/api/
	cd backend && GOOS=linux GOARCH=amd64 go build -o cmd/youtube/bootstrap ./cmd/youtube/
	cd backend && GOOS=linux GOARCH=amd64 go build -o cmd/migrate/bootstrap ./cmd/migrate/

deploy: build
	cd infra && cdk deploy --all

migrate:
	aws lambda invoke --function-name dj-sets-migrate --payload '{}' response.json && cat response.json