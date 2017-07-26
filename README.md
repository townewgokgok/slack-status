# slack-status

Changes your Slack user custom status from CLI

# Install

```bash
go install github.com/townewgokgok/slack-status
```

# Customize your settings

Your `token` can be created at [Slack "Legacy tokens" page](https://api.slack.com/custom-integrations/legacy-tokens).

```bash
vi $GOPATH/src/github.com/townewgokgok/slack-status/settings.yml
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
Usage: slack-status <template ID>

Templates:
- home : 🏡  Working remotely
- lunch : 🍴  Having lunch
```
