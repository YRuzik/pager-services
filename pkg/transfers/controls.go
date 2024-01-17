package transfers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	pager_transfers "pager-services/pkg/api/pager_api/transfers"
	"pager-services/pkg/mongo_ops"
	"pager-services/pkg/utils"
)

type StreamItem struct {
	*pager_transfers.TransferObject
	streamError error
}

func (v *StreamItem) IsError() error {
	return v.streamError
}

func InsertData(ctx context.Context, collection *mongo.Collection, sectionId string, streamType string, payload interface{}, customId primitive.ObjectID) error {
	if serializedData, err := utils.CustomMarshal(&payload); err == nil {
		item := &pager_transfers.TransferObject{
			Id:        customId.Hex(),
			SectionId: sectionId,
			Data:      serializedData,
			Type:      streamType,
		}
		bsonItem := mongo_ops.ProtoTObjectToBSON(item)
		if _, err := collection.InsertOne(ctx, bsonItem); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

// ReadStream /TODO refactor/fix StreamItem repeat
func ReadStream(ctx context.Context, collection *mongo.Collection, sectionId string) <-chan StreamItem {
	res := make(chan StreamItem, 10)

	pipeline := mongo.Pipeline{bson.D{
		{"$match",
			bson.D{
				{"fullDocument.section_id", sectionId},
			},
		},
	}}

	streamOptions := options.ChangeStream().SetFullDocument(options.UpdateLookup)

	go func() {
		defer close(res)
		if stream, err := collection.Watch(ctx, pipeline, streamOptions); err == nil {
			var changeDoc map[string]interface{}
			for stream.Next(ctx) {
				if stream.Err() == nil {
					if err := stream.Decode(&changeDoc); err != nil {
						res <- StreamItem{TransferObject: nil, streamError: err}
					} else {
						if transferObject, err := mongo_ops.MapTObjectToProto(changeDoc); err == nil {
							res <- StreamItem{TransferObject: transferObject, streamError: err}
						} else {
							res <- StreamItem{TransferObject: nil, streamError: err}
						}
					}
				} else {
					res <- StreamItem{TransferObject: nil, streamError: err}
				}
			}
		} else {
			res <- StreamItem{TransferObject: nil, streamError: err}
		}
	}()

	return res
}
