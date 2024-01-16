package transfers

import (
	"log"
	pager_transfers "pager-services/pkg/api/pager_api/transfers"
	"pager-services/pkg/mongo_ops"
	"pager-services/pkg/utils"
)

var _ pager_transfers.PagerStreamsServer = (*PagerStreams)(nil)

type PagerStreams struct {
}

func (p PagerStreams) StreamProfile(request *pager_transfers.ProfileStreamRequest, server pager_transfers.PagerStreams_StreamProfileServer) error {
	//TODO implement me
	panic("implement me")
}

func (p PagerStreams) StreamChat(request *pager_transfers.ChatStreamRequest, server pager_transfers.PagerStreams_StreamChatServer) error {
	ctx := server.Context()
	watch := utils.WatchFlag(ctx)

	for item := range ReadStream(server.Context(), mongo_ops.CollectionsPoll.ChatCollection, "test", watch) {
		if err := item.IsError(); err != nil {
			log.Default().Println(err)
			return err
		} else {
			if err := server.Send(item.TransferObject); err != nil {
				log.Fatal(err)
			}
		}
	}

	return nil
}
