api:
	goctl api go --api ./api/desc/video.api --dir ./api/ --style=goZero --home=../../../deploy/goctl

rpc:
	goctl rpc protoc ./rpc/pb/video.proto --go_out=./rpc --go-grpc_out=./rpc --zrpc_out=./rpc --style=goZero --home=../../../deploy/goctl
	sed -i "" 's/,omitempty//g' ./rpc/pb/*.pb.go

doc:
	goctl api doc --dir=./api/desc --o=../../../docs/video

gen:
	goctl model mysql datasource -url="root:my-secret-pw@tcp(192.168.2.248:3306)/go_zero_douyin" -table="video"  -dir="./rpc/internal/model" -cache=true --style=goZero --home=../../../deploy/goctl
	goctl model mysql datasource -url="root:my-secret-pw@tcp(192.168.2.248:3306)/go_zero_douyin" -table="tag"  -dir="./rpc/internal/model" -cache=true --style=goZero --home=../../../deploy/goctl
	goctl model mysql datasource -url="root:my-secret-pw@tcp(192.168.2.248:3306)/go_zero_douyin" -table="section"  -dir="./rpc/internal/model" -cache=true --style=goZero --home=../../../deploy/goctl

mock:
	mockgen -source=./rpc/internal/model/videoModel_gen.go -destination=./rpc/mock/videoModel_mock.go -package=mock
	mockgen -source=./rpc/internal/model/tagModel_gen.go -destination=./rpc/mock/tagModel_mock.go -package=mock
	mockgen -source=./rpc/internal/model/sectionModel_gen.go -destination=./rpc/mock/sectionModel_mock.go -package=mock
	mockgen -source=./rpc/video/video.go -destination=./rpc/mock/video_rpc_mock.go -package=mock

dockerfile-api:
	cd api && goctl docker -go video.go --version 1.19 --home=../../../deploy/goctl

dockerfile-rpc:
	cd rpc && goctl docker -go video.go --version 1.19 --home=../../../deploy/goctl

.PHONY: api rpc doc gen mock dockerfile-api dockerfile-rpc