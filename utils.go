package main

import (
	"fmt"
	"strings"
)

func extractCommand(args []string) (command, error) {
	if len(args) < 2 {
		return command{}, fmt.Errorf("No commands with arguments were provided, exiting...")
	}
	commandName := args[1]
	var commandArgs []string
	if len(args) > 2 {
		commandArgs = args[2:]
	}
	return command{
		name: commandName,
		args: commandArgs,
	}, nil
}

func formatContentWithTitle(title []string, content [][]string) (string, error) {
	if len(content) == 0 {
		return "", nil
	}
	if len(title) != len(content[0]) {
		return "", fmt.Errorf("Title cannot has les entries that content! Title %v\nContent %v", title, content)
	}
	measures := make([]int, len(content[0]))
	bordersChars := 2
	for _, entry := range content {
		for idx, it := range entry {
			curLen := len(it) + bordersChars
			if curLen <= len(title[idx]) {
				curLen = len(title[idx]) + bordersChars*2
			}
			if curLen > measures[idx] {
				measures[idx] = curLen
			}
		}
	}
	titleBody := ""
	for idx, it := range title {
		titleBody += assembleFormattedString(it, measures[idx], "|", "-")
	}
	contentBody := ""
	for _, entry := range content {
		entryStr := ""
		for idx, it := range entry {
			entryStr += assembleFormattedString(it, measures[idx], " ", " ")
		}
		contentBody += entryStr + "\n"
	}
	return titleBody + "\n" + contentBody, nil
}

func assembleFormattedString(str string, length int, border string, padding string) string {
	borderLen := len(border)
	rem := length - borderLen - len(str)
	half := int(float64(rem) / 2)
	return fmt.Sprintf(
		"%v%s%s%s%v",
		border,
		strings.Repeat(padding, half),
		str,
		strings.Repeat(padding, half+rem-(2*half)),
		border,
	)
}
