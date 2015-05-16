# Yesterday

[![Build Status](https://travis-ci.org/mewmew/yesterday.svg?branch=master)](https://travis-ci.org/mewmew/yesterday)
[![Coverage Status](https://img.shields.io/coveralls/mewmew/yesterday.svg)](https://coveralls.io/r/mewmew/yesterday?branch=master)
[![GoDoc](https://godoc.org/github.com/mewmew/yesterday?status.svg)](https://godoc.org/github.com/mewmew/yesterday)

Yesterday is a procrastination tool which allows you to send emails up to 24 hours in the past.

It has two modes.

Without the -http flag, it runs in command-line mode and sends email from the terminal.

    yesterday -from="john.doe@student.uni.edu" -to="jane.roe@uni.edu" -subject="Report" -message="See attached report." report.pdf

With the -http flag, it runs as a web server and sends email from a web page.

    yesterday -http=:6565

## Installation

    go get github.com/mewmew/yesterday

## Usage

TODO

## Public domain

The source code and any original content of this repository is hereby released into the [public domain].

[public domain]: https://creativecommons.org/publicdomain/zero/1.0/
