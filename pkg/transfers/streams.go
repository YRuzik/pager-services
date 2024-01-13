package transfers

import (
	pager_transfers "pager-services/pkg/api/pager_api/transfers"
	"pager-services/pkg/mongo_ops"
)

var _ pager_transfers.PagerStreamsServer = (*PagerStreams)(nil)

type PagerStreams struct {
}

func (p PagerStreams) StreamProfile(request *pager_transfers.ProfileStreamRequest, server pager_transfers.PagerStreams_StreamProfileServer) error {
	//TODO implement me
	panic("implement me")
}

func (p PagerStreams) StreamChat(request *pager_transfers.ChatStreamRequest, server pager_transfers.PagerStreams_StreamChatServer) error {
	ReadStream(server.Context(), mongo_ops.Client.Database("test_streams").Collection("transfers"), "test")
	return nil
}
