.PHONY: test test-unit test-integration proto run

test: test-unit

test-unit:
	go test ./internal/... -count=1 -timeout 60s

test-integration:
	go test ./internal/repository/ -tags=integration -count=1 -timeout 5m -v

proto:
	protoc --proto_path=api/proto --proto_path=/usr/local/include \
		--go_out=api/gen --go_opt=paths=source_relative \
		--go-grpc_out=api/gen --go-grpc_opt=paths=source_relative \
		api/proto/task/v1/task.proto

run:
	go run .
