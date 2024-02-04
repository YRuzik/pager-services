package client

import (
	"context"
	"log"
	pagerClient "pager-services/pkg/api/pager_api/client"
	common "pager-services/pkg/api/pager_api/common"
	"pager-services/pkg/transfers"
	"sync"
)

var _ pagerClient.ClientServiceServer = (*PagerClient)(nil)

type PagerClient struct {
	m sync.Mutex
}

func (p *PagerClient) SearchUsersByIdentifier(ctx context.Context, request *pagerClient.SearchUsersRequest) (*pagerClient.SearchUsersResponse, error) {
	userIds, err := transfers.FindUserIDsByIdentifier(ctx, request.GetIdentifier())
	if err != nil {
		return nil, err
	}

	// Возвращаем список ID в ответе
	response := &pagerClient.SearchUsersResponse{
		UserIds: userIds,
	}

	return response, nil
}

func (p *PagerClient) ChangeDataProfile(ctx context.Context, request *common.PagerProfile) (*common.PagerProfile, error) {
	userId := ctx.Value("user_id").(string)
	log.Print(request)
	transfers.PagerLocker.Lock(request.UserId)
	response, err := transfers.UpdateUserProfile(ctx, userId, request)
	if err != nil {
		log.Print(err)
	}
	transfers.PagerLocker.Unlock(request.UserId)
	return response, nil
}
