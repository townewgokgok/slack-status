package internal

import "github.com/nlopes/slack"

const SlackUserStatusMaxLength = 100

func SetSlackUserStatus(text, emoji string) {
	api := slack.New(Settings.Slack.Token)
	err := api.SetUserCustomStatus(text, emoji)
	if err != nil {
		panic("Failed to change status: " + err.Error())
	}
}

func GetSlackUserStatus() string {
	api := slack.New(Settings.Slack.Token)
	res, err := api.AuthTest()
	if err != nil {
		panic("Failed to get status: " + err.Error())
	}
	user, err := api.GetUserInfo(res.UserID)
	if err != nil {
		panic("Failed to get status: " + err.Error())
	}
	return user.Profile.StatusEmoji + " " + user.Profile.StatusText
}
