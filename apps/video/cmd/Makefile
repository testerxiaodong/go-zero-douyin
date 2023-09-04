api:
	goctl api go --api ./api/desc/video.api --dir ./api/ --style=goZero --home=../../../deploy/goctl

rpc:
	goctl rpc protoc ./rpc/pb/video.proto --go_out=./rpc --go-grpc_out=./rpc --zrpc_out=./rpc --style=goZero --home=../../../deploy/goctl

gen:
	go run ../gen/generate.go

mock:
	mockgen -source=./rpc/internal/dao/videoRepository.go -destination=./rpc/mock/videoDo_mock.go -package=mock
	mockgen -source=./rpc/video/video.go -destination=./rpc/mock/video_rpc_mock.go -package=mock

dockerfile-api:
	cd api && goctl docker -go video.go --version 1.19 --home=../../../deploy/goctl

dockerfile-rpc:
	cd rpc && goctl docker -go video.go --version 1.19 --home=../../../deploy/goctl

.PHONY: api rpc gen mock dockerfile-api dockerfile-rpc