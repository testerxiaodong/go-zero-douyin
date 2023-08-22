api:
	goctl api go --api ./api/desc/comment.api --dir ./api/ --style=goZero --home=../../../deploy/goctl

rpc:
	goctl rpc protoc ./rpc/pb/commentrpc.proto --go_out=./rpc --go-grpc_out=./rpc --zrpc_out=./rpc --style=goZero --home=../../../deploy/goctl

gen:
	go run ../gen/generate.go

.PHONY: api rpc gen