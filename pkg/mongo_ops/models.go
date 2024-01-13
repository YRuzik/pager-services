package mongo_ops

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	pager_transfers "pager-services/pkg/api/pager_api/transfers"
)

func ProtoTObjectToBSON(item *pager_transfers.TransferObject) TransferObjectBSON {
	return TransferObjectBSON{
		SectionID: item.SectionId,
		Data:      item.Data,
		Type:      item.Type,
	}
}

type TransferObjectBSON struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	SectionID string             `bson:"section_id"`
	Data      []byte             `bson:"data"`
	Type      string             `bson:"type"`
}
