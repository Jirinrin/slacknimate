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
