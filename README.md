English | [日本語](README.ja.md)

# slack-status

[![MIT License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square)](https://github.com/goreleaser)

Updates your Slack user status from CLI.
Your own templates can be managed by editting the settings file.
Also has features to append the music information being played on your iTunes etc.

# Install

## manually

Download the binary from [the releases page](releases) and move it to a directory in your `PATH` variable.

## by [Homebrew](https://brew.sh/)

```
brew tap townewgokgok/tap
brew install slack-status
```

## by [Go 1.8](https://golang.org/)

```
go get github.com/townewgokgok/slack-status
```

# Configure your settings

Your settings file will be saved at `$HOME/.slack-status.yml`.

`slack-status edit` to edit it like

```yaml
slack:
  token: xoxp-...

templates:
  home: ':house: Working remotely today'
  lunch: ':fork_and_knife: Having lunch'
```

- The emoji at the beginning of the template string will be used as the status icon.
- Your `token` can be created at [Slack "Legacy tokens" page](https://api.slack.com/custom-integrations/legacy-tokens).

# Usage

`slack-status` can be used in subcommand style like `slack-status <command> ...`.

- `slack-status edit` … Opens your settings file in the editor
- `slack-status list` … Lists your templates
- `slack-status get` … Shows your current status
- `slack-status set [options...] <template ID>...` … Updates your status
  - **Options**
    - `--dryrun`, `-d` … just print the composed status text (your status will be not changed)
    - `--watch`, `-w` … watch changes (with `-i` or `-l`)
  - **Special template IDs**
    - `itunes` … appends information about the music playing on iTunes
    - `lastfm` … appends information about the music scrobbled to last.fm
- `slack-status help [<command>]` … Shows a list of commands or help for one command

# Examples

```
$ slack-status set home
🏠 Working remotely today
```

```
$ slack-status set lunch
🍴 Having lunch
```

```
$ slack-status set lunch home
🍴 Having lunch 🏠 Working remotely today
```

↑ `🍴` will be used as the emoji icon and the rest as the text.

```
$ slack-status set itunes
🎵 Satellite Young - Break! Break! Tic! Tac! (from "Satellite Young")
```

```
$ slack-status set home itunes
🏠 Working remotely today 🎵 Satellite Young - Break! Break! Tic! Tac! (from "Satellite Young")
```

```
$ slack-status set -w home itunes
[10:25:39] 🏠 Working remotely today 🎵 Satellite Young - Break! Break! Tic! Tac! (from "Satellite Young")
[10:30:16] 🏠 Working remotely today 🎵 Satellite Young - Geeky Boyfriend (from "Satellite Young")
[10:33:51] 🏠 Working remotely today 🎵 Satellite Young - AI Threnody (from "Satellite Young")
```
