package transfers

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	pager_transfers "pager-services/pkg/api/pager_api/transfers"
	"pager-services/pkg/utils"
)

func InsertData(ctx context.Context, collection *mongo.Collection, sectionId string, streamType string, payload interface{}) error {
	if serializedData, err := utils.CustomMarshal(payload); err == nil {
		item := &pager_transfers.TransferObject{
			Id:        uuid.NewString(),
			SectionId: sectionId,
			Data:      serializedData,
			Type:      streamType,
		}
		if _, err := collection.InsertOne(ctx, item); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

//func ReadOneData(collection *mgo.Collection, id string, payload interface{}) (interface{}, error) {
//	obj := collection.Find(bson.D{{"id", id}})
//	if obj != nil {
//		deserializedObject := utils.CustomUnmarshal(obj, payload)
//		return obj, nil
//	} else {
//		return nil, status.Error(codes.NotFound, "object not found")
//	}
//}

func ReadStream(ctx context.Context, collection *mongo.Collection, sectionId string) {
	pipeline := mongo.Pipeline{bson.D{
		{"$match",
			bson.D{
				{"fullDocument.section_id", sectionId},
				{},
			},
		},
	}}
	streamOptions := options.ChangeStream().SetFullDocument(options.UpdateLookup)

	stream, err := collection.Watch(ctx, pipeline, streamOptions)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("waiting for changes")
	var changeDoc map[string]interface{}
	for stream.Next(ctx) {
		if e := stream.Decode(&changeDoc); e != nil {
			log.Printf("error decoding: %s", e)
		}
		log.Printf("change: %+v", changeDoc)
	}
}
