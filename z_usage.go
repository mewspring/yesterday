/*
Usage: yesterday [OPTION]... FILE...
Send emails up to 24 hours in the past.

Flags:
  -from string
    	Sender email address.
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
*/
package main
