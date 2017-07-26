# slack-status

Changes your Slack user custom status from CLI

# Install

```
$ go install github.com/townewgokgok/slack-status
```

# Customize your settings

Your `token` can be created at [Slack "Legacy tokens" page](https://api.slack.com/custom-integrations/legacy-tokens).

```
$ vi $GOPATH/src/github.com/townewgokgok/slack-status/settings.yml
```

```yaml
token: xoxp-...
templates:
  home:
    emoji: house
    text: Working remotely
  lunch:
    emoji: fork_and_knife
    text: Having lunch
lastfm:
  user_name: ...
  api_key: ...
  secret: ...
```

# Usage

```
$ slack-status
Usage: slack-status [options..] <template ID>

Options:
  -d  Dry run
  -i  Append information of the music playing on iTunes
  -l  Append information of the music playing on last.fm
  -w  Watch changes (with -i or -l)

Templates:
  - home : :house: Working remotely
  - lunch : :fork_and_knife: Having lunch
```

# Examples

```
$ slack-status home
:house: Working remotely
```

```
$ slack-status lunch
:fork_and_knife: Having lunch
```

```
$ slack-status -i home
:house: Working remotely :musical_note: Satellite Young - Sniper Rouge (feat. Mitch Murder)
```

```
$ slack-status -i -w home
:house: Working remotely :musical_note: Satellite Young - Sniper Rouge (feat. Mitch Murder)
:house: Working remotely :musical_note: Satellite Young - Break! Break! Tic! Tac!
:house: Working remotely :musical_note: Satellite Young - Geeky Boyfriend
```
