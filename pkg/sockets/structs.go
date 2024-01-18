package sockets

import "github.com/gorilla/websocket"

// UserDetailsStruct is a universal struct for mapping the user details
type UserDetailsStruct struct {
	ID       string `bson:"_id,omitempty"`
	Login    string
	Online   bool
	SocketID string
}

// ConversationStruct is a universal struct for mapping the conversations
type ConversationStruct struct {
	ID         string `json:"id" bson:"_id,omitempty"`
	Message    string `json:"message"`
	ToUserID   string `json:"toUserID"`
	FromUserID string `json:"fromUserID"`
}

// UserDetailsRequestPayloadStruct represents payload for Login and Registration request
type UserDetailsRequestPayloadStruct struct {
	Username string
	Password string
}

// UserDetailsResponsePayloadStruct represents payload for Login and Registration response
type UserDetailsResponsePayloadStruct struct {
	Login  string `json:"login"`
	UserId string `json:"user_id"`
	Online bool   `json:"online"`
}

// SocketEventStruct struct of socket events
type SocketEventStruct struct {
	EventName    string      `json:"eventName"`
	EventPayload interface{} `json:"eventPayload"`
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub                 *Hub
	webSocketConnection *websocket.Conn
	send                chan SocketEventStruct
	Id                  string
}

// MessagePayloadStruct is a struct used for message Payload
type MessagePayloadStruct struct {
	FromUserID string `json:"fromUserID"`
	ToUserID   string `json:"toUserID"`
	Message    string `json:"message"`
}
