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

func strlen(input string) int {
	return strings.Count(input, "") - 1
}

func createStringifiedTable(title []string, content [][]string) (string, error) {
	if len(content) == 0 {
		return "", nil
	}
	if len(title) != len(content[0]) {
		return "", fmt.Errorf("Title cannot has les entries that content! Title %v\nContent %v", title, content)
	}
	measures := make([]int, len(content[0]))
	for _, entry := range content {
		entryLen := len(entry)
		lastEntryIdx := entryLen - 1
		for idx, it := range entry {
			bordersChars := 2
			if idx != lastEntryIdx {
				bordersChars = 1
			}
			curLen := max(strlen(it)+bordersChars, strlen(title[idx])+bordersChars)
			if curLen > measures[idx] {
				measures[idx] = curLen
			}
		}
	}
	titleBody := assembleFormattedRow(title, measures, "|", "-")
	var contentBody strings.Builder
	for _, entry := range content {
		entryStr := assembleFormattedRow(entry, measures, " ", " ")
		contentBody.WriteString(entryStr + "\n")
	}
	return titleBody + "\n" + contentBody.String(), nil
}

func assembleFormattedRow(row []string, measures []int, border string, padding string) string {
	var res strings.Builder
	rowLen := len(row)
	lastIndex := rowLen - 1
	for idx, it := range row {
		leftBorder := border
		rightBorder := border
		if idx != lastIndex {
			rightBorder = ""
		}
		rem := measures[idx] - strlen(it) - (len(leftBorder) + len(rightBorder))
		if rem < 0 {
			res.WriteString(it)
			continue
		}
		half := int(float64(rem) / 2)
		paddingLen := len(padding)
		fmt.Fprintf(&res, "%v%s%s%s%v",
			leftBorder,
			strings.Repeat(padding, half/paddingLen),
			it,
			strings.Repeat(padding, (half+rem-(2*half))/paddingLen),
			rightBorder,
		)
	}
	return res.String()
}
