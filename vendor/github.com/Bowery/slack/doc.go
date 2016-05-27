// Package slack is a client library for Slack.
//
// Usage:
/*
  package main

  import (
    "github.com/Bowery/slack"
  )

  var (
    client *slack.Client
  )

  func main() {
    client = slack.NewClient("API_TOKEN")
    err := client.SendMessage("#mychannel", "message", "username")
    if err != nil {
      log.Fatal(err)
    }
  }
*/
package slack
