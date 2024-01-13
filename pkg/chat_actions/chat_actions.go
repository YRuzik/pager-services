package chat_actions

import (
	context "context"
	pagerChat "pager-services/pkg/api/pager_api/chat"
)

var _ pagerChat.ChatActionsServer = (*PagerChat)(nil)

type PagerChat struct {
}

func (p PagerChat) SendMessage(ctx context.Context, message *pagerChat.ChatMessage) (*interface{}, error) {
	//TODO implement me
	panic("implement me")
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
