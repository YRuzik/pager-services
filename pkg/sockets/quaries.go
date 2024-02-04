package sockets

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"pager-services/pkg/mongo_ops"
	"time"
)

// UpdateUserOnlineStatusByUserID will update the online status of the user
func UpdateUserOnlineStatusByUserID(userID string, status bool) error {
	docID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil
	}

	collection := mongo_ops.CollectionsPoll.ProfileCollection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	_, queryError := collection.UpdateOne(ctx, bson.M{"_id": docID}, bson.M{"$set": bson.M{"online": status}})
	//_, queryError := memberCollection.UpdateOne(ctx, bson.M{"_id": docID}, bson.M{"$set": bson.M{"online": status}})
	defer cancel()

	if queryError != nil {
		return errors.New("server response error")
	}
	return nil
}

//func GetUserByUserID(userID string) *pager_common.PagerProfile {
//	var userDetails *pager_common.PagerProfile
//
//	docID, err := primitive.ObjectIDFromHex(userID)
//	if err != nil {
//		return &pager_common.PagerProfile{}
//	}
//
//	collection := mongo_ops.CollectionsPoll.ProfileCollection
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//
//	_ = collection.FindOne(ctx, bson.M{
//		"_id": docID,
//	}).Decode(&userDetails)
//
//	defer cancel()
//
//	return userDetails
//}
//
//// GetAllOnlineUsers function will return the all online users
//func GetUserChats(userID string) []chat_actions.PagerChat {
//	var chats []chat_actions.PagerChat
//
//	collection := mongo_ops.CollectionsPoll.ProfileCollection
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//
//	cursor, queryError := collection.Find(ctx, bson.M{
//		"section_id": namespaces.ChatSection(userID),
//	})
//
//	defer cancel()
//
//	if queryError != nil {
//		return chats
//	}
//
//	for cursor.Next(context.TODO()) {
//		var transferObj pager_transfers.TransferObject
//		err := cursor.Decode(&transferObj)
//
//		var singleChat chat_actions.PagerChat
//		utils.CustomUnmarshal()
//
//		if err == nil {
//			onlineUsers = append(onlineUsers, UserDetailsResponsePayloadStruct{
//				ID:     singleOnlineUser.ID,
//				Online: singleOnlineUser.Online,
//				Login:  singleOnlineUser.Login,
//			})
//		}
//	}
//
//	return onlineUsers
//}
