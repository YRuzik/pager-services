package transfers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"log"
	pager_chat "pager-services/pkg/api/pager_api/chat"
	common "pager-services/pkg/api/pager_api/common"
	pager_transfers "pager-services/pkg/api/pager_api/transfers"
	"pager-services/pkg/mongo_ops"
	"pager-services/pkg/namespaces"
	"pager-services/pkg/utils"
)

func UpdateUserProfile(ctx context.Context, userID string, updatedProfile *common.PagerProfile) (*common.PagerProfile, error) {
	mongoUserId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, utils.MentorError("failed set userId to objectID", codes.Internal, &common.PagerError{
			Code:    common.PagerError_INTERNAL,
			Details: err.Error(),
		})
	}
	nProfile := &common.PagerProfile{
		UserId:         updatedProfile.UserId,
		Email:          updatedProfile.Email,
		Avatar:         updatedProfile.Avatar,
		Login:          updatedProfile.Login,
		Online:         updatedProfile.Online,
		LastSeenMillis: updatedProfile.LastSeenMillis,
	}
	if err := UpdateData(ctx, mongo_ops.CollectionsPoll.ProfileCollection, namespaces.ProfileSection(userID), pager_transfers.ProfileStreamRequest_profile_info.String(), nProfile, mongoUserId); err != nil {
		log.Print(err)
		return nil, utils.MentorError("failed to update data profile", codes.Internal, &common.PagerError{
			Code:    common.PagerError_INTERNAL,
			Details: err.Error(),
		})
	}
	member := &pager_chat.ChatMember{
		Id:             updatedProfile.UserId,
		Email:          updatedProfile.Email,
		Avatar:         updatedProfile.Avatar,
		Login:          updatedProfile.Login,
		Online:         updatedProfile.Online,
		LastSeenMillis: updatedProfile.LastSeenMillis,
	}
	err = UpdateData(ctx, mongo_ops.CollectionsPoll.MembersCollection, namespaces.MemberSection(userID), pager_transfers.ChatMemberRequest_member_info.String(), member, mongoUserId)
	if err != nil {
		return nil, utils.MentorError("failed to update data profile", codes.Internal, &common.PagerError{
			Code:    common.PagerError_INTERNAL,
			Details: err.Error(),
		})
	}
	return updatedProfile, nil
}
