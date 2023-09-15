mock:
	mockgen -source=./common/cache/redis.go -destination=./mock/redis_mock.go -package=mock
	mockgen -source=./common/utils/oss.go -destination=./mock/oss_mock.go -package=mock
	mockgen -source=./common/utils/validator.go -destination=./mock/validator_mock.go -package=mock
	mockgen -source=./common/rabbitmq/sender.go -destination=./mock/sender_mock.go -package=mock
	mockgen -source=./common/elasticService/elasticsearch.go -destination=./mock/elasticsearch_mock.go -package=mock
	mockgen -source=./common/asynq/asynq.go -destination=./mock/asynq_mock.go -package=mock

user-api:
	docker build -t user-api:v1.0 -f ./apps/user/cmd/api/Dockerfile .; \
	docker tag user-api:v1.0 47.99.140.12:8077/go-zero-douyin/user-api:v1.0; \
	docker push 47.99.140.12:8077/go-zero-douyin/user-api:v1.0

user-rpc:
	docker build -t user-rpc:v1.0 -f ./apps/user/cmd/rpc/Dockerfile .; \
	docker tag user-rpc:v1.0 47.99.140.12:8077/go-zero-douyin/user-rpc:v1.0; \
	docker push 47.99.140.12:8077/go-zero-douyin/user-rpc:v1.0

video-api:
	docker build -t video-api:v1.0 -f ./apps/video/cmd/api/Dockerfile .; \
	docker tag video-api:v1.0 47.99.140.12:8077/go-zero-douyin/video-api:v1.0; \
	docker push 47.99.140.12:8077/go-zero-douyin/video-api:v1.0

video-rpc:
	docker build -t video-rpc:v1.0 -f ./apps/video/cmd/rpc/Dockerfile .; \
	docker tag video-rpc:v1.0 47.99.140.12:8077/go-zero-douyin/video-rpc:v1.0; \
	docker push 47.99.140.12:8077/go-zero-douyin/video-rpc:v1.0

social-api:
	docker build -t social-api:v1.0 -f ./apps/social/cmd/api/Dockerfile .; \
	docker tag social-api:v1.0 47.99.140.12:8077/go-zero-douyin/social-api:v1.0; \
	docker push 47.99.140.12:8077/go-zero-douyin/social-api:v1.0

social-rpc:
	docker build -t social-rpc:v1.0 -f ./apps/social/cmd/rpc/Dockerfile .; \
	docker tag social-rpc:v1.0 47.99.140.12:8077/go-zero-douyin/social-rpc:v1.0; \
	docker push 47.99.140.12:8077/go-zero-douyin/social-rpc:v1.0

mqueue:
	docker build -t mqueue:v1.0 -f ./apps/mqueue/cmd/consumer/Dockerfile .; \
	docker tag mqueue:v1.0 47.99.140.12:8077/go-zero-douyin/mqueue:v1.0; \
	docker push 47.99.140.12:8077/go-zero-douyin/mqueue:v1.0

k8s-user-api:
	cd deploy/kubernetes; \
	rm -rf user-api.yaml; \
	goctl kube deploy -secret docker-login -replicas 2 -minReplicas 2 -requestCpu 200 -requestMem 50 -limitCpu 300 -name user-api -namespace go-zero-douyin -image 47.99.140.12:8077/go-zero-douyin/user-api:v1.0 -o user-api.yaml -port 1002 -nodePort 31002 --serviceAccount find-endpoints --home=../goctl

k8s-user-rpc:
	cd deploy/kubernetes; \
	rm -rf user-rpc.yaml; \
	goctl kube deploy -secret docker-login -replicas 2 -minReplicas 2 -requestCpu 200 -requestMem 50 -limitCpu 300 -name user-rpc -namespace go-zero-douyin -image 47.99.140.12:8077/go-zero-douyin/user-rpc:v1.0 -o user-rpc.yaml -port 1102 -nodePort 31102 --serviceAccount find-endpoints --home=../goctl

k8s-video-api:
	cd deploy/kubernetes; \
	rm -rf video-api.yaml; \
	goctl kube deploy -secret docker-login -replicas 2 -minReplicas 2 -requestCpu 200 -requestMem 50 -limitCpu 300 -name video-api -namespace go-zero-douyin -image 47.99.140.12:8077/go-zero-douyin/video-api:v1.0 -o video-api.yaml -port 1003 -nodePort 31003 --serviceAccount find-endpoints --home=../goctl

k8s-video-rpc:
	cd deploy/kubernetes; \
	rm -rf video-rpc.yaml; \
	goctl kube deploy -secret docker-login -replicas 2 -minReplicas 2 -requestCpu 200 -requestMem 50 -limitCpu 300 -name video-rpc -namespace go-zero-douyin -image 47.99.140.12:8077/go-zero-douyin/video-rpc:v1.0 -o video-rpc.yaml -port 1103 -nodePort 31103 --serviceAccount find-endpoints --home=../goctl

k8s-social-api:
	cd deploy/kubernetes; \
	rm -rf social-api.yaml; \
	goctl kube deploy -secret docker-login -replicas 2 -minReplicas 2 -requestCpu 200 -requestMem 50 -limitCpu 300 -name social-api -namespace go-zero-douyin -image 47.99.140.12:8077/go-zero-douyin/social-api:v1.0 -o social-api.yaml -port 1004 -nodePort 31004 --serviceAccount find-endpoints --home=../goctl

k8s-social-rpc:
	cd deploy/kubernetes; \
	rm -rf social-rpc.yaml; \
	goctl kube deploy -secret docker-login -replicas 2 -minReplicas 2 -requestCpu 200 -requestMem 50 -limitCpu 300 -name social-rpc -namespace go-zero-douyin -image 47.99.140.12:8077/go-zero-douyin/social-rpc:v1.0 -o social-rpc.yaml -port 1104 -nodePort 31104 --serviceAccount find-endpoints --home=../goctl

k8s-mqueue:
	cd deploy/kubernetes; \
	rm -rf mqueue.yaml; \
	goctl kube deploy -secret docker-login -replicas 2 -minReplicas 2 -requestCpu 200 -requestMem 50 -limitCpu 300 -name mqueue -namespace go-zero-douyin -image 47.99.140.12:8077/go-zero-douyin/mqueue:v1.0 -o mqueue.yaml -port 2000 --serviceAccount find-endpoints --home=../goctl

.PHONY: mock user-api user-rpc video-api video-rpc social-api social-rpc mqueue k8s-user-api k8s-user-rpc k8s-video-api k8s-video-rpc k8s-social-api k8s-social-rpc k8s-mqueue