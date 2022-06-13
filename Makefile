build:
	@cd cmd/tech-test; go build -o service main.go

buildDocker:
	@cd cmd/tech-test; env GOOS=linux CGO_ENABLED=0 go build -o service main.go
.PHONY: mocks
mocks:
	@go generate ./...

.PHONY: test
test:
	@go test ./...

.PHONY: run
run: mocks
	@go run cmd/tech-test/main.go

.PHONY: docker
docker: buildDocker
	@cd cmd/tech-test; docker build -t tech-test:latest .; rm ./service

.PHONY: docker-run
docker-run: docker
	@docker run -d -p 8080:8080 tech-test


