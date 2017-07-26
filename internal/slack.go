package internal

import "github.com/nlopes/slack"

func SetSlackUserStatus(token, text, emoji string) {
	api := slack.New(token)
	err := api.SetUserCustomStatus(text, emoji)
	if err != nil {
		panic("Failed to change status: " + err.Error())
	}
}
