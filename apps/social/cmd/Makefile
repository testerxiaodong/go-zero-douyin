api:
	goctl api go --api ./api/desc/social.api --dir ./api/ --style=goZero --home=../../../deploy/goctl

rpc:
	goctl rpc protoc ./rpc/pb/social.proto --go_out=./rpc --go-grpc_out=./rpc --zrpc_out=./rpc --style=goZero --home=../../../deploy/goctl

gen:
	go run ../gen/generate.go

.PHONY: api rpc gen