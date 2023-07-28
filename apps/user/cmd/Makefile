api:
	goctl api go --api ./api/desc/user.api --dir ./api/ --style=goZero

rpc:
	goctl rpc protoc ./rpc/pb/userrpc.proto --go_out=./rpc --go-grpc_out=./rpc --zrpc_out=./rpc --style=goZero

gen:
	go run ../gen/generate.go

.PHONY: api rpc gen