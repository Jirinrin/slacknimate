package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/nlopes/slack"
)

func main() {
	app := cli.NewApp()
	app.Name = "slacknimate"
	app.Usage = "text animation for Slack messages"
	app.Version = "1.0.1"
	app.UsageText = "slacknimate [options]"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "api-token, a",
			Usage:  "API token*",
			EnvVar: "SLACK_TOKEN",
		},
		cli.Float64Flag{
			Name:  "delay, d",
			Usage: "minimum delay between frames",
			Value: 1,
		},
		cli.StringFlag{
			Name:   "channel, c",
			Usage:  "channel/destination*",
			EnvVar: "SLACK_CHANNEL",
		},
		cli.BoolFlag{
			Name:  "loop, l",
			Usage: "loop content upon reaching end",
		},
		cli.BoolFlag{
			Name:  "preview",
			Usage: "preview on terminal instead of posting",
		},
		cli.BoolFlag{
			Name:  "backandforth, bf",
			Usage: "play the content forth and back",
		},
	}
	app.Action = func(c *cli.Context) {
		apiToken := c.String("api-token")

		channelSlice := strings.Split(c.String("channel"), "/")
		channel := channelSlice[len(channelSlice)-1]

		delay := c.Float64("delay")

		noop := c.Bool("preview")

		if !noop {
			stderr := log.New(os.Stderr, "", 0) // log to stderr with no timestamps
			if apiToken == "" {
				stderr.Fatal("API token is required.",
					" Use --api-token or set SLACK_TOKEN env variable.")
			}
			if channel == "" {
				stderr.Fatal("Destination is required.",
					" Use --channel or set SLACK_CHANNEL env variable.")
			}
			if delay < 0.001 {
				stderr.Fatal("You must have a delay >=0.001 to avoid creating a time paradox.")
			}
		}

		frames := ScanFrames(c.Bool("backandforth"), c.Bool("loop"))
		var framesChan chan string
		if c.Bool("loop") {
			framesChan = LoopingFramesIterator(frames)
		} else {
			framesChan = FramesIterator(frames)
		}

		api := slack.New(apiToken)

		LoopPostMessage(framesChan, channel, delay, noop, api)

		fmt.Println("\nDone!")
	}

	app.Run(os.Args)
}
