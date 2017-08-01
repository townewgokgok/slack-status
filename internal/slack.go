package internal

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

const SlackUserStatusMaxLength = 100

func SetSlackUserStatus(text, emoji string) error {
	spn := spinner.New(spinner.CharSets[14], time.Second/30)
	spn.Start()
	defer spn.Stop()
	api := slack.New(Settings.Slack.Token)
	err := api.SetUserCustomStatus(text, emoji)
	if err != nil {
		return errors.Wrap(err, "Failed to update status")
	}
	return nil
}

func GetSlackUserStatus() (string, string, error) {
	api := slack.New(Settings.Slack.Token)
	res, err := api.AuthTest()
	if err != nil {
		return "", "", errors.Wrap(err, "Failed to get status")
	}
	user, err := api.GetUserInfo(res.UserID)
	if err != nil {
		return "", "", errors.Wrap(err, "Failed to get status")
	}
	return user.Profile.StatusEmoji, user.Profile.StatusText, nil
}
