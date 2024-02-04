package chat_actions

import (
	context "context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	pagerChat "pager-services/pkg/api/pager_api/chat"
	common "pager-services/pkg/api/pager_api/common"
	pager_transfers "pager-services/pkg/api/pager_api/transfers"
	"pager-services/pkg/mongo_ops"
	"pager-services/pkg/namespaces"
	"pager-services/pkg/transfers"
	"pager-services/pkg/utils"
)

var _ pagerChat.ChatActionsServer = (*PagerChat)(nil)

type PagerChat struct {
}

func (p PagerChat) UpdateManyMessages(ctx context.Context, request *pagerChat.ManyMessagesRequest) (*common.Empty, error) {
	for _, message := range request.Messages {
		messageId, err := primitive.ObjectIDFromHex(message.Id)
		if err != nil {
			return nil, err
		}
		if err := transfers.UpdateData(ctx, mongo_ops.CollectionsPoll.ChatCollection, namespaces.ChatSection(message.LinkedChatId), pager_transfers.ChatStreamRequest_messages.String(), message, messageId); err != nil {
			return nil, err
		}
	}
	return &common.Empty{}, nil
}

func (p PagerChat) SendMessage(ctx context.Context, message *pagerChat.ChatMessage) (*common.Empty, error) {
	id := utils.GenerateUniqueID()
	var updatedMessage *pagerChat.ChatMessage

	if len(message.Id) == 0 {
		updatedMessage = &pagerChat.ChatMessage{
			Id:           id.Hex(),
			Text:         message.Text,
			StampMillis:  message.StampMillis,
			Status:       pagerChat.ChatMessage_unread,
			AuthorId:     message.AuthorId,
			LinkedChatId: message.LinkedChatId,
			Updated:      false,
		}
		if err := transfers.InsertData(ctx, mongo_ops.CollectionsPoll.ChatCollection, namespaces.ChatSection(updatedMessage.LinkedChatId), pager_transfers.ChatStreamRequest_messages.String(), updatedMessage, id); err != nil {
			return nil, err
		}
	} else {
		messageId, err := primitive.ObjectIDFromHex(message.Id)
		if err != nil {
			return nil, err
		}
		updatedMessage = &pagerChat.ChatMessage{
			Id:           message.Id,
			Text:         message.Text,
			StampMillis:  message.StampMillis,
			Status:       message.Status,
			AuthorId:     message.AuthorId,
			LinkedChatId: message.LinkedChatId,
			Updated:      message.Updated,
		}
		if err := transfers.UpdateData(ctx, mongo_ops.CollectionsPoll.ChatCollection, namespaces.ChatSection(updatedMessage.LinkedChatId), pager_transfers.ChatStreamRequest_messages.String(), updatedMessage, messageId); err != nil {
			return nil, err
		}
	}

	return &common.Empty{}, nil
}

func (p PagerChat) UpdateChat(ctx context.Context, request *pagerChat.Chat) (*pagerChat.Chat, error) {
	id := utils.GenerateUniqueID()
	newChat := &pagerChat.Chat{}
	if len(request.Id) == 0 {
		newChat = &pagerChat.Chat{
			Id:        id.Hex(),
			Type:      request.Type,
			Metadata:  request.Metadata,
			MembersId: request.MembersId,
		}

		for _, memberId := range request.MembersId {
			role := &pagerChat.ChatRole{
				Id:   id.Hex(),
				Role: pagerChat.ChatRole_member,
			}
			if err := transfers.InsertData(ctx, mongo_ops.CollectionsPoll.ProfileCollection, namespaces.ProfileSection(memberId), pager_transfers.ProfileStreamRequest_chats_role.String(), role, primitive.NilObjectID); err != nil {
				return nil, err
			}
		}

		if err := transfers.InsertData(ctx, mongo_ops.CollectionsPoll.ChatCollection, namespaces.ChatSection(id.Hex()), pager_transfers.ChatStreamRequest_chat_info.String(), newChat, id); err != nil {
			return nil, err
		}
	} else {
		chatId, err := primitive.ObjectIDFromHex(request.Id)
		newChat = &pagerChat.Chat{
			Id:        request.Id,
			Type:      request.Type,
			Metadata:  request.Metadata,
			MembersId: request.MembersId,
		}
		if err != nil {
			return nil, err
		}
		if err := transfers.UpdateData(ctx, mongo_ops.CollectionsPoll.ChatCollection, namespaces.ChatSection(request.Id), pager_transfers.ChatStreamRequest_messages.String(), newChat, chatId); err != nil {
			return nil, err
		}
	}

	return newChat, nil
}
