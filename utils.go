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

func createStringifiedTable(title []string, content [][]string) (string, error) {
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
			curLen := max(len(it)+bordersChars, len(title[idx])+bordersChars)
			if curLen > measures[idx] {
				measures[idx] = curLen
			}
		}
	}
	titleBody := assembleFormattedRow(title, measures, "|", "-")
	contentBody := ""
	for _, entry := range content {
		entryStr := assembleFormattedRow(entry, measures, " ", " ")
		contentBody += entryStr + "\n"
	}
	return titleBody + "\n" + contentBody, nil
}

func assembleFormattedRow(row []string, measures []int, border string, padding string) string {
	res := ""
	rowLen := len(row)
	lastIndex := rowLen - 1
	borderLen := len(border) * 2
	for idx, it := range row {
		leftBorder := border
		rightBorder := border
		if idx == 0 && rowLen > 2 {
			rightBorder = ""
		} else if idx == lastIndex {
			leftBorder = ""
		}
		rem := measures[idx] - len(it) - borderLen
		if rem < 0 {
			res += it
			continue
		}
		half := int(float64(rem) / 2)
		half /= len(padding)
		res += fmt.Sprintf(
			"%v%s%s%s%v",
			leftBorder,
			strings.Repeat(padding, half),
			it,
			strings.Repeat(padding, half+rem-(2*half)),
			rightBorder,
		)
	}
	return res
}

func assembleFormattedString(str string, length int, border string, padding string) string {
	if len(padding) > 1 {
		return str
	}
	borderLen := len(border) * 2
	rem := length - borderLen - len(str)
	if rem < 0 {
		return str
	}
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
