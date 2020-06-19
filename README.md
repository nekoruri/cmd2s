# cmd2s
Operate Google Chrome to automatically send a message to Slack.

## Description
- Reads the specified file and sends the string to Slack.
- You can automatically send using Google Chrome.
- This time, it is written assuming the format of `/feed subscribe [RSS Feed URL]`.
- If you want to use it for general purposes, please modify the source code.

***see also:***

- Getting Started with Headless Chrome 
	- https://developers.google.com/web/updates/2017/04/headless-chrome
- Add RSS feeds to Slack 
	- https://slack.com/help/articles/218688467-Add-RSS-feeds-to-Slack

## Features
- It is made by golang so it supports multi form.
- You can automatically send using Google Chrome.
- It simply sends the string in the file, so you can also send and run Slack slash command.
	- This source is specialized for `/feed subscribe [RSS Feed URL]`.
	- If you want to generalize, modify the source

## Requirement
- Go 1.14+
- Packages in use
	- chromedp/chromedp: A faster, simpler way to drive browsers supporting the Chrome DevTools Protocol.
		- https://github.com/chromedp/chromedp

## Usage
Just run the only one command.

```	sh
$ ./cmd2s
```

However, setting is necessary to execute.

### Setting Example

1. In the same place as the binary file create execution settings file.

1. Execution settings are done with `config.json` file.

```sh
{
	"LoginURL": "https://xxxxx.slack.com/sso/saml/start?redir=%2Fmessages",
	"LoginID": "YOUR_EMAIL",
	"LoginPass": "YOUR_PASSWORD",
	"ChannelURL": "https://app.slack.com/client/XXXXXXX/XXXXXXX",
	"CmdFile": "/cmd2s/_region_feed.txt"
}
```

- About setting items
	- `LoginURL`: String
		- Specify the URL to log in to Slack (This source is SMAL login in Azure AD)
	- `LoginID`: String
		- Specify the ID for logging in to Slack (This source is SMAL login in Azure AD)
	- `LoginPass`: String
		- Specify password to log in to Slack (This source is SMAL login in Azure AD).
	- `ChannelURL`: String
		- Specify the URL of the Slack channel to send the string to.
	- `CmdFile`: String
		- Specify the full path of the file that describes the transmission contents.

## Installation

If you build from source yourself.

```	console
$ go get -u github.com/uchimanajet7/cmd2s
$ go build
```

### When you want to know the execution result of the slash command.

This source does not acquire the sequential execution result, but it can be dealt with if changed.

In this source, we get the result of `/feed list` and confirm the execution.

```	console
Only visible to you
Slackbot  8:28 PM
ID: 1210548705408 - Title: Amazon Virtual Private Cloud (Osaka-Local) Service Status
URL: https://status.aws.amazon.com/rss/vpc-ap-northeast-3.rss
ID: 1186712303954 - Title: Amazon Simple Workflow Service (Osaka-Local) Service Status
URL: https://status.aws.amazon.com/rss/swf-ap-northeast-3.rss
ID: 1185333543893 - Title: Amazon Simple Queue Service (Osaka-Local) Service Status
URL: https://status.aws.amazon.com/rss/sqs-ap-northeast-3.rss
... more
```

## Author
[uchimanajet7](https://github.com/uchimanajet7)

## Licence
[Apache License 2.0](https://github.com/uchimanajet7/cmd2s/blob/master/LICENSE)

## As reference information
- AWS Service Health Dashboard のRSS Feed をSlack に自動登録する #aws #slack - uchimanajet7のメモ
	- https://uchimanajet7.hatenablog.com/entry/2020/06/19/180000
