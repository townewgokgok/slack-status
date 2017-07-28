English | [日本語](README.ja.md)

# slack-status

Changes your Slack user status from CLI.
Your own templates can be managed by editting the settings file.

# Requirements

- [Go 1.8](https://golang.org/)

# Install

```
$ go get github.com/townewgokgok/slack-status
$ go install github.com/townewgokgok/slack-status
```

# Configure your settings

Your settings file will be saved at `$HOME/.slack-status.yml`.

Please

```
$ slack-status -e
```

and edit it like

```yaml
slack:
  token: xoxp-...

templates:
  home:
    emoji: house
    text: Working remotely
  lunch:
    emoji: fork_and_knife
    text: Having lunch
```

Your `token` can be created at [Slack "Legacy tokens" page](https://api.slack.com/custom-integrations/legacy-tokens).

# Usage

```
$ slack-status
Usage: slack-status [options..] <template ID>

Options:
  -d  Dry run
  -e  Edit settings
  -i  Append information about the music playing on iTunes
  -l  Append information about the music playing on last.fm
  -v  View current status
  -w  Watch changes (with -i or -l)

Templates:
  home  = 🏠 Working remotely
  lunch = 🍴 Having lunch
```

# Examples

```
$ slack-status home
🏠 Working remotely
```

```
$ slack-status lunch
🍴 Having lunch
```

```
$ slack-status -i
🎵 Satellite Young - Break! Break! Tic! Tac! (from "Satellite Young")
```

```
$ slack-status -i home
🏠 Working remotely 🎵 Satellite Young - Break! Break! Tic! Tac! (from "Satellite Young")
```

```
$ slack-status -i -w home
[10:25:39] 🏠 Working remotely 🎵 Satellite Young - Break! Break! Tic! Tac! (from "Satellite Young")
[10:30:16] 🏠 Working remotely 🎵 Satellite Young - Geeky Boyfriend (from "Satellite Young")
[10:33:51] 🏠 Working remotely 🎵 Satellite Young - AI Threnody (from "Satellite Young")
```
