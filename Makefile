proto-win:
	protoc --proto_path=./proto/pager_proto \
		--go-grpc_opt=require_unimplemented_servers=false,paths=source_relative --go-grpc_out=./pkg/api/pager_api \
		--go_opt=paths=source_relative --go_out=./pkg/api/pager_api \
		./proto/pager_proto/common/common.proto ./proto/pager_proto/common/errors.proto ./proto/pager_proto/chat/chat_actions.proto ./proto/pager_proto/transfers/item.proto \
		./proto/pager_proto/transfers/streams.proto ./proto/pager_proto/auth/auth.proto

grpc-proxy:
	grpcwebproxy --server_tls_cert_file=./cmd/pager/certs/server.crt --server_tls_key_file=./cmd/pager/certs/server.key --backend_addr=localhost:4001 --backend_tls_noverify --use_websockets --allow_all_origins --server_http_debug_port=4561