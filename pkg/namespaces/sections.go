package namespaces

func ProfileSection(userId string) string {
	return userId + ".profiles"
}

func ChatSection(chatId string) string {
	return chatId + ".chats"
}
