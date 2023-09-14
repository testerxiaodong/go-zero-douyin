api:
	goctl api go --api ./api/desc/search.api --dir ./api/ --style=goZero --home=../../../deploy/goctl

rpc:
	goctl rpc protoc ./rpc/pb/search.proto --go_out=./rpc --go-grpc_out=./rpc --zrpc_out=./rpc --style=goZero --home=../../../deploy/goctl
	sed -i "" 's/,omitempty//g' ./rpc/pb/*.pb.go

gen:
	go run ../gen/generate.go

mock:
	mockgen -source=./rpc/search/search.go -destination=./rpc/mock/search_rpc_mock.go -package=mock

dockerfile-api:
	cd api && goctl docker -go search.go --version 1.19 --home=../../../deploy/goctl

dockerfile-rpc:
	cd rpc && goctl docker -go search.go --version 1.19 --home=../../../deploy/goctl

.PHONY: api rpc gen mock dockerfile-api dockerfile-rpc