Developed during the 3rd [HackPompey](https://twitter.com/hackpompey) hackathon.

# Yesterday

[![Build Status](https://travis-ci.org/mewmew/yesterday.svg?branch=master)](https://travis-ci.org/mewmew/yesterday)
[![Coverage Status](https://img.shields.io/coveralls/mewmew/yesterday.svg)](https://coveralls.io/r/mewmew/yesterday?branch=master)
[![GoDoc](https://godoc.org/github.com/mewmew/yesterday?status.svg)](https://godoc.org/github.com/mewmew/yesterday)

Yesterday is a procrastination tool which allows you to send emails up to 24 hours in the past.

It has two modes.

Without the -http flag, it runs in command-line mode and sends email from the terminal.

    yesterday -to="jane.roe@uni.edu" -subject="Report" -message="See attachment." report.pdf

With the -http flag, it runs as a web server and sends email from a web page.

    yesterday -http=:6565

## Installation

    go get github.com/mewmew/yesterday
    cp $GOPATH/src/github.com/mewmew/yesterday/example_auth.json auth.json
    # Edit username and password for SMTP authentication in auth.json

## Usage

```
yesterday [OPTION]... FILE...

Flags:
  -auth string
      JSON file with SMTP authentication information. (default "auth.json")
  -d  Enable debug output.
  -http string
      HTTP service address (e.g. ":6565").
  -message string
      Email message.
  -past duration
      Spoof date in number of hours in the past. (default 24h0m0s)
  -subject string
      Email subject.
  -to string
      Recipient email address.

```

## Web Interface

> *Yesterday, all my troubles seemed so far away*

![Yesterday, all my troubles seemed so far away](https://raw.githubusercontent.com/mewmew/yesterday/master/examples/yesterday_1.png)

> *Oh, I believe in yesterday*

![Oh, I believe in yesterday](https://raw.githubusercontent.com/mewmew/yesterday/master/examples/yesterday_2.png)

> *Yesterday, love was such an easy game to play*

![Yesterday, love was such an easy game to play](https://raw.githubusercontent.com/mewmew/yesterday/master/examples/yesterday_3.png)

## Public domain

The source code and any original content of this repository is hereby released into the [public domain].

[public domain]: https://creativecommons.org/publicdomain/zero/1.0/
