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
	"pager-services/pkg/namespaces"
	"pager-services/pkg/utils"
	"strings"
	"time"
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
	UserId string `bson:"user_id"`
	Email  string `bson:"email"`
	Avatar []byte `bson:"avatar"`
	Login  string `bson:"login"`
	Online bool   `bson:"online"`
}
type AuthDataForCollection2 struct {
	ID           primitive.ObjectID `bson:"_id"`
	Password     string             `bson:"password"`
	RefreshToken string             `bson:"refreshToken"`
}

func InsertAuthData(ctx context.Context, payload *AuthRegisterData) error {
	uniqueID := utils.GenerateUniqueID()

	payloadForCollection1 := &AuthDataForCollection1{
		UserId: uniqueID.Hex(),
		Email:  payload.Email,
		Avatar: nil,
		Login:  payload.Login,
		Online: false,
	}

	err := InsertData(ctx, mongo_ops.CollectionsPoll.ProfileCollection, namespaces.ProfileSection(uniqueID.Hex()), pager_transfers.ProfileStreamRequest_profile_info.String(), payloadForCollection1, uniqueID)
	if err != nil {
		return err
	}

	refreshToken, err := utils.NewRefreshToken(uniqueID, payload.Email, time.Hour*24*30)
	if err != nil {
		return err
	}
	payloadForCollection2 := &AuthDataForCollection2{
		ID:           uniqueID,
		Password:     payload.Password,
		RefreshToken: refreshToken,
	}

	if _, err := mongo_ops.CollectionsPoll.UsersCollection.InsertOne(ctx, payloadForCollection2); err != nil {
		if _, err := mongo_ops.CollectionsPoll.ProfileCollection.DeleteOne(ctx, bson.M{"_id": uniqueID}); err != nil {
			return err
		}
		return err
	}

	return nil
}

func IsUserExistsWithData(ctx context.Context, email, login string) (bool, error) {
	filter := bson.D{
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

		if userData.Email == email || userData.Login == login {
			return true, nil
		}
	}

	return false, nil
}

func FindUserIDByIdentifier(ctx context.Context, identifier string) (primitive.ObjectID, error) {
	filter := bson.D{
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

func GetHashedPasswordByIDAndRefreshToken(ctx context.Context, userID primitive.ObjectID) ([]byte, string, error) {
	filter := bson.D{
		{"_id", userID},
	}

	var result AuthDataForCollection2
	if err := mongo_ops.CollectionsPoll.UsersCollection.FindOne(ctx, filter).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, "", errors.New("user not found")
		}
		return nil, "", err
	}

	return []byte(result.Password), result.RefreshToken, nil
}

func FindUserIDsByIdentifier(ctx context.Context, identifier string) ([]string, error) {
	filter := bson.D{
		{"type", "profile_info"},
	}

	cursor, err := mongo_ops.CollectionsPoll.ProfileCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var userIDs []string

	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}

		dataBase64, ok := result["data"].(primitive.Binary)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "data field is not of type Binary")
		}

		var userData AuthDataForCollection1
		if err := utils.CustomUnmarshal(dataBase64.Data, &userData); err != nil {
			return nil, err
		}

		if strings.Contains(userData.Email, identifier) || strings.Contains(userData.Login, identifier) {
			userID, ok := result["_id"].(primitive.ObjectID)
			if !ok {
				return nil, status.Error(codes.Internal, "failed to convert ObjectID to string")
			}
			userIDs = append(userIDs, userID.Hex())
		}
	}

	if len(userIDs) == 0 {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return userIDs, nil
}

func CheckRefreshToken(ctx context.Context, token string) (bool, error) {
	filter := bson.D{
		{"refreshToken", token},
	}
	_, err := mongo_ops.CollectionsPoll.ProfileCollection.Find(ctx, filter)
	if err != nil {
		return false, status.Error(codes.NotFound, "refreshToken not found in db")
	}
	return true, nil
}
