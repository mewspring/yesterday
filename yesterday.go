// TODO: Generate usage doc using go:generate.

// Yesterday is a procrastination tool which allows you to send emails up to 24
// hours in the past.
//
// It has two modes.
//
// Without the -http flag, it runs in command-line mode and sends email from the
// terminal.
//
//    yesterday -from="john.doe@student.uni.edu" -to="jane.roe@uni.edu" -subject="Report" -message="See attached report." report.pdf
//
// With the -http flag, it runs as a web server and sends email from a web page.
//
//    yesterday -http=:6565
package main

func main() {
}
