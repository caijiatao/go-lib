SERVICE_NAME = chat_app_svr
PID_FILE := $(SERVICE_NAME).pid

api:
	echo "generate all protocol"
	cd rpc/receive_svr && goctl rpc protoc receive_svr.proto --go_out=. --go-grpc_out=. --zrpc_out=.
	cd rpc/push_svr && goctl rpc protoc push_svr.proto --go_out=. --go-grpc_out=. --zrpc_out=.
	cd rpc/user && goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=.
	cd rpc/send_svr && goctl rpc protoc send_svr.proto --go_out=. --go-grpc_out=. --zrpc_out=.
	goctl api go -api chat-app-svr.api -dir .

run	:
	echo "run all server"
	cd rpc/receive_svr && go run receivesvr.go

clean:
	echo "stop all server"

.PHONY: rpc run stop