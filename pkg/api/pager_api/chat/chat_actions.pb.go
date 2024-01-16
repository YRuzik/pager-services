// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.1
// source: chat/chat_actions.proto

package pager_chat

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	common "pager-services/pkg/api/pager_api/common"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ChatType int32

const (
	ChatType_group    ChatType = 0
	ChatType_personal ChatType = 1
)

// Enum value maps for ChatType.
var (
	ChatType_name = map[int32]string{
		0: "group",
		1: "personal",
	}
	ChatType_value = map[string]int32{
		"group":    0,
		"personal": 1,
	}
)

func (x ChatType) Enum() *ChatType {
	p := new(ChatType)
	*p = x
	return p
}

func (x ChatType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ChatType) Descriptor() protoreflect.EnumDescriptor {
	return file_chat_chat_actions_proto_enumTypes[0].Descriptor()
}

func (ChatType) Type() protoreflect.EnumType {
	return &file_chat_chat_actions_proto_enumTypes[0]
}

func (x ChatType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ChatType.Descriptor instead.
func (ChatType) EnumDescriptor() ([]byte, []int) {
	return file_chat_chat_actions_proto_rawDescGZIP(), []int{0}
}

type ChatMessage_MessageStatus int32

const (
	ChatMessage_uploading ChatMessage_MessageStatus = 0
	ChatMessage_sent      ChatMessage_MessageStatus = 1
	ChatMessage_seen      ChatMessage_MessageStatus = 2
	ChatMessage_error     ChatMessage_MessageStatus = 3
	ChatMessage_deleted   ChatMessage_MessageStatus = 4
)

// Enum value maps for ChatMessage_MessageStatus.
var (
	ChatMessage_MessageStatus_name = map[int32]string{
		0: "uploading",
		1: "sent",
		2: "seen",
		3: "error",
		4: "deleted",
	}
	ChatMessage_MessageStatus_value = map[string]int32{
		"uploading": 0,
		"sent":      1,
		"seen":      2,
		"error":     3,
		"deleted":   4,
	}
)

func (x ChatMessage_MessageStatus) Enum() *ChatMessage_MessageStatus {
	p := new(ChatMessage_MessageStatus)
	*p = x
	return p
}

func (x ChatMessage_MessageStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ChatMessage_MessageStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_chat_chat_actions_proto_enumTypes[1].Descriptor()
}

func (ChatMessage_MessageStatus) Type() protoreflect.EnumType {
	return &file_chat_chat_actions_proto_enumTypes[1]
}

func (x ChatMessage_MessageStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ChatMessage_MessageStatus.Descriptor instead.
func (ChatMessage_MessageStatus) EnumDescriptor() ([]byte, []int) {
	return file_chat_chat_actions_proto_rawDescGZIP(), []int{2, 0}
}

type Chat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string         `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`                                  //идентификатор чата
	Type      ChatType       `protobuf:"varint,2,opt,name=type,proto3,enum=com.pager.api.ChatType" json:"type,omitempty"` //тип чата
	Metadata  *ChatMetadata  `protobuf:"bytes,3,opt,name=metadata,proto3,oneof" json:"metadata,omitempty"`                //дополнительная информация
	Messages  []*ChatMessage `protobuf:"bytes,4,rep,name=messages,proto3" json:"messages,omitempty"`                      //сообщения в чате
	MembersId []string       `protobuf:"bytes,5,rep,name=members_id,json=membersId,proto3" json:"members_id,omitempty"`   //идентификаторы участников
}

func (x *Chat) Reset() {
	*x = Chat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_chat_actions_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Chat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Chat) ProtoMessage() {}

func (x *Chat) ProtoReflect() protoreflect.Message {
	mi := &file_chat_chat_actions_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Chat.ProtoReflect.Descriptor instead.
func (*Chat) Descriptor() ([]byte, []int) {
	return file_chat_chat_actions_proto_rawDescGZIP(), []int{0}
}

func (x *Chat) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Chat) GetType() ChatType {
	if x != nil {
		return x.Type
	}
	return ChatType_group
}

func (x *Chat) GetMetadata() *ChatMetadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *Chat) GetMessages() []*ChatMessage {
	if x != nil {
		return x.Messages
	}
	return nil
}

func (x *Chat) GetMembersId() []string {
	if x != nil {
		return x.MembersId
	}
	return nil
}

type ChatMetadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title     string  `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`                                // название группового чата
	AvatarUrl *string `protobuf:"bytes,2,opt,name=avatar_url,json=avatarUrl,proto3,oneof" json:"avatar_url,omitempty"` // обложка для группового чата
}

func (x *ChatMetadata) Reset() {
	*x = ChatMetadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_chat_actions_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChatMetadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatMetadata) ProtoMessage() {}

func (x *ChatMetadata) ProtoReflect() protoreflect.Message {
	mi := &file_chat_chat_actions_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatMetadata.ProtoReflect.Descriptor instead.
func (*ChatMetadata) Descriptor() ([]byte, []int) {
	return file_chat_chat_actions_proto_rawDescGZIP(), []int{1}
}

func (x *ChatMetadata) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *ChatMetadata) GetAvatarUrl() string {
	if x != nil && x.AvatarUrl != nil {
		return *x.AvatarUrl
	}
	return ""
}

type ChatMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Text         *string                   `protobuf:"bytes,1,opt,name=text,proto3,oneof" json:"text,omitempty"`                                             //текст сообщения
	StampMillis  int64                     `protobuf:"varint,2,opt,name=stamp_millis,json=stampMillis,proto3" json:"stamp_millis,omitempty"`                 //время отправки сообщения
	Status       ChatMessage_MessageStatus `protobuf:"varint,3,opt,name=status,proto3,enum=com.pager.api.ChatMessage_MessageStatus" json:"status,omitempty"` //статус сообщения
	AuthorId     string                    `protobuf:"bytes,4,opt,name=author_id,json=authorId,proto3" json:"author_id,omitempty"`                           //автор сообщения
	LinkedChatId string                    `protobuf:"bytes,5,opt,name=linked_chat_id,json=linkedChatId,proto3" json:"linked_chat_id,omitempty"`             //связанный чат
}

func (x *ChatMessage) Reset() {
	*x = ChatMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_chat_actions_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChatMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatMessage) ProtoMessage() {}

func (x *ChatMessage) ProtoReflect() protoreflect.Message {
	mi := &file_chat_chat_actions_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatMessage.ProtoReflect.Descriptor instead.
func (*ChatMessage) Descriptor() ([]byte, []int) {
	return file_chat_chat_actions_proto_rawDescGZIP(), []int{2}
}

func (x *ChatMessage) GetText() string {
	if x != nil && x.Text != nil {
		return *x.Text
	}
	return ""
}

func (x *ChatMessage) GetStampMillis() int64 {
	if x != nil {
		return x.StampMillis
	}
	return 0
}

func (x *ChatMessage) GetStatus() ChatMessage_MessageStatus {
	if x != nil {
		return x.Status
	}
	return ChatMessage_uploading
}

func (x *ChatMessage) GetAuthorId() string {
	if x != nil {
		return x.AuthorId
	}
	return ""
}

func (x *ChatMessage) GetLinkedChatId() string {
	if x != nil {
		return x.LinkedChatId
	}
	return ""
}

type CreateChatRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type      ChatType      `protobuf:"varint,1,opt,name=type,proto3,enum=com.pager.api.ChatType" json:"type,omitempty"` //тип чата
	Metadata  *ChatMetadata `protobuf:"bytes,2,opt,name=metadata,proto3,oneof" json:"metadata,omitempty"`                //дополнительная информация
	MembersId []string      `protobuf:"bytes,3,rep,name=members_id,json=membersId,proto3" json:"members_id,omitempty"`   //идентификаторы участников
}

func (x *CreateChatRequest) Reset() {
	*x = CreateChatRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_chat_actions_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateChatRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateChatRequest) ProtoMessage() {}

func (x *CreateChatRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chat_chat_actions_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateChatRequest.ProtoReflect.Descriptor instead.
func (*CreateChatRequest) Descriptor() ([]byte, []int) {
	return file_chat_chat_actions_proto_rawDescGZIP(), []int{3}
}

func (x *CreateChatRequest) GetType() ChatType {
	if x != nil {
		return x.Type
	}
	return ChatType_group
}

func (x *CreateChatRequest) GetMetadata() *ChatMetadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *CreateChatRequest) GetMembersId() []string {
	if x != nil {
		return x.MembersId
	}
	return nil
}

var File_chat_chat_actions_proto protoreflect.FileDescriptor

var file_chat_chat_actions_proto_rawDesc = []byte{
	0x0a, 0x17, 0x63, 0x68, 0x61, 0x74, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x5f, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x63, 0x6f, 0x6d, 0x2e, 0x70,
	0x61, 0x67, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x1a, 0x13, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe5, 0x01,
	0x0a, 0x04, 0x43, 0x68, 0x61, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2b, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x61, 0x67, 0x65, 0x72,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x12, 0x3c, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x61, 0x67, 0x65,
	0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x48, 0x00, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x88, 0x01,
	0x01, 0x12, 0x36, 0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x61, 0x67, 0x65, 0x72, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52,
	0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x73, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x6d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x49, 0x64, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x6d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x22, 0x57, 0x0a, 0x0c, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x22, 0x0a, 0x0a, 0x61,
	0x76, 0x61, 0x74, 0x61, 0x72, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x00, 0x52, 0x09, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x55, 0x72, 0x6c, 0x88, 0x01, 0x01, 0x42,
	0x0d, 0x0a, 0x0b, 0x5f, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x5f, 0x75, 0x72, 0x6c, 0x22, 0xa3,
	0x02, 0x0a, 0x0b, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x17,
	0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04,
	0x74, 0x65, 0x78, 0x74, 0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x5f, 0x6d, 0x69, 0x6c, 0x6c, 0x69, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x4d, 0x69, 0x6c, 0x6c, 0x69, 0x73, 0x12, 0x40, 0x0a, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x28, 0x2e, 0x63, 0x6f, 0x6d,
	0x2e, 0x70, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1b, 0x0a, 0x09,
	0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0e, 0x6c, 0x69, 0x6e,
	0x6b, 0x65, 0x64, 0x5f, 0x63, 0x68, 0x61, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x6c, 0x69, 0x6e, 0x6b, 0x65, 0x64, 0x43, 0x68, 0x61, 0x74, 0x49, 0x64, 0x22,
	0x4a, 0x0a, 0x0d, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x0d, 0x0a, 0x09, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x10, 0x00, 0x12,
	0x08, 0x0a, 0x04, 0x73, 0x65, 0x6e, 0x74, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x73, 0x65, 0x65,
	0x6e, 0x10, 0x02, 0x12, 0x09, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x03, 0x12, 0x0b,
	0x0a, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x10, 0x04, 0x42, 0x07, 0x0a, 0x05, 0x5f,
	0x74, 0x65, 0x78, 0x74, 0x22, 0xaa, 0x01, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43,
	0x68, 0x61, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2b, 0x0a, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x70,
	0x61, 0x67, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x3c, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x63, 0x6f, 0x6d, 0x2e,
	0x70, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x48, 0x00, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x88, 0x01, 0x01, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73,
	0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x73, 0x49, 0x64, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x2a, 0x23, 0x0a, 0x08, 0x43, 0x68, 0x61, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x09, 0x0a,
	0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x70, 0x65, 0x72, 0x73,
	0x6f, 0x6e, 0x61, 0x6c, 0x10, 0x01, 0x32, 0x93, 0x01, 0x0a, 0x0b, 0x43, 0x68, 0x61, 0x74, 0x41,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x43, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x43, 0x68, 0x61, 0x74, 0x12, 0x20, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x61, 0x67, 0x65, 0x72,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x61, 0x67,
	0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x12, 0x3f, 0x0a, 0x0b, 0x53,
	0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1a, 0x2e, 0x63, 0x6f, 0x6d,
	0x2e, 0x70, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x14, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x61, 0x67,
	0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x32, 0x5a, 0x30,
	0x70, 0x61, 0x67, 0x65, 0x72, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x61, 0x67, 0x65, 0x72, 0x5f, 0x61, 0x70, 0x69,
	0x2f, 0x63, 0x68, 0x61, 0x74, 0x3b, 0x70, 0x61, 0x67, 0x65, 0x72, 0x5f, 0x63, 0x68, 0x61, 0x74,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_chat_chat_actions_proto_rawDescOnce sync.Once
	file_chat_chat_actions_proto_rawDescData = file_chat_chat_actions_proto_rawDesc
)

func file_chat_chat_actions_proto_rawDescGZIP() []byte {
	file_chat_chat_actions_proto_rawDescOnce.Do(func() {
		file_chat_chat_actions_proto_rawDescData = protoimpl.X.CompressGZIP(file_chat_chat_actions_proto_rawDescData)
	})
	return file_chat_chat_actions_proto_rawDescData
}

var file_chat_chat_actions_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_chat_chat_actions_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_chat_chat_actions_proto_goTypes = []interface{}{
	(ChatType)(0),                  // 0: com.pager.api.ChatType
	(ChatMessage_MessageStatus)(0), // 1: com.pager.api.ChatMessage.MessageStatus
	(*Chat)(nil),                   // 2: com.pager.api.Chat
	(*ChatMetadata)(nil),           // 3: com.pager.api.ChatMetadata
	(*ChatMessage)(nil),            // 4: com.pager.api.ChatMessage
	(*CreateChatRequest)(nil),      // 5: com.pager.api.CreateChatRequest
	(*common.Empty)(nil),           // 6: com.pager.api.Empty
}
var file_chat_chat_actions_proto_depIdxs = []int32{
	0, // 0: com.pager.api.Chat.type:type_name -> com.pager.api.ChatType
	3, // 1: com.pager.api.Chat.metadata:type_name -> com.pager.api.ChatMetadata
	4, // 2: com.pager.api.Chat.messages:type_name -> com.pager.api.ChatMessage
	1, // 3: com.pager.api.ChatMessage.status:type_name -> com.pager.api.ChatMessage.MessageStatus
	0, // 4: com.pager.api.CreateChatRequest.type:type_name -> com.pager.api.ChatType
	3, // 5: com.pager.api.CreateChatRequest.metadata:type_name -> com.pager.api.ChatMetadata
	5, // 6: com.pager.api.ChatActions.CreateChat:input_type -> com.pager.api.CreateChatRequest
	4, // 7: com.pager.api.ChatActions.SendMessage:input_type -> com.pager.api.ChatMessage
	2, // 8: com.pager.api.ChatActions.CreateChat:output_type -> com.pager.api.Chat
	6, // 9: com.pager.api.ChatActions.SendMessage:output_type -> com.pager.api.Empty
	8, // [8:10] is the sub-list for method output_type
	6, // [6:8] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_chat_chat_actions_proto_init() }
func file_chat_chat_actions_proto_init() {
	if File_chat_chat_actions_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_chat_chat_actions_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Chat); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chat_chat_actions_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChatMetadata); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chat_chat_actions_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChatMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chat_chat_actions_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateChatRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_chat_chat_actions_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_chat_chat_actions_proto_msgTypes[1].OneofWrappers = []interface{}{}
	file_chat_chat_actions_proto_msgTypes[2].OneofWrappers = []interface{}{}
	file_chat_chat_actions_proto_msgTypes[3].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_chat_chat_actions_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_chat_chat_actions_proto_goTypes,
		DependencyIndexes: file_chat_chat_actions_proto_depIdxs,
		EnumInfos:         file_chat_chat_actions_proto_enumTypes,
		MessageInfos:      file_chat_chat_actions_proto_msgTypes,
	}.Build()
	File_chat_chat_actions_proto = out.File
	file_chat_chat_actions_proto_rawDesc = nil
	file_chat_chat_actions_proto_goTypes = nil
	file_chat_chat_actions_proto_depIdxs = nil
}
