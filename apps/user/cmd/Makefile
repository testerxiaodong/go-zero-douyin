api:
	goctl api go --api ./api/desc/user.api --dir ./api/ --style=goZero --home=../../../deploy/goctl

rpc:
	goctl rpc protoc ./rpc/pb/user.proto --go_out=./rpc --go-grpc_out=./rpc --zrpc_out=./rpc --style=goZero --home=../../../deploy/goctl
	sed -i "" 's/,omitempty//g' ./rpc/pb/*.pb.go

gen:
	goctl model mysql datasource -url="root:my-secret-pw@tcp(127.0.0.1:3306)/go_zero_douyin" -table="user"  -dir="./rpc/internal/model" -cache=true --style=goZero --home=../../../deploy/goctl

mock:
	mockgen -source=./rpc/internal/model/userModel_gen.go -destination=./rpc/mock/userModel_mock.go -package=mock
	mockgen -source=./rpc/user/user.go -destination=./rpc/mock/user_rpc_mock.go -package=mock

dockerfile-api:
	cd api && goctl docker -go user.go --version 1.19 --home=../../../deploy/goctl

dockerfile-rpc:
	cd rpc && goctl docker -go user.go --version 1.19 --home=../../../deploy/goctl

.PHONY: api rpc gen mock dockerfile-api dockerfile-rpc