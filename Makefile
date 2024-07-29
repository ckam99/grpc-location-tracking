PROTO_OUTPUT_PATH=proto/gen
PROTO_PATH=proto

install:
	go mod tidy
clean:
	go clean -cache
	go clean -i
	go clean -testcache
	go clean -modcache
generate:
	rm -rf ${PROTO_OUTPUT_PATH}/*.go
	protoc --proto_path=${PROTO_PATH} \
	--go_out=${PROTO_OUTPUT_PATH} \
	--go_opt=paths=source_relative \
	--go-grpc_out=${PROTO_OUTPUT_PATH} \
	--go-grpc_opt=paths=source_relative \
	${PROTO_PATH}/*.proto

.PHONY: install clean generate
