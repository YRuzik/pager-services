package chat_actions

import (
	context "context"
	pagerChat "pager-services/pkg/api/pager_api/chat"
)

var _ pagerChat.ChatActionsServer = (*PagerChat)(nil)

type PagerChat struct {
}

func (p PagerChat) CreateChat(ctx context.Context, request *pagerChat.CreateChatRequest) (*pagerChat.Chat, error) {
	//TODO implement me
	panic("implement me")
}

func (p PagerChat) Chatting(server pagerChat.ChatActions_ChattingServer) error {
	//TODO implement me
	panic("implement me")
}
