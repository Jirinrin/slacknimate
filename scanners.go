package main

import (
	"bufio"
	"fmt"
	"os"
)

// ScanFrames reads from StdIn and when it encounters and EOF it stops
// and returns all lines in the form of an array.
func ScanFrames(backAndForth bool) []string {
	var frames []string
	reader := bufio.NewScanner(os.Stdin)
	for reader.Scan() {
		frames = append(frames, reader.Text())
	}
	if err := reader.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}

	if backAndForth {
		backAndForthFrames := frames
		for i := len(frames)-2; i > 0; i-- {
			backAndForthFrames = append(backAndForthFrames, frames[len(frames)-i-1])
		}
		return backAndForthFrames
	}
	return frames
}

// FramesIterator is cool
func FramesIterator(frames []string) chan string {
	ch := make(chan string)
	go func() {
		for _, frame := range frames {
			ch <- frame
		}
		close(ch)
	}()
	return ch
}

// LoopingFramesIterator is only suitable for input that will end, and will continue
// consuming memory while never sending anything if STDIN is a process that
// generates continuous output.
func LoopingFramesIterator(frames []string) chan string {
	ch := make(chan string)
	go func() {
		for {
			for _, frame := range frames {
				ch <- frame
			}
		}
	}()
	return ch
}
