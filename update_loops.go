package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nlopes/slack"
)

type processFrame = func(frame string)

// LoopOverChannel is kinda generic
func LoopOverChannel(framesChan chan string, delay float64, noop bool, callback processFrame) {
	tickerChan := time.Tick(time.Millisecond * time.Duration(delay*1000))

	for frame := range framesChan {
		<-tickerChan
		if noop {
			fmt.Printf("\033[2K\r%s", frame)
		} else {
			callback(frame)
		}
	}
}

// LoopPostMessage posts a message and does stuff
func LoopPostMessage(framesChan chan string, channel string, delay float64, noop bool, slackAPI *slack.Client) {
	var dst, ts, txt string

	callback := func(frame string) {
		if dst == "" || ts == "" {

			var err error
			dst, ts, err = slackAPI.PostMessage(channel, slack.MsgOptionText(frame, false), slack.MsgOptionAsUser(true))
			if err != nil {
				log.Fatal("FATAL: Could not post initial frame to Slack: ", err)
			}
			log.Printf("initial frame %v/%v: %v\n", dst, ts, frame)

		} else {

			var err error
			_, _, txt, err = slackAPI.UpdateMessage(dst, ts, slack.MsgOptionText(frame, false))
			if err != nil {
				log.Printf("ERROR updating %v/%v with frame %v: %v", dst, ts, frame, err)
			} else {
				log.Printf("updated frame %v/%v: %v", dst, ts, txt)
			}
		}
	}

	LoopOverChannel(framesChan, delay, noop, callback)
}

// In order for this to work, paste this code in vendor\github.com\nlopes\slack\users.go:

// // SetUserName sets the current user's display name.
// // CUSTOM ADDED BY JIRI SWEN 2019
// func (api *Client) SetUserName(newUserName string) error {
// 	values := url.Values{
// 		"token": {api.token},
// 		"value": {newUserName},
// 	}
// 	values.Add("name", "display_name")

// 	resp := &getUserProfileResponse{}

// 	err := api.postMethod(context.Background(), "users.profile.set", values, &resp)

// 	if err != nil {
// 		return err
// 	}
// 	if err := resp.Err(); err != nil {
// 		fmt.Println(err)
// 		return err
// 	}

// 	return nil
// }

// LoopUpdateProfile is cool
func LoopUpdateProfile(framesChan chan string, delay float64, noop bool, slackAPI *slack.Client) {
	callback := func(frame string) {
		err := slackAPI.SetUserName(frame)
		if err != nil {
			log.Printf("ERROR updating with frame %v: %v", frame, err)
		} else {
			log.Printf("updated frame %v", frame)
		}
	}

	LoopOverChannel(framesChan, delay, noop, callback)
}
