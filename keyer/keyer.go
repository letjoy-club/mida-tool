package keyer

func User(userID string) string {
	return "user:" + userID
}

func UserMatching(userID string) string {
	return "user:" + userID + ":matching"
}

func UserMotion(userID string) string {
	return "user:" + userID + ":motion"
}

func UserQuota(userID string) string {
	return "user:" + userID + ":quota"
}

func UserMatchingInvitation(userID string) string {
	return "user:" + userID + ":matching-invitation"
}
