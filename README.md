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
```

# Usage

```
$ slack-status
Usage: slack-status [options..] <template ID>

Options:
  -d  Dry run
  -i  Append information of the music playing on iTunes
  -w  Watch changes (with -i)

Templates:
  - home : :house: Working remotely
  - lunch : :fork_and_knife: Having lunch
```
