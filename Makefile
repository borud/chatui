all: test vet lint

test:
	@echo "*** $@"
	@go test ./...

vet:
	@echo "*** $@"
	@go vet ./...

lint:
	@echo "*** $@"
	@revive ./...

