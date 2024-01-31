package client

import (
	"context"
	pagerClient "pager-services/pkg/api/pager_api/client"
	common "pager-services/pkg/api/pager_api/common"
	"pager-services/pkg/transfers"
)

var _ pagerClient.ClientServiceServer = (*PagerClient)(nil)

type PagerClient struct {
}

func (p PagerClient) SearchUsersByIdentifier(ctx context.Context, request *pagerClient.SearchUsersRequest) (*pagerClient.SearchUsersResponse, error) {
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

func (p PagerClient) ChangeDataProfile(ctx context.Context, request *common.PagerProfile) (*common.PagerProfile, error) {
	response, err := transfers.UpdateUserProfile(ctx, request.UserId, request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
