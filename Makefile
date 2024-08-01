PROTO_OUTPUT_PATH=proto/gen
PROTO_PATH=proto
SWAGGER_DOC_DIR=docs/swagger

install:
	go mod tidy
clean:
	go clean -cache
	go clean -i
	go clean -testcache
	go clean -modcache
generate:
	rm -rf ${PROTO_OUTPUT_PATH}/go/*.go
	rm -f ${SWAGGER_DOC_DIR}/*.swagger.json
	protoc --proto_path=${PROTO_PATH} \
	--go_out=${PROTO_OUTPUT_PATH}/go \
	--go_opt=paths=source_relative \
	--go-grpc_out=${PROTO_OUTPUT_PATH}/go \
	--go-grpc_opt=paths=source_relative \
	--grpc-gateway_out ${PROTO_OUTPUT_PATH}/go --grpc-gateway_opt paths=source_relative \
	--openapiv2_out=${SWAGGER_DOC_DIR} \
	--openapiv2_opt=allow_merge=true,merge_file_name=text_grpc_app \
	${PROTO_PATH}/*.proto

.PHONY: install clean generate
