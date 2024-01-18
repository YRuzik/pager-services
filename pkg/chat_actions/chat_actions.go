package chat_actions

import (
	context "context"
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

func (p PagerChat) SendMessage(ctx context.Context, message *pagerChat.ChatMessage) (*common.Empty, error) {
	id := utils.GenerateUniqueID()

	updatedMessage := &pagerChat.ChatMessage{
		Id:           id.Hex(),
		Text:         message.Text,
		StampMillis:  message.StampMillis,
		Status:       pagerChat.ChatMessage_sent,
		AuthorId:     message.AuthorId,
		LinkedChatId: message.LinkedChatId,
	}

	if err := transfers.InsertData(ctx, mongo_ops.CollectionsPoll.ChatCollection, namespaces.ChatSection(updatedMessage.LinkedChatId), pager_transfers.ChatStreamRequest_messages.String(), updatedMessage, id); err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}

func (p PagerChat) CreateChat(ctx context.Context, request *pagerChat.CreateChatRequest) (*pagerChat.Chat, error) {
	//userId := ctx.Value("user_id")
	var newChat *pagerChat.Chat

	//if request.Type == pagerChat.ChatType_group {
	//	newChat = &pagerChat.Chat{
	//		Id:       uuid.NewString(),
	//		Type:     pagerChat.ChatType_group,
	//		Rules:    request.Rules,
	//		Messages: nil,
	//	}
	//} else if request.Type == pagerChat.ChatType_personal {
	//	messageList := []pagerChat.ChatMessage{}
	//	slice := append(messageList, request.Rules.GetPersonalChat().Message)
	//	newChat = &pagerChat.Chat{
	//		Id:       "",
	//		Type:     pagerChat.ChatType_personal,
	//		Rules:    request.Rules,
	//		Messages: []pagerChat.{},
	//	}
	//}

	return newChat, nil
}
