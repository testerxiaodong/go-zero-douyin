api:
	goctl api go --api ./api/desc/social.api --dir ./api/ --style=goZero --home=../../../deploy/goctl

rpc:
	goctl rpc protoc ./rpc/pb/social.proto --go_out=./rpc --go-grpc_out=./rpc --zrpc_out=./rpc --style=goZero --home=../../../deploy/goctl
	sed -i "" 's/,omitempty//g' ./rpc/pb/*.pb.go

gen:
	go run ../gen/generate.go

mock:
	mockgen -source=./rpc/internal/dao/commentRepository.go -destination=./rpc/mock/commentDo_mock.go -package=mock
	mockgen -source=./rpc/internal/dao/followRepository.go -destination=./rpc/mock/followDo_mock.go -package=mock
	mockgen -source=./rpc/internal/dao/likeRepository.go -destination=./rpc/mock/likeDo_mock.go -package=mock
	mockgen -source=./rpc/social/social.go -destination=./rpc/mock/social_rpc_mock.go -package=mock

dockerfile-api:
	cd api && goctl docker -go social.go --version 1.19 --home=../../../deploy/goctl

dockerfile-rpc:
	cd rpc && goctl docker -go social.go --version 1.19 --home=../../../deploy/goctl

.PHONY: api rpc gen mock dockerfile-api dockerfile-rpc