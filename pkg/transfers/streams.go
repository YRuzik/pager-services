package transfers

import (
	"log"
	pager_transfers "pager-services/pkg/api/pager_api/transfers"
	"pager-services/pkg/mongo_ops"
)

var _ pager_transfers.PagerStreamsServer = (*PagerStreams)(nil)

type PagerStreams struct {
}

func (p PagerStreams) StreamProfile(request *pager_transfers.ProfileStreamRequest, server pager_transfers.PagerStreams_StreamProfileServer) error {
	for item := range ReadStream(server.Context(), mongo_ops.CollectionsPoll.ProfileCollection, "test") {
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

func (p PagerStreams) StreamChat(request *pager_transfers.ChatStreamRequest, server pager_transfers.PagerStreams_StreamChatServer) error {
	for item := range ReadStream(server.Context(), mongo_ops.CollectionsPoll.ChatCollection, "test") {
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
