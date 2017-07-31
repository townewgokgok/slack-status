[English](README.md) | 日本語

# slack-status

[![MIT License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square)](https://github.com/goreleaser)

SlackのユーザステータスをCLIから更新するツールです。
設定ファイルを編集することで、自分専用のテンプレートを管理することができます。

# 必要環境

- [Go 1.8](https://golang.org/)

# インストール

```
go get github.com/townewgokgok/slack-status
```

# 設定ファイルの編集

設定ファイルは `$HOME/.slack-status.yml` に保存されます。

`slack-status edit` と入力するとエディタが起動します。以下のように項目を編集してください。

```yaml
slack:
  token: xoxp-...

templates:
  home: ':house: 在宅作業中'
  lunch: ':fork_and_knife: お昼ごはん中'
```

- テンプレートの文頭にemojiがあるときは、ステータスアイコンとして使用されます。
- `token` は [Slackの "Legacy tokens" ページ](https://api.slack.com/custom-integrations/legacy-tokens) で発行することができます。

# 使用方法

`slack-status` コマンドは、 `slack-status <command> ...` のようなサブコマンドスタイルで使用します。

- `slack-status edit` … 設定ファイルをエディタで開きます
- `slack-status list` … テンプレート一覧を表示します
- `slack-status get` … 現在のユーザステータスを表示します
- `slack-status set [options...] [<template ID>]` … ユーザステータスを更新します
  - `--dryrun`, `-d` … ステータステキストの表示のみ（実際のステータスは変更されません）
  - `--itunes`, `-i` … iTunes で再生中の音楽情報を付加
  - `--lastfm`, `-l` … last.fm で再生中の音楽情報を付加
  - `--watch`, `-w` … 状態変化を監視（`-i` または `-l` と併せて使用）
- `slack-status help [<command>]` … コマンド一覧 または 指定されたコマンドのヘルプを表示します

# 使用例

```
$ slack-status set home
🏠 在宅作業中
```

```
$ slack-status set lunch
🍴 お昼ごはん中
```

```
$ slack-status set -i
🎵 Satellite Young - Break! Break! Tic! Tac! (from "Satellite Young")
```

```
$ slack-status set -i home
🏠 在宅作業中 🎵 Satellite Young - Break! Break! Tic! Tac! (from "Satellite Young")
```

```
$ slack-status set -i -w home
[10:25:39] 🏠 在宅作業中 🎵 Satellite Young - Break! Break! Tic! Tac! (from "Satellite Young")
[10:30:16] 🏠 在宅作業中 🎵 Satellite Young - Geeky Boyfriend (from "Satellite Young")
[10:33:51] 🏠 在宅作業中 🎵 Satellite Young - AI Threnody (from "Satellite Young")
```
