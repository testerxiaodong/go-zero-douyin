api:
	goctl api go --api ./api/desc/video.api --dir ./api/ --style=goZero --home=../../../deploy/goctl

rpc:
	goctl rpc protoc ./rpc/pb/video.proto --go_out=./rpc --go-grpc_out=./rpc --zrpc_out=./rpc --style=goZero --home=../../../deploy/goctl

gen:
	go run ../gen/generate.go

.PHONY: api rpc gen