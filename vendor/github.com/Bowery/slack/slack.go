// Copyright 2015 Bowery, Inc.

package slack

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

var (
	slackAddr = "https://slack.com/api"
)

const (
	postMessageURI = "chat.postMessage"
)

type slackPostMessageRes struct {
	Ok    bool
	Error string
}

// Client represents a slack api client. A Client is
// used for making requests to the slack api.
type Client struct {
	token string
}

// NewClient returns a Client with the provided api token.
func NewClient(token string) *Client {
	return &Client{token}
}

// SendMessage sends a text message to a specific channel
// with a specific username.
func (c *Client) SendMessage(channel, message, username string) error {
	if channel == "" || message == "" || username == "" {
		return errors.New("channel, message, and username required")
	}

	payload := url.Values{}
	payload.Set("token", c.token)
	payload.Set("channel", channel)
	payload.Set("text", message)
	payload.Set("username", username)

	res, err := http.PostForm(fmt.Sprintf("%s/%s", slackAddr, postMessageURI), payload)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody := new(slackPostMessageRes)
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(resBody)
	if err != nil {
		return err
	}

	if !resBody.Ok {
		return errors.New(resBody.Error)
	}

	return nil
}
