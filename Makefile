gen:
	protoc -I=. \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	 	--go_out=. \
		--go-grpc_out=. \
		logger.proto

node1:
	go run github.com/GolubAlexander/grpc-logs/cmd/node -port=7001 -label=node7001
node2:
	go run github.com/GolubAlexander/grpc-logs/cmd/node -port=7002 -label=node7002

nodes:
	go run github.com/GolubAlexander/grpc-logs/cmd/node -port=7001 -label=node7001 &\
	go run github.com/GolubAlexander/grpc-logs/cmd/node -port=7002 -label=node7002

collector:
	go run github.com/GolubAlexander/grpc-logs/cmd/collector localhost:7001 localhost:7002