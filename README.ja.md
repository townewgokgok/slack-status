[English](README.md) | 日本語

# slack-status

SlackのユーザステータスをCLIから変更するツールです。
設定ファイルを編集することで、自分専用のテンプレートを管理することができます。

# 必要環境

- [Go 1.8](https://golang.org/)

# インストール

```
go get github.com/townewgokgok/slack-status
```

# 設定ファイルの編集

設定ファイルは `$HOME/.slack-status.yml` に保存されます。

`slack-status -e` と入力するとエディタが起動します。以下のように項目を編集してください。

```yaml
slack:
  token: xoxp-...

templates:
  home:
    emoji: house
    text: 在宅作業中
  lunch:
    emoji: fork_and_knife
    text: お昼ごはん中
```

`token` は [Slackの "Legacy tokens" ページ](https://api.slack.com/custom-integrations/legacy-tokens) で発行することができます。

# 使用方法

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

# 使用例

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
