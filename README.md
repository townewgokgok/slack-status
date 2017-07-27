# slack-status

Changes your Slack user custom status from CLI

# Install

```
$ go install github.com/townewgokgok/slack-status
```

# Configure your settings

```
$ slack-status -e
```

or

```
$ vi $GOPATH/src/github.com/townewgokgok/slack-status/settings.yml
```

to edit your settings file like

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
  -i  Append information of the music playing on iTunes
  -l  Append information of the music playing on last.fm
  -v  View current status
  -w  Watch changes (with -i or -l)

Templates:
  - home : 🏠 Working remotely
  - lunch : 🍴 Having lunch
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
🎵 Satellite Young - Sniper Rouge (feat. Mitch Murder)
```

```
$ slack-status -i home
🏠 Working remotely 🎵 Satellite Young - Sniper Rouge (feat. Mitch Murder)
```

```
$ slack-status -i -w home
🏠 Working remotely 🎵 Satellite Young - Sniper Rouge (feat. Mitch Murder)
🏠 Working remotely 🎵 Satellite Young - Break! Break! Tic! Tac!
🏠 Working remotely 🎵 Satellite Young - Geeky Boyfriend
```
