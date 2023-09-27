api:
	goctl api go --api ./api/desc/recommend.api --dir ./api/ --style=goZero --home=../../../deploy/goctl

rpc:
	goctl rpc protoc ./rpc/pb/recommend.proto --go_out=./rpc --go-grpc_out=./rpc --zrpc_out=./rpc --style=goZero --home=../../../deploy/goctl
	sed -i "" 's/,omitempty//g' ./rpc/pb/*.pb.go

doc:
	goctl api doc --dir=./api/desc --o=../../../docs/recommend

mock:
	mockgen -source=./rpc/recommend/recommend.go -destination=./rpc/mock/recommend_rpc_mock.go -package=mock

dockerfile-api:
	cd api && goctl docker -go recommend.go --version 1.19 --home=../../../deploy/goctl

dockerfile-rpc:
	cd rpc && goctl docker -go recommend.go --version 1.19 --home=../../../deploy/goctl

.PHONY: api rpc doc mock dockerfile-api dockerfile-rpc