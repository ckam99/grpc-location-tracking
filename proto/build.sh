#!/bin/bash

PROTO_DIR=./proto

PROTO_OUT_DIR=./proto/gen

## Generate JavaScript code
#yarn run grpc_tools_node_protoc \
#    --js_out=import_style=commonjs,binary:${PROTO_DIR} \
#    --grpc_out=${PROTO_DIR} \
#    --plugin=protoc-gen-grpc=./node_modules/.bin/grpc_tools_node_protoc_plugin \
#    -I ./proto \
#    proto/*.proto
#
## Generate TypeScript code (d.ts)
#yarn run grpc_tools_node_protoc \
#    --plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts \
#    --ts_out=${PROTO_DIR} \
#    -I ./proto \
#    proto/*.proto



# Generate JavaScript code
yarn run grpc_tools_node_protoc \
    --js_out=import_style=commonjs,binary:${PROTO_OUT_DIR} \
    --grpc_out=${PROTO_OUT_DIR} \
    --plugin=protoc-gen-grpc=./node_modules/.bin/grpc_tools_node_protoc_plugin \
    --grpc-web_out=import_style=typescript,mode=grpcwebtext:${PROTO_OUT_DIR} \
    -I ${PROTO_DIR} \
    ${PROTO_DIR}/*.proto

# Generate JavaScript code
yarn run grpc_tools_node_protoc \
    --js_out=import_style=commonjs,binary:${PROTO_OUT_DIR} \
    --grpc_out=${PROTO_OUT_DIR} \
    --plugin=protoc-gen-grpc=./node_modules/.bin/grpc_tools_node_protoc_plugin \
    --grpc-web_out=import_style=typescript,mode=grpcwebtext:${PROTO_OUT_DIR} \
    -I ${PROTO_DIR} \
    ${PROTO_DIR}/*.proto