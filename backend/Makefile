soft-run:
	@go run ./cmd/api/main.go
build:
	@go build -C cmd/api/ -o ../../bin/app
run: build
	@./bin/app
test:
	@go test -v ./...
silent-test:
	@go test ./...
create-feature:
	@mkdir internal/features/${feature} && \
	touch internal/features/${feature}/${feature}.go && \
	touch internal/features/${feature}/router.go && \
	mkdir internal/features/${feature}/internal && \
	touch internal/features/${feature}/internal/handler.go && \
	touch internal/features/${feature}/internal/service.go && \
	touch internal/features/${feature}/internal/repository.go 
create-platform:
	@mkdir internal/platform/${feature}