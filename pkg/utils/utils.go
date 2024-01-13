package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	pager_chat "pager-services/pkg/api/pager_api/chat"
	"pager-services/pkg/chat_actions"
)

func RegisterGrpcServices(registrar grpc.ServiceRegistrar) {
	pager_chat.RegisterChatActionsServer(registrar, &chat_actions.PagerChat{})
}

func CustomMarshal(v interface{}) ([]byte, error) {
	if p, ok := v.(proto.Message); ok {
		return (&proto.MarshalOptions{Deterministic: true}).Marshal(p)
	}
	return json.Marshal(v)
}

func CustomUnmarshal(data []byte, v interface{}) error {
	if p, ok := v.(proto.Message); ok {
		return proto.Unmarshal(data, p)
	}
	return json.Unmarshal(data, v)
}

func WatchFlag(ctx context.Context) bool {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		fmt.Println("error while parsing md")
	}
	_, has := md["watch"]
	if has {
		return true
	} else {
		return false
	}
}
