api:
	goctl api go --api ./api/desc/user.api --dir ./api/ --style=goZero --home=../../../deploy/goctl

rpc:
	goctl rpc protoc ./rpc/pb/user.proto --go_out=./rpc --go-grpc_out=./rpc --zrpc_out=./rpc --style=goZero --home=../../../deploy/goctl
	sed -i "" 's/,omitempty//g' ./rpc/pb/*.pb.go

gen:
	go run ../gen/generate.go --home=../../../deploy/goctl

mock:
	mockgen -source=./rpc/internal/dao/userRepository.go -destination=./rpc/mock/userDo_mock.go -package=mock
	mockgen -source=./rpc/user/user.go -destination=./rpc/mock/user_rpc_mock.go -package=mock

dockerfile-api:
	cd api && goctl docker -go user.go --version 1.19 --home=../../../deploy/goctl

dockerfile-rpc:
	cd rpc && goctl docker -go user.go --version 1.19 --home=../../../deploy/goctl

.PHONY: api rpc gen mock dockerfile-api dockerfile-rpc