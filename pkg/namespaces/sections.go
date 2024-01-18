package namespaces

func ProfileSection(userId string) string {
	return userId + ".profiles"
}

func ChatSection(userId string) string {
	return userId + ".chats"
}
