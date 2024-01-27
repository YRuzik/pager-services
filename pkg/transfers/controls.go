package transfers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"seq_number", -1}}).SetLimit(1)

	pip := mongo.Pipeline{bson.D{
		{"$group", bson.D{
			{"_id", ""},
			{"maxid", bson.D{{"$max", "$id"}}},
		},
		}},
	}

	var foundElement *mongo_ops.TransferObjectBSON
	err := cursor.All(ctx, &foundElement)
	if err != nil {
		return err
	}
	//if err := singleElement.Decode(&foundElement); err != nil {
	//	return err
	//}

	if serializedData, err := utils.CustomMarshal(&payload); err == nil {
		item := &mongo_ops.TransferObjectBSON{
			ID:        customId,
			SectionID: sectionId,
			Data:      serializedData,
			Type:      streamType,
			SeqNumber: foundElement.SeqNumber + 1,
		}
		if _, err := collection.InsertOne(ctx, item); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func ReadDataByID(ctx context.Context, collection *mongo.Collection, id string, payload interface{}) error {
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return status.Error(codes.Canceled, "error while building id")
	}

	var foundElement *mongo_ops.TransferObjectBSON

	_ = collection.FindOne(ctx, bson.M{
		"_id": docID,
	}).Decode(&foundElement)

	if err := utils.CustomUnmarshal(foundElement.Data, &payload); err != nil {
		return err
	}
	return nil
}

// ReadStream /TODO refactor/fix StreamItem repeat
func ReadStream(ctx context.Context, collection *mongo.Collection, sectionId string, watch bool, additionalInfo int64) <-chan StreamItem {
	res := make(chan StreamItem, 10)

	pipeline := mongo.Pipeline{bson.D{
		{"$match",
			bson.D{
				{"fullDocument.section_id", sectionId},
			},
		},
	}}

	filter := bson.D{{"section_id", sectionId}}

	streamOptions := options.ChangeStream().SetFullDocument(options.UpdateLookup)

	findOptions := options.Find()

	go func() {
		defer close(res)
		if watch {
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
		} else {
			if current, err := collection.Find(ctx, filter, findOptions); err == nil {
				var foundElement *mongo_ops.TransferObjectBSON
				for current.Next(ctx) {
					if current.Err() == nil {
						if err := current.Decode(&foundElement); err != nil {
							res <- StreamItem{TransferObject: nil, streamError: err}
						} else {
							transferObject := mongo_ops.BSONToProtoTObject(foundElement)
							if transferObject.Type == pager_transfers.ChatStreamRequest_messages.String() {
								limit := transferObject.SeqNumber - additionalInfo
								if transferObject.SeqNumber < limit {
									res <- StreamItem{TransferObject: transferObject, streamError: err}
								}
							} else {
								res <- StreamItem{TransferObject: transferObject, streamError: err}
							}
						}
					} else {
						res <- StreamItem{TransferObject: nil, streamError: err}
					}
				}
			} else {
				res <- StreamItem{TransferObject: nil, streamError: err}
			}
		}
	}()

	return res
}
