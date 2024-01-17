package transfers

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pager_transfers "pager-services/pkg/api/pager_api/transfers"
	"pager-services/pkg/mongo_ops"
	"pager-services/pkg/utils"
)

type AuthRegisterData struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
	Login    string             `bson:"login"`
}
type AuthLoginData struct {
	Identity string
	Password string `bson:"password"`
}
type AuthDataForCollection1 struct {
	Email string `bson:"email"`
	Login string `bson:"login"`
}
type AuthDataForCollection2 struct {
	ID       primitive.ObjectID `bson:"_id"`
	Password string             `bson:"password"`
}

func InsertDataWithID(ctx context.Context, collection *mongo.Collection, sectionId string, streamType string, payload interface{}) (primitive.ObjectID, error) {
	if serializedData, err := utils.CustomMarshal(&payload); err == nil {
		item := &pager_transfers.TransferObject{
			SectionId: sectionId,
			Data:      serializedData,
			Type:      streamType,
		}
		bsonItem := mongo_ops.ProtoTObjectToBSON(item)
		result, err := collection.InsertOne(ctx, bsonItem)
		if err != nil {
			return primitive.NilObjectID, err
		}
		return result.InsertedID.(primitive.ObjectID), nil
	} else {
		return primitive.NilObjectID, err
	}
}

func InsertAuthData(ctx context.Context, payload *AuthRegisterData) error {

	payloadForCollection1 := &AuthDataForCollection1{
		Email: payload.Email,
		Login: payload.Login,
	}

	result1Id, err := InsertDataWithID(ctx, mongo_ops.CollectionsPoll.ProfileCollection, "test", pager_transfers.ProfileStreamRequest_profile_info.String(), payloadForCollection1)
	if err != nil {
		return err
	}

	payloadForCollection2 := &AuthDataForCollection2{
		ID:       result1Id,
		Password: payload.Password,
	}

	if _, err := mongo_ops.CollectionsPoll.UsersCollection.InsertOne(ctx, payloadForCollection2); err != nil {
		if _, err := mongo_ops.CollectionsPoll.ProfileCollection.DeleteOne(ctx, bson.M{"_id": result1Id}); err != nil {
			return err
		}
		return err
	}

	return nil
}
func IsUserExistsWithData(ctx context.Context, sectionId string, email, login string) (bool, error) {
	filter := bson.D{
		{"section_id", sectionId},
		{"type", "profile_info"},
	}

	cursor, err := mongo_ops.CollectionsPoll.ProfileCollection.Find(ctx, filter)
	if err != nil {
		return false, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, ctx)

	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return false, err
		}

		dataBase64, ok := result["data"].(primitive.Binary)
		if !ok {
			return false, status.Error(codes.InvalidArgument, "data field is not of type Binary")
		}

		var userData AuthDataForCollection1
		if err := utils.CustomUnmarshal(dataBase64.Data, &userData); err != nil {
			return false, err
		}

		if userData.Email == email && userData.Login == login {
			return true, nil
		}
	}

	return false, nil
}

func FindUserIDByIdentifier(ctx context.Context, sectionId, identifier string) (primitive.ObjectID, error) {
	filter := bson.D{
		{"section_id", sectionId},
		{"type", "profile_info"},
	}

	cursor, err := mongo_ops.CollectionsPoll.ProfileCollection.Find(ctx, filter)
	if err != nil {
		return primitive.NilObjectID, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, ctx)

	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return primitive.NilObjectID, err
		}

		dataBase64, ok := result["data"].(primitive.Binary)
		if !ok {
			return primitive.NilObjectID, status.Error(codes.InvalidArgument, "data field is not of type Binary")
		}

		var userData AuthDataForCollection1
		if err := utils.CustomUnmarshal(dataBase64.Data, &userData); err != nil {
			return primitive.NilObjectID, err
		}

		if userData.Email == identifier || userData.Login == identifier {
			return result["_id"].(primitive.ObjectID), nil
		}
	}

	return primitive.NilObjectID, status.Error(codes.NotFound, "user not found")
}

func GetHashedPasswordByID(ctx context.Context, userID primitive.ObjectID) ([]byte, error) {
	filter := bson.D{
		{"_id", userID},
	}

	var result AuthDataForCollection2
	if err := mongo_ops.CollectionsPoll.UsersCollection.FindOne(ctx, filter).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return []byte(result.Password), nil
}
