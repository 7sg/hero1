generate:
	protoc \
    		-I. \
    		-I=${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway\@v1.14.6/third_party/googleapis \
    		--go_out=plugins=grpc:go-services/historical-events/api/proto \
    		--grpc-gateway_out=logtostderr=true:go-services/historical-events/api/proto \
    		--swagger_out=logtostderr=true:. \
    		go-services/historical-events/api/proto/*.proto;

start:
	docker-compose -f infra/docker-compose.yml up --build

stop:
	docker-compose -f infra/docker-compose.yml down