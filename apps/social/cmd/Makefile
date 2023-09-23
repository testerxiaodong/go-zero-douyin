api:
	goctl api go --api ./api/desc/social.api --dir ./api/ --style=goZero --home=../../../deploy/goctl

rpc:
	goctl rpc protoc ./rpc/pb/social.proto --go_out=./rpc --go-grpc_out=./rpc --zrpc_out=./rpc --style=goZero --home=../../../deploy/goctl
	sed -i "" 's/,omitempty//g' ./rpc/pb/*.pb.go

gen:
	goctl model mysql datasource -url="root:my-secret-pw@tcp(127.0.0.1:3306)/go_zero_douyin" -table="comment"  -dir="./rpc/internal/model" -cache=true --style=goZero --home=../../../deploy/goctl
	goctl model mysql datasource -url="root:my-secret-pw@tcp(127.0.0.1:3306)/go_zero_douyin" -table="comment_count"  -dir="./rpc/internal/model" -cache=true --style=goZero --home=../../../deploy/goctl
	goctl model mysql datasource -url="root:my-secret-pw@tcp(127.0.0.1:3306)/go_zero_douyin" -table="like"  -dir="./rpc/internal/model" -cache=true --style=goZero --home=../../../deploy/goctl
	goctl model mysql datasource -url="root:my-secret-pw@tcp(127.0.0.1:3306)/go_zero_douyin" -table="like_count"  -dir="./rpc/internal/model" -cache=true --style=goZero --home=../../../deploy/goctl
	goctl model mysql datasource -url="root:my-secret-pw@tcp(127.0.0.1:3306)/go_zero_douyin" -table="follow"  -dir="./rpc/internal/model" -cache=true --style=goZero --home=../../../deploy/goctl
	goctl model mysql datasource -url="root:my-secret-pw@tcp(127.0.0.1:3306)/go_zero_douyin" -table="follow_count"  -dir="./rpc/internal/model" -cache=true --style=goZero --home=../../../deploy/goctl

mock:
	mockgen -source=./rpc/internal/model/commentCountModel_gen.go -destination=./rpc/mock/commentCountModel_mock.go -package=mock
	mockgen -source=./rpc/internal/model/commentModel_gen.go -destination=./rpc/mock/commentModel_mock.go -package=mock
	mockgen -source=./rpc/internal/model/likeCountModel_gen.go -destination=./rpc/mock/likeCountModel_mock.go -package=mock
	mockgen -source=./rpc/internal/model/likeModel_gen.go -destination=./rpc/mock/likeModel_mock.go -package=mock
	mockgen -source=./rpc/internal/model/followCountModel_gen.go -destination=./rpc/mock/followCountModel_mock.go -package=mock
	mockgen -source=./rpc/internal/model/followModel_gen.go -destination=./rpc/mock/followModel_mock.go -package=mock
	mockgen -source=./rpc/social/social.go -destination=./rpc/mock/social_rpc_mock.go -package=mock

dockerfile-api:
	cd api && goctl docker -go social.go --version 1.19 --home=../../../deploy/goctl

dockerfile-rpc:
	cd rpc && goctl docker -go social.go --version 1.19 --home=../../../deploy/goctl

.PHONY: api rpc gen mock dockerfile-api dockerfile-rpc