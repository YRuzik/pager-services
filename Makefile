proto-win:
	protoc --proto_path=./proto/pager_proto \
		--go-grpc_opt=require_unimplemented_servers=false,paths=source_relative --go-grpc_out=./pkg/api/pager_api \
		--go_opt=paths=source_relative --go_out=./pkg/api/pager_api ./proto/pager_proto/chat/chat.proto