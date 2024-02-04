package client

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"log"
	pager_chat "pager-services/pkg/api/pager_api/chat"
	pagerClient "pager-services/pkg/api/pager_api/client"
	common "pager-services/pkg/api/pager_api/common"
	pager_transfers "pager-services/pkg/api/pager_api/transfers"
	"pager-services/pkg/mongo_ops"
	"pager-services/pkg/namespaces"
	"pager-services/pkg/transfers"
	"pager-services/pkg/utils"
	"sync"
)

var _ pagerClient.ClientServiceServer = (*PagerClient)(nil)

type PagerClient struct {
	m sync.Mutex
}

func (p *PagerClient) ChangeConnectionState(ctx context.Context, request *pagerClient.ConnectionRequest) (*pagerClient.ConnectionRequest, error) {
	userId := ctx.Value("user_id").(string)
	log.Print(request)
	if _, err := ChangeConnectionStateBody(ctx, userId, request); err != nil {
		return nil, err
	}
	return request, nil
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

func ChangeConnectionStateBody(ctx context.Context, userId string, request *pagerClient.ConnectionRequest) (*pagerClient.ConnectionRequest, error) {
	mongoUserId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	existProfile := &common.PagerProfile{}

	if err := transfers.ReadDataByID(ctx, mongo_ops.CollectionsPoll.ProfileCollection, userId, existProfile); err != nil {
		return nil, err
	}

	existProfile.Online = request.Online
	existProfile.LastSeenMillis = request.LastStampMillis

	if err := transfers.UpdateData(ctx, mongo_ops.CollectionsPoll.ProfileCollection, namespaces.ProfileSection(userId), pager_transfers.ProfileStreamRequest_profile_info.String(), existProfile, mongoUserId); err != nil {
		log.Print(err)
		return nil, utils.MentorError("failed to update data profile", codes.Internal, &common.PagerError{
			Code:    common.PagerError_INTERNAL,
			Details: err.Error(),
		})
	}

	member := &pager_chat.ChatMember{
		Id:             existProfile.UserId,
		Email:          existProfile.Email,
		Avatar:         existProfile.Avatar,
		Login:          existProfile.Login,
		Online:         request.Online,
		LastSeenMillis: request.LastStampMillis,
	}
	err = transfers.UpdateData(ctx, mongo_ops.CollectionsPoll.MembersCollection, namespaces.MemberSection(userId), pager_transfers.ChatMemberRequest_member_info.String(), member, mongoUserId)
	if err != nil {
		return nil, utils.MentorError("failed to update data profile", codes.Internal, &common.PagerError{
			Code:    common.PagerError_INTERNAL,
			Details: err.Error(),
		})
	}
	return request, nil
}
