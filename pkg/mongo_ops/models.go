package mongo_ops

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pager_transfers "pager-services/pkg/api/pager_api/transfers"
)

func BSONToProtoTObject(item *TransferObjectBSON) *pager_transfers.TransferObject {
	return &pager_transfers.TransferObject{
		Id:        item.ID.Hex(),
		SectionId: item.SectionID,
		Data:      item.Data,
		Type:      item.Type,
		SeqNumber: item.SeqNumber,
	}
}

func MapTObjectToProto(item map[string]interface{}) (*pager_transfers.TransferObject, error) {
	if fullDocument, ok := item["fullDocument"].(map[string]interface{}); !ok {
		return nil, status.Error(codes.Unknown, "unknown map format")
	} else {
		return &pager_transfers.TransferObject{
			Id:        fullDocument["_id"].(primitive.ObjectID).Hex(),
			SectionId: fullDocument["section_id"].(string),
			Data:      (fullDocument["data"].(primitive.Binary)).Data,
			Type:      fullDocument["type"].(string),
			SeqNumber: fullDocument["seq_number"].(int64),
		}, nil
	}
}

type TransferObjectBSON struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	SectionID string             `bson:"section_id"`
	Data      []byte             `bson:"data"`
	Type      string             `bson:"type"`
	SeqNumber int64              `bson:"seq_number"`
}
