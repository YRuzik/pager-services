package transfers

import (
	"log"
	pager_transfers "pager-services/pkg/api/pager_api/transfers"
	"pager-services/pkg/mongo_ops"
	"pager-services/pkg/namespaces"
	"pager-services/pkg/utils"
)

var _ pager_transfers.PagerStreamsServer = (*PagerStreams)(nil)

type PagerStreams struct {
}

func (p PagerStreams) StreamChatMember(request *pager_transfers.ChatMemberRequest, server pager_transfers.PagerStreams_StreamChatMemberServer) error {
	ctx := server.Context()
	watch := utils.WatchFlag(ctx)

	for item := range ReadStream(server.Context(), mongo_ops.CollectionsPoll.MembersCollection, namespaces.MemberSection(request.MemberId), watch, 0) {
		if err := item.IsError(); err != nil {
			log.Default().Println(err)
			return err
		} else {
			if err := server.Send(item.TransferObject); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p PagerStreams) StreamProfile(request *pager_transfers.ProfileStreamRequest, server pager_transfers.PagerStreams_StreamProfileServer) error {
	ctx := server.Context()
	userId := ctx.Value("user_id").(string)
	watch := utils.WatchFlag(ctx)

	for item := range ReadStream(server.Context(), mongo_ops.CollectionsPoll.ProfileCollection, namespaces.ProfileSection(userId), watch, 0) {
		if err := item.IsError(); err != nil {
			log.Default().Println(err)
			return err
		} else {
			if err := server.Send(item.TransferObject); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p PagerStreams) StreamChat(request *pager_transfers.ChatStreamRequest, server pager_transfers.PagerStreams_StreamChatServer) error {
	ctx := server.Context()
	watch := utils.WatchFlag(ctx)

	for item := range ReadStream(server.Context(), mongo_ops.CollectionsPoll.ChatCollection, namespaces.ChatSection(request.ChatId), watch, request.RequestedMessages) {
		if err := item.IsError(); err != nil {
			log.Default().Println(err)
			return err
		} else {
			if err := server.Send(item.TransferObject); err != nil {
				return err
			}
		}
	}

	return nil
}
