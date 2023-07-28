package keyer

import "strconv"

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

func Invitation(invitationID string) string {
	return "invitation:" + invitationID
}

func MotionOffer(offerId int) string {
	return "motion-offer:" + strconv.Itoa(offerId)
}
