api:
	goctl api go --api ./api/desc/like.api --dir ./api/ --style=goZero --home=../../../deploy/goctl

rpc:
	goctl rpc protoc ./rpc/pb/likerpc.proto --go_out=./rpc --go-grpc_out=./rpc --zrpc_out=./rpc --style=goZero --home=../../../deploy/goctl

gen:
	go run ../gen/generate.go

.PHONY: api rpc gen